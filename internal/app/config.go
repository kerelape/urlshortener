package app

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config is the application configuration.
type Config struct {
	// ServerAddress is the address that the app runs on.
	ServerAddress string `env:"SERVER_ADDRESS"`

	// BaseURL is the base url for the shortened urls.
	BaseURL string `env:"BASE_URL"`

	// ShortenerPath is the root of the app.
	ShortenerPath string `env:"SHORTENER_PATH"`

	// APIPath is the root of the app's API.
	APIPath string `env:"API_PATH"`

	// FileStoragePath is path to the file database.
	FileStoragePath string `env:"FILE_STORAGE_PATH"`

	// DatabaseDSN is the DSN to connect to.
	DatabaseDSN string `env:"DATABASE_DSN"`

	// EnableHTTPS is used to enable https.
	EnableHTTPS bool `env:"ENABLE_HTTPS"`

	// ConfigFile is path to config file.
	ConfigFile string `env:"CONFIG"`

	// TrustedSubnet is ips for internal api.
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
}

// InitConfig initializes Config and returns it.
func InitConfig() (Config, error) {
	var environment Config
	parseError := env.Parse(&environment)
	if parseError != nil {
		return environment, parseError
	}
	flags := Config{}
	flag.StringVar(&flags.ServerAddress, "a", "localhost:8080", "Server address")
	flag.StringVar(&flags.BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&flags.FileStoragePath, "f", "urlshortener.db", "Path file DB")
	flag.StringVar(&flags.APIPath, "api-path", "/api", "API root")
	flag.StringVar(&flags.ShortenerPath, "app-path", "/", "Shortener root")
	flag.StringVar(&flags.DatabaseDSN, "d", "", "")
	flag.BoolVar(&flags.EnableHTTPS, "s", false, "Enable HTTPS")
	flag.StringVar(&flags.ConfigFile, "c", "", "Path to config file.")
	flag.StringVar(&flags.TrustedSubnet, "t", "", "Trusted subnet for internal apis.")
	flag.Parse()
	if environment.ServerAddress == "" {
		environment.ServerAddress = flags.ServerAddress
	}
	if environment.BaseURL == "" {
		environment.BaseURL = flags.BaseURL
	}
	if environment.FileStoragePath == "" {
		environment.FileStoragePath = flags.FileStoragePath
	}
	if environment.APIPath == "" {
		environment.APIPath = flags.APIPath
	}
	if environment.ShortenerPath == "" {
		environment.ShortenerPath = flags.ShortenerPath
	}
	if environment.DatabaseDSN == "" {
		environment.DatabaseDSN = flags.DatabaseDSN
	}
	if !environment.EnableHTTPS {
		environment.EnableHTTPS = flags.EnableHTTPS
	}
	if environment.ConfigFile == "" {
		environment.ConfigFile = flags.ConfigFile
	}
	if environment.TrustedSubnet == "" {
		environment.TrustedSubnet = flags.TrustedSubnet
	}
	if err := readConfigFile(&environment); err != nil {
		return Config{}, err
	}
	return environment, parseError
}

func readConfigFile(config *Config) error {
	if config.ConfigFile == "" {
		return nil
	}

	file, openError := os.Open(config.ConfigFile)
	if openError != nil {
		return openError
	}

	var cfg struct {
		ServerAddress   string `json:"server_address"`
		BaseURL         string `json:"base_url"`
		FileStoragePath string `json:"file_storage_path"`
		DatabaseDSN     string `json:"database_dsn"`
		EnableHTTPS     bool   `json:"enable_https"`
		TrustedSubnet   string `json:"trusted_subnet"`
	}
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return err
	}

	if config.ServerAddress == "" {
		config.ServerAddress = cfg.ServerAddress
	}
	if config.BaseURL == "" {
		config.BaseURL = cfg.BaseURL
	}
	if config.FileStoragePath == "" {
		config.FileStoragePath = cfg.FileStoragePath
	}
	if config.DatabaseDSN == "" {
		config.DatabaseDSN = cfg.DatabaseDSN
	}
	if !config.EnableHTTPS {
		config.EnableHTTPS = cfg.EnableHTTPS
	}
	if config.TrustedSubnet == "" {
		config.TrustedSubnet = cfg.TrustedSubnet
	}

	return nil
}
