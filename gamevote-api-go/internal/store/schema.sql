-- Applied on every start; everything is idempotent.

CREATE TABLE IF NOT EXISTS users (
    username   text PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    last_login timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS drink_types (
    id              integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name            text NOT NULL UNIQUE,
    volume_ml       integer NOT NULL,
    alcohol_percent real NOT NULL,
    beer_equivalent real NOT NULL,
    unit_name       text NOT NULL
);

CREATE TABLE IF NOT EXISTS parties (
    code          text PRIMARY KEY,
    status        text NOT NULL DEFAULT 'NOMINATION' CHECK (status IN ('NOMINATION', 'VOTING', 'RESULTS')),
    current_round integer NOT NULL DEFAULT 1,
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS attendees (
    party_code text NOT NULL REFERENCES parties ON DELETE CASCADE,
    username   text NOT NULL REFERENCES users,
    joined_at  timestamptz NOT NULL DEFAULT clock_timestamp(),
    PRIMARY KEY (party_code, username)
);

CREATE TABLE IF NOT EXISTS options (
    party_code text NOT NULL REFERENCES parties ON DELETE CASCADE,
    name       text NOT NULL,
    app_id     integer,
    image_url  text,
    created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    PRIMARY KEY (party_code, name)
);
CREATE UNIQUE INDEX IF NOT EXISTS options_party_lower_name ON options (party_code, lower(name));

-- One ballot per attendee and round; choices is {"option name": -1|0|1}.
CREATE TABLE IF NOT EXISTS ballots (
    party_code text NOT NULL REFERENCES parties ON DELETE CASCADE,
    round      integer NOT NULL,
    username   text NOT NULL,
    choices    jsonb NOT NULL,
    PRIMARY KEY (party_code, round, username)
);

CREATE TABLE IF NOT EXISTS beers (
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    party_code      text NOT NULL,
    attendee        text NOT NULL,
    drink_name      text NOT NULL DEFAULT 'Beer',
    beer_equivalent real NOT NULL DEFAULT 1.0,
    created_at      timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (party_code, attendee) REFERENCES attendees ON DELETE CASCADE
);
ALTER TABLE beers ADD COLUMN IF NOT EXISTS drink_name text NOT NULL DEFAULT 'Beer';
ALTER TABLE beers ADD COLUMN IF NOT EXISTS beer_equivalent real NOT NULL DEFAULT 1.0;
CREATE INDEX IF NOT EXISTS beers_party ON beers (party_code);

-- Default presets, only when the table is empty.
INSERT INTO drink_types (name, volume_ml, alcohol_percent, beer_equivalent, unit_name)
SELECT * FROM (VALUES
    ('Beer',                500,  5.0, 1.0, 'Halber Liter'),
    ('Vodka',                60, 40.0, 3.0, 'Shot'),
    ('Wein',                210, 12.0, 1.5, 'Achtel'),
    ('Spritzer',            420,  6.0, 1.5, 'Viertel'),
    ('Jägermeister',         70, 35.0, 3.5, 'Shot'),
    ('Berliner Luft',       140, 18.0, 1.0, 'Achtel'),
    ('Klopfer',             160, 16.0, 1.5, 'Achtel'),
    ('Gin',                  70, 40.0, 3.5, 'Shot'),
    ('Nussschnaps',          80, 30.0, 4.0, 'Shot'),
    ('Jägermeister Orange',  80, 35.0, 4.0, 'Shot'),
    ('Tatra Tea Black',      50, 52.0, 2.5, 'Shot'),
    ('Weißer Rum',           70, 37.5, 3.5, 'Shot'),
    ('Xuxu',                250, 15.0, 1.0, 'Viertel'),
    ('Jägermeister Winter', 100, 35.0, 5.0, 'Shot')
) AS v (name, volume_ml, alcohol_percent, beer_equivalent, unit_name)
WHERE NOT EXISTS (SELECT 1 FROM drink_types);
