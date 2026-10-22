package store

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"

	"gamevote-api-go/api"
)

// Unambiguous characters only (no 0/O, 1/I).
const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func newCode() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = codeAlphabet[int(b[i])%len(codeAlphabet)]
	}
	return string(b)
}

// ── reads ───────────────────────────────────────────────────────────────

func (s *Store) ListParties(ctx context.Context) ([]api.PartySummary, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT p.code, p.status, p.created_at,
		       (SELECT count(*) FROM attendees a WHERE a.party_code = p.code),
		       (SELECT count(*) FROM options o WHERE o.party_code = p.code),
		       (SELECT count(*) FROM beers b WHERE b.party_code = p.code)
		FROM parties p
		ORDER BY p.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []api.PartySummary{}
	for rows.Next() {
		var p api.PartySummary
		if err := rows.Scan(&p.Code, &p.Status, &p.CreatedAt, &p.AttendeeCount, &p.OptionCount, &p.BeerCount); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetParty loads the full state of one party using a single round trip.
func (s *Store) GetParty(ctx context.Context, code string) (*api.Party, error) {
	return loadParty(ctx, s.pool, code)
}

func loadParty(ctx context.Context, q querier, code string) (*api.Party, error) {
	// A pipelined batch: six small queries, one network round trip.
	b := &pgx.Batch{}
	b.Queue(`SELECT status, current_round, created_at FROM parties WHERE code = $1`, code)
	b.Queue(`SELECT username FROM attendees WHERE party_code = $1 ORDER BY joined_at, username`, code)
	b.Queue(`SELECT name, app_id, image_url FROM options WHERE party_code = $1 ORDER BY created_at, name`, code)
	b.Queue(`SELECT attendee, count(*) FROM beers WHERE party_code = $1 GROUP BY attendee`, code)
	b.Queue(`SELECT round, username, choices FROM ballots WHERE party_code = $1`, code)
	b.Queue(`SELECT id, attendee, drink_name, beer_equivalent, created_at FROM beers WHERE party_code = $1 ORDER BY created_at ASC, id ASC`, code)

	res := q.SendBatch(ctx, b)
	defer res.Close()

	p := &api.Party{
		Code:            code,
		Attendees:       []string{},
		Options:         []api.PartyOption{},
		PendingVoters:   []string{},
		RoundResults:    map[string]map[string]int{},
		BeerPerAttendee: map[string]int{},
		BeerLog:         []api.BeerEntry{},
	}
	if err := res.QueryRow().Scan(&p.Status, &p.CurrentRound, &p.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFound("party %s not found", code)
		}
		return nil, err
	}

	rows, err := res.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return nil, err
		}
		p.Attendees = append(p.Attendees, name)
		p.BeerPerAttendee[name] = 0
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = res.Query()
	if err != nil {
		return nil, err
	}
	optionNames := map[string]bool{}
	for rows.Next() {
		var o api.PartyOption
		if err := rows.Scan(&o.Name, &o.AppId, &o.ImageUrl); err != nil {
			rows.Close()
			return nil, err
		}
		optionNames[o.Name] = true
		p.Options = append(p.Options, o)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = res.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			name string
			n    int
		)
		if err := rows.Scan(&name, &n); err != nil {
			rows.Close()
			return nil, err
		}
		p.BeerPerAttendee[name] = n
		p.BeerCount += n
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = res.Query()
	if err != nil {
		return nil, err
	}
	voted := map[string]bool{} // voters in the current round
	for rows.Next() {
		var (
			round   int
			user    string
			choices map[string]int
		)
		if err := rows.Scan(&round, &user, &choices); err != nil {
			rows.Close()
			return nil, err
		}
		key := strconv.Itoa(round)
		scores := p.RoundResults[key]
		if scores == nil {
			scores = map[string]int{}
			p.RoundResults[key] = scores
		}
		for name, v := range choices {
			if optionNames[name] {
				scores[name] += v
			}
		}
		if round == p.CurrentRound {
			voted[user] = true
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = res.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e api.BeerEntry
		if err := rows.Scan(&e.Id, &e.Attendee, &e.DrinkName, &e.BeerEquivalent, &e.CreatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		p.BeerLog = append(p.BeerLog, e)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if p.Status == api.VOTING {
		for _, a := range p.Attendees {
			if !voted[a] {
				p.PendingVoters = append(p.PendingVoters, a)
			}
		}
	}
	return p, nil
}

// ── writes ──────────────────────────────────────────────────────────────

// CreateParty creates a party with a fresh code and the given attendees.
func (s *Store) CreateParty(ctx context.Context, attendees []string) (string, error) {
	for range 10 {
		code := newCode()
		err := s.createParty(ctx, code, attendees)
		if pgCode(err) == codeUniqueViolation {
			continue // code collision, try another
		}
		return code, err
	}
	return "", errors.New("could not generate a unique party code")
}

func (s *Store) createParty(ctx context.Context, code string, attendees []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `INSERT INTO parties (code) VALUES ($1)`, code); err != nil {
		return err
	}
	if len(attendees) > 0 {
		_, err := tx.Exec(ctx, `
			INSERT INTO attendees (party_code, username)
			SELECT $1, u FROM unnest($2::text[]) AS u
			ON CONFLICT DO NOTHING`, code, attendees)
		if pgCode(err) == codeForeignKey {
			return Invalid("unknown user in attendees")
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

type partyRow struct {
	Status api.PartyStatus
	Round  int
}

// inParty runs fn in a transaction holding a lock on the party row, which
// serialises all state changes of one party (votes vs. phase changes etc.).
func (s *Store) inParty(ctx context.Context, code string, fn func(tx pgx.Tx, p partyRow) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var p partyRow
	err = tx.QueryRow(ctx, `SELECT status, current_round FROM parties WHERE code = $1 FOR UPDATE`, code).
		Scan(&p.Status, &p.Round)
	if errors.Is(err, pgx.ErrNoRows) {
		return NotFound("party %s not found", code)
	}
	if err != nil {
		return err
	}
	if err := fn(tx, p); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) AddAttendee(ctx context.Context, code, username string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO attendees (party_code, username) VALUES ($1, $2)
		ON CONFLICT DO NOTHING`, code, username)
	if pgCode(err) == codeForeignKey {
		// Either the party or the user is unknown.
		var exists bool
		if qerr := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM parties WHERE code = $1)`, code).Scan(&exists); qerr != nil {
			return qerr
		}
		if !exists {
			return NotFound("party %s not found", code)
		}
		return Invalid("unknown user %q", username)
	}
	return err
}

func (s *Store) AddOption(ctx context.Context, code string, o api.PartyOption) error {
	return s.inParty(ctx, code, func(tx pgx.Tx, p partyRow) error {
		if p.Status != api.NOMINATION {
			return Conflict("games can only be changed during nomination")
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO options (party_code, name, app_id, image_url) VALUES ($1, $2, $3, $4)`,
			code, o.Name, o.AppId, o.ImageUrl)
		if pgCode(err) == codeUniqueViolation {
			return Conflict("%q is already nominated", o.Name)
		}
		return err
	})
}

func (s *Store) RemoveOption(ctx context.Context, code, name string) error {
	return s.inParty(ctx, code, func(tx pgx.Tx, p partyRow) error {
		if p.Status != api.NOMINATION {
			return Conflict("games can only be changed during nomination")
		}
		tag, err := tx.Exec(ctx, `DELETE FROM options WHERE party_code = $1 AND name = $2`, code, name)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return NotFound("%q is not nominated", name)
		}
		return nil
	})
}

var allowedTransitions = map[api.PartyStatus][]api.PartyStatus{
	api.NOMINATION: {api.VOTING},
	api.VOTING:     {api.NOMINATION, api.RESULTS},
	api.RESULTS:    {api.NOMINATION}, // starts the next round
}

func (s *Store) SetStatus(ctx context.Context, code string, to api.PartyStatus) error {
	return s.inParty(ctx, code, func(tx pgx.Tx, p partyRow) error {
		if p.Status == to {
			return nil
		}
		allowed := false
		for _, t := range allowedTransitions[p.Status] {
			allowed = allowed || t == to
		}
		if !allowed {
			return Conflict("cannot go from %s to %s", p.Status, to)
		}
		if to == api.VOTING {
			var n int
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM options WHERE party_code = $1`, code).Scan(&n); err != nil {
				return err
			}
			if n == 0 {
				return Conflict("nominate at least one game before voting")
			}
		}
		round := p.Round
		if p.Status == api.RESULTS && to == api.NOMINATION {
			round++
		}
		_, err := tx.Exec(ctx, `UPDATE parties SET status = $2, current_round = $3 WHERE code = $1`, code, to, round)
		return err
	})
}

// SubmitVotes stores the ballot and, once every attendee has voted, moves the
// party on to RESULTS. All in one transaction so the last two voters cannot race.
func (s *Store) SubmitVotes(ctx context.Context, code, username string, votes api.Votes) error {
	for name, v := range votes {
		if v < -1 || v > 1 {
			return Invalid("vote for %q must be -1, 0 or 1", name)
		}
	}
	return s.inParty(ctx, code, func(tx pgx.Tx, p partyRow) error {
		if p.Status != api.VOTING {
			return Conflict("voting is not open")
		}
		var attends bool
		if err := tx.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM attendees WHERE party_code = $1 AND username = $2)`,
			code, username).Scan(&attends); err != nil {
			return err
		}
		if !attends {
			return Invalid("%q is not an attendee of this party", username)
		}

		// Store an explicit value for every option; unknown names are dropped.
		rows, err := tx.Query(ctx, `SELECT name FROM options WHERE party_code = $1`, code)
		if err != nil {
			return err
		}
		choices := map[string]int{}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				rows.Close()
				return err
			}
			choices[name] = votes[name]
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}
		blob, err := json.Marshal(choices)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO ballots (party_code, round, username, choices) VALUES ($1, $2, $3, $4)
			ON CONFLICT (party_code, round, username) DO UPDATE SET choices = EXCLUDED.choices`,
			code, p.Round, username, blob); err != nil {
			return err
		}

		var pending int
		if err := tx.QueryRow(ctx, `
			SELECT count(*) FROM attendees a
			WHERE a.party_code = $1 AND NOT EXISTS (
				SELECT 1 FROM ballots b
				WHERE b.party_code = a.party_code AND b.round = $2 AND b.username = a.username)`,
			code, p.Round).Scan(&pending); err != nil {
			return err
		}
		if pending == 0 {
			_, err = tx.Exec(ctx, `UPDATE parties SET status = 'RESULTS' WHERE code = $1`, code)
		}
		return err
	})
}

func (s *Store) AddBeer(ctx context.Context, code, attendee, drinkName string, beerEquivalent float32) error {
	if drinkName == "" {
		drinkName = "Beer"
	}
	if beerEquivalent <= 0 {
		beerEquivalent = 1.0
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO beers (party_code, attendee, drink_name, beer_equivalent) VALUES ($1, $2, $3, $4)`, code, attendee, drinkName, beerEquivalent)
	if pgCode(err) == codeForeignKey {
		var exists bool
		if qerr := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM parties WHERE code = $1)`, code).Scan(&exists); qerr != nil {
			return qerr
		}
		if !exists {
			return NotFound("party %s not found", code)
		}
		return Invalid("%q is not an attendee of this party", attendee)
	}
	return err
}
