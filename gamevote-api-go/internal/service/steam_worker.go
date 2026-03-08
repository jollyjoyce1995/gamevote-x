package service

import (
	"fmt"
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/storage"
	"log/slog"
	"os"
	"time"

	"github.com/Jleagle/steam-go/steamapi"
)

type SteamWorker struct {
	GameRepo *storage.GameRepository
	client   *steamapi.Client
}

func NewSteamWorker(gameRepo *storage.GameRepository) *SteamWorker {
	apiKey := os.Getenv("STEAM_API_KEY")
	client := steamapi.NewClient()
	client.SetKey(apiKey)

	return &SteamWorker{
		GameRepo: gameRepo,
		client:   client,
	}
}

// Start launches the background goroutine that syncs Steam games daily
func (w *SteamWorker) Start() {
	go func() {
		// Run once on startup if the DB is empty
		w.GameRepo.InitTable()
		count, err := w.GameRepo.Count()
		if err != nil || count == 0 {
			slog.Info("Steam Worker: fetching initial game list...")
			if err := w.FetchAndStore(); err != nil {
				slog.Error("Steam Worker: initial fetch failed", "error", err)
			}
		} else {
			slog.Info("Steam Worker: skipping initial fetch", "count", count)
		}

		// Daily refresh as a background ticker
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			slog.Info("Steam Worker: daily refresh starting...")
			if err := w.FetchAndStore(); err != nil {
				slog.Error("Steam Worker: daily refresh failed", "error", err)
			}
		}
	}()
}

// FetchAndStore downloads the Steam app list using IStoreService/GetAppList (paginated)
func (w *SteamWorker) FetchAndStore() error {
	if os.Getenv("STEAM_API_KEY") == "" {
		return fmt.Errorf("STEAM_API_KEY is not set. IStoreService requires an API key")
	}

	slog.Info("Steam Worker: starting paginated fetch via IStoreService...")

	// Clear existing before reinserting
	// Note: For a very large DB, we might want a more sophisticated sync (upsert),
	// but keeping the original "clear and reload" logic for now.
	if err := w.GameRepo.DeleteAll(); err != nil {
		return fmt.Errorf("failed to clear games table: %w", err)
	}

	lastAppID := 0
	maxResults := 10000 // Library default is 10k, max is 50k for the API
	totalStored := 0

	for {
		// GetAppList(limit int, offset int, afterDate int64, language LanguageCode)
		// offset is used as last_appid in the library implementation
		appList, err := w.client.GetAppList(maxResults, lastAppID, 0, "")
		if err != nil {
			return fmt.Errorf("failed to fetch app list at offset %d: %w", lastAppID, err)
		}

		if len(appList.Apps) == 0 {
			break
		}

		games := make([]models.Game, 0, len(appList.Apps))
		for _, app := range appList.Apps {
			if app.Name == "" {
				continue
			}
			games = append(games, models.Game{
				AppID:    app.AppID,
				Name:     app.Name,
				ImageURL: fmt.Sprintf("https://cdn.akamai.steamstatic.com/steam/apps/%d/header.jpg", app.AppID),
			})
		}

		if len(games) > 0 {
			if err := w.GameRepo.BulkInsert(games); err != nil {
				slog.Error("Steam Worker: batch insert error", "appid", lastAppID, "error", err)
			}
			totalStored += len(games)
		}

		slog.Info("Steam Worker: processed games", "stored", totalStored, "last_appid", appList.LastAppID)

		if !appList.HaveMoreResults {
			break
		}
		lastAppID = appList.LastAppID
	}

	slog.Info("Steam Worker: sync complete", "total_stored", totalStored)
	return nil
}

// Search finds games matching a query string
func (w *SteamWorker) Search(query string) ([]models.Game, error) {
	return w.GameRepo.Search(query, 20)
}
