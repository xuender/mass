package pb

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"github.com/xuender/kit/los"
)

func NewConfig() *Config {
	cfg := &Config{}

	los.Must0(viper.Unmarshal(cfg))

	return cfg
}

func (p *Config) Save() {
	path := viper.ConfigFileUsed()
	if path == "" {
		path = filepath.Join(los.Must(os.UserHomeDir()), "mass.toml")
	}

	file := los.Must(os.Create(path))
	defer file.Close()

	encoder := toml.NewEncoder(file)
	encoder.SetArraysMultiline(true)
	los.Must0(encoder.Encode(p))

	slog.Info("Save config", "file", path)
}
