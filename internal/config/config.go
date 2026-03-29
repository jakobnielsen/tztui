package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const appName = "tztui"

// Config holds all user-configurable settings persisted to disk.
type Config struct {
	Favourites []string `toml:"favourites"`
}

func configPath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("home dir: %w", err)
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, appName, "config.toml"), nil
}

// Load reads the configuration from disk.
// If the file does not exist a default (empty) Config is returned without error.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	var cfg Config
	_, err = toml.DecodeFile(path, &cfg)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{Favourites: []string{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return &cfg, nil
}

// Save writes the configuration to disk, creating intermediate directories as needed.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	return nil
}

// AddFavourite inserts iana into cfg.Favourites if not already present.
func (cfg *Config) AddFavourite(iana string) {
	for _, f := range cfg.Favourites {
		if f == iana {
			return
		}
	}
	cfg.Favourites = append(cfg.Favourites, iana)
}

// RemoveFavourite removes iana from cfg.Favourites if present.
func (cfg *Config) RemoveFavourite(iana string) {
	filtered := cfg.Favourites[:0]
	for _, f := range cfg.Favourites {
		if f != iana {
			filtered = append(filtered, f)
		}
	}
	cfg.Favourites = filtered
}

// IsFavourite reports whether iana is in cfg.Favourites.
func (cfg *Config) IsFavourite(iana string) bool {
	for _, f := range cfg.Favourites {
		if f == iana {
			return true
		}
	}
	return false
}
