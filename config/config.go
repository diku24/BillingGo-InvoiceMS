package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort       string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabaseName     string
	NetworkProtocol  string
}

func LoadDefault() (Config, error) {
	if path := os.Getenv("BILLINGGO_CONFIG_FILE"); path != "" {
		return Load(path)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("get working directory: %w", err)
	}
	if path, err := findConfigFile(workingDirectory); err == nil {
		return Load(path)
	}

	executable, err := os.Executable()
	if err == nil {
		if path, err := findConfigFile(filepath.Dir(executable)); err == nil {
			return Load(path)
		}
	}

	return Config{}, fmt.Errorf("configuration file config/config.env not found")
}

func Load(path string) (Config, error) {
	settings := viper.New()
	settings.SetConfigFile(path)
	if err := settings.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read configuration %q: %w", path, err)
	}

	value := func(key string) (string, error) {
		result := strings.TrimSpace(settings.GetString(key))
		if result == "" {
			return "", fmt.Errorf("missing required configuration %q", key)
		}
		return result, nil
	}

	serverPort, err := value("SERVERPORT")
	if err != nil {
		return Config{}, err
	}
	if err := validateListenAddress(serverPort); err != nil {
		return Config{}, fmt.Errorf("invalid SERVERPORT: %w", err)
	}
	databaseUser, err := value("DATABASEUSER")
	if err != nil {
		return Config{}, err
	}
	databasePassword, err := value("DATABASEPASS")
	if err != nil {
		return Config{}, err
	}
	databaseHost, err := value("DBCONNECTION")
	if err != nil {
		return Config{}, err
	}
	databaseName, err := value("DATABASE")
	if err != nil {
		return Config{}, err
	}
	networkProtocol, err := value("NETPROTOCOL")
	if err != nil {
		return Config{}, err
	}
	if networkProtocol != "tcp" && networkProtocol != "tcp4" && networkProtocol != "tcp6" {
		return Config{}, fmt.Errorf("invalid NETPROTOCOL %q", networkProtocol)
	}

	return Config{
		ServerPort:       serverPort,
		DatabaseUser:     databaseUser,
		DatabasePassword: databasePassword,
		DatabaseHost:     databaseHost,
		DatabaseName:     databaseName,
		NetworkProtocol:  networkProtocol,
	}, nil
}

func validateListenAddress(address string) error {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("%q must be a host:port listen address: %w", address, err)
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("port %q must be an integer from 1 to 65535", port)
	}
	return nil
}

func findConfigFile(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve start directory: %w", err)
	}

	for {
		candidate := filepath.Join(current, "config", "config.env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("config/config.env not found from %q", start)
}
