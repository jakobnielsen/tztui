package config

import (
	"os"
	"path/filepath"
	"testing"
)

func withTempConfig(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	withTempConfig(t)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Favourites) != 0 {
		t.Fatalf("expected empty favourites, got %v", cfg.Favourites)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	withTempConfig(t)
	cfg := &Config{Favourites: []string{"Europe/London", "Asia/Tokyo"}}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.Favourites) != 2 {
		t.Fatalf("expected 2 favourites, got %d", len(got.Favourites))
	}
	if got.Favourites[0] != "Europe/London" || got.Favourites[1] != "Asia/Tokyo" {
		t.Errorf("unexpected favourites: %v", got.Favourites)
	}
}

func TestLoad_InvalidTOML(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgDir := filepath.Join(dir, "tztui")
	if err := os.MkdirAll(cfgDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte("favourites = [unclosed"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid TOML, got nil")
	}
}

func TestAddFavourite(t *testing.T) {
	cfg := &Config{}
	cfg.AddFavourite("Europe/London")
	cfg.AddFavourite("Europe/London") // duplicate
	if len(cfg.Favourites) != 1 {
		t.Fatalf("expected 1 favourite, got %d", len(cfg.Favourites))
	}
}

func TestRemoveFavourite(t *testing.T) {
	cfg := &Config{Favourites: []string{"Europe/London", "Asia/Tokyo"}}
	cfg.RemoveFavourite("Europe/London")
	if len(cfg.Favourites) != 1 || cfg.Favourites[0] != "Asia/Tokyo" {
		t.Fatalf("unexpected favourites: %v", cfg.Favourites)
	}
}

func TestIsFavourite(t *testing.T) {
	cfg := &Config{Favourites: []string{"UTC"}}
	if !cfg.IsFavourite("UTC") {
		t.Error("expected UTC to be favourite")
	}
	if cfg.IsFavourite("Asia/Tokyo") {
		t.Error("expected Asia/Tokyo to not be favourite")
	}
}
