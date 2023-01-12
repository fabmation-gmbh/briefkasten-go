package config

import (
	"io/ioutil"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/pkg/errors"
)

// Config is the configuration structure.
type Config struct {
	General struct {
		// Listen is the IP and port where the server should listen to.
		Listen string `koanf:"listen"`
		// Environment defines the running environment for the log output.
		Environment string `koanf:"environment"`
		// EnableCompression enables compression of HTTP communication.
		EnableCompression bool `koanf:"enable_compression"`
		// CompressionLevel is the level of compression used, if enabled.
		// The following values can be used:
		//
		//  LevelDefault:          0
		//  LevelBestSpeed:        1
		//  LevelBestCompression:  2
		CompressionLevel uint `koanf:"compression_level"`
		// SecureCookie defines if cookies shall be marked as "secure".
		SecureCookie bool `koanf:"secure_cookie"`
	} `koanf:"general"`
	Debug struct {
		EnableSQLDebug bool `koanf:"enable_sql_debug"`
		EnableTracing  bool `koanf:"enable_tracing"`
		Tracing        struct {
			// Tracing endpoint. For example host+port of Jaeger.
			Endpoint string `koanf:"endpoint"`
		} `koanf:"tracing"`
	} `koanf:"debug"`
	// DB holds the database specific configuration parameter.
	DB struct {
		// URI is the connection URI.
		URI string `koanf:"uri"`
	} `koanf:"db"`
}

// C holds the current configuration.
var C Config

var (
	k      = koanf.New(".")
	parser = yaml.Parser()
)

// LoadConfig loads and parses the configuration from the given path.
func LoadConfig(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err = ioutil.WriteFile(path, []byte(""), 0o644); err != nil {
			return errors.Wrap(err, "unable to create config file")
		}
	} else if err != nil {
		return errors.Wrap(err, "unable to create config file")
	}

	loadDefaultValues()

	if err := k.Load(file.Provider(path), parser); err != nil {
		return errors.Wrap(err, "unable to load file from path")
	}

	if err := k.Unmarshal("", &C); err != nil {
		return errors.Wrap(err, "unable to parse configuration file")
	}

	return nil
}

func loadDefaultValues() {
	k.Load(confmap.Provider(map[string]any{}, "."), nil)
}
