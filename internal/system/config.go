package system

import (
	"errors"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

type Config struct {
	SiteName     string
	Location     types.Location
	Coordinates  models.Coordinate
	EnableSignUp bool
}

var defaultConfig = Config{
	SiteName: "",
	Location: types.NewLocation(time.Local),
	Coordinates: models.Coordinate{
		Latitude:  0,
		Longitude: 0,
	},
	EnableSignUp: true,
}

func readConfig(filePath string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(filePath, &cfg)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultConfig, nil
		}
		return Config{}, err
	}
	return cfg, nil
}

func writeConfig(filePath string, cfg Config) error {
	filePathTmp := filePath + ".tmp"
	file, err := os.OpenFile(filePathTmp, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		file.Close()
		return err
	}
	file.Close()

	return os.Rename(filePathTmp, filePath)
}

func NewConfigProvider(filePath string) (ConfigProvider, error) {
	if exist, err := core.FileExists(filePath); err != nil {
		return ConfigProvider{}, err
	} else if !exist {
		if err := writeConfig(filePath, defaultConfig); err != nil {
			return ConfigProvider{}, err
		}
	}
	return ConfigProvider{
		filePath: filePath,
	}, nil
}

type ConfigProvider struct {
	filePath string
}

func (p ConfigProvider) GetConfig() (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(p.filePath, &cfg)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultConfig, nil
		}
		return Config{}, err
	}
	return cfg, err
}

func (p ConfigProvider) UpdateConfig(fn func(cfg Config) (Config, error)) error {
	cfg, err := p.GetConfig()
	if err != nil {
		return err
	}

	cfg, err = fn(cfg)
	if err != nil {
		return err
	}

	return writeConfig(p.filePath, cfg)
}
