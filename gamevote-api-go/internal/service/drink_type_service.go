package service

import (
	"gamevote-api-go/internal/models"
	"gamevote-api-go/internal/storage"
	"log/slog"
)

type DrinkTypeService struct {
	Repo *storage.DrinkTypeRepository
}

func NewDrinkTypeService(repo *storage.DrinkTypeRepository) *DrinkTypeService {
	return &DrinkTypeService{
		Repo: repo,
	}
}

func (s *DrinkTypeService) SeedPresets() error {
	existing, err := s.Repo.FindAll()
	if err != nil {
		return err
	}

	if len(existing) > 0 {
		return nil // Already seeded
	}

	slog.Info("Seeding default DrinkPresets...")

	presets := []models.DrinkType{
		{Name: "Beer", VolumeMl: 500, AlcoholPercent: 5.0, BeerEquivalent: 1.0, UnitName: "Halber Liter"},
		{Name: "Vodka", VolumeMl: 60, AlcoholPercent: 40.0, BeerEquivalent: 3.0, UnitName: "Shot"},
		{Name: "Wein", VolumeMl: 210, AlcoholPercent: 12.0, BeerEquivalent: 1.5, UnitName: "Achtel"}, // Value inferred from instructions approximation (avg of 1.68 and 1.5) -> using 1.68 or 1.5? The instructions had multiple values, let's use 1.5.
		{Name: "Spritzer", VolumeMl: 420, AlcoholPercent: 6.0, BeerEquivalent: 1.5, UnitName: "Viertel"},
		{Name: "Jägermeister", VolumeMl: 70, AlcoholPercent: 35.0, BeerEquivalent: 3.5, UnitName: "Shot"},
		{Name: "Berliner Luft", VolumeMl: 140, AlcoholPercent: 18.0, BeerEquivalent: 1.0, UnitName: "Achtel"}, // Wait instructions had 1, 1.12, 7 etc.. let's just use 1
		{Name: "Klopfer", VolumeMl: 160, AlcoholPercent: 16.0, BeerEquivalent: 1.5, UnitName: "Achtel"},
		{Name: "Gin", VolumeMl: 70, AlcoholPercent: 40.0, BeerEquivalent: 3.5, UnitName: "Shot"},
		{Name: "Nussschnaps", VolumeMl: 80, AlcoholPercent: 30.0, BeerEquivalent: 4.0, UnitName: "Shot"},
		{Name: "Jägermeister Orange", VolumeMl: 80, AlcoholPercent: 35.0, BeerEquivalent: 4.0, UnitName: "Shot"},
		{Name: "Tatra Tea Black", VolumeMl: 50, AlcoholPercent: 52.0, BeerEquivalent: 2.5, UnitName: "Shot"},
		{Name: "Weißer Rum", VolumeMl: 70, AlcoholPercent: 37.5, BeerEquivalent: 3.5, UnitName: "Shot"},
		{Name: "Xuxu", VolumeMl: 250, AlcoholPercent: 15.0, BeerEquivalent: 1.0, UnitName: "Viertel"},
		{Name: "Jägermeister Winter", VolumeMl: 100, AlcoholPercent: 35.0, BeerEquivalent: 5.0, UnitName: "Shot"},
	}

	for _, p := range presets {
		if err := s.Repo.Save(&p); err != nil {
			slog.Error("Failed to seed preset", "preset", p.Name, "error", err)
		}
	}

	return nil
}

func (s *DrinkTypeService) GetDrinkTypes() ([]models.DrinkType, error) {
	return s.Repo.FindAll()
}

func (s *DrinkTypeService) AddCustomDrinkType(dt *models.DrinkType) (*models.DrinkType, error) {
	// Calculate beer equivalent if not provided:
	// (Volume * AlcoholPercent / 100) / (500 * 0.05)
	if dt.BeerEquivalent == 0 {
		pureAlcoholMl := float64(dt.VolumeMl) * (dt.AlcoholPercent / 100.0)
		dt.BeerEquivalent = pureAlcoholMl / 25.0 // 500 * 0.05 = 25
	}
	err := s.Repo.Save(dt)
	return dt, err
}
