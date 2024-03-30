package gonfik

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// LoaderConfig is an alais for map[string]any that is being used in gonfik as the values from configuration files
type LoaderConfig map[string]any

// Gonfik is the interface type to fetch configurations by key after initialization
//
// Config method takes an argument(key) and returns corresponding value and a boolean if not found
// Key can be a dot seperated value e.g. gonfik.path is a key for nested path key as an element of the gonfik key
type Gonfik interface {
	Config(keyPath string) (string, bool)
}

type gonfik struct {
	config   LoaderConfig
	fileName string
}

var konfik *gonfik

func init() {
	if err := loadDotEnv(); err != nil {
		fmt.Printf("DotEnv Configuration is not loaded: %s", err)
	}
}

func loadDotEnv() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %s", err)
	}

	envFiles, err := filepath.Glob(filepath.Join("gonfik.env"))
	if err != nil {
		return fmt.Errorf("error listing .env files: %s", err)
	}
	configFiles, err := filepath.Glob(filepath.Join(currentDir+"/config", "*.env"))
	if err != nil {
		return fmt.Errorf("error listing .env files in config dir: %s", err)
	}

	envFiles = append(envFiles, configFiles...)
	// Load .env file
	err = godotenv.Load(envFiles...)
	if err != nil {
		return fmt.Errorf("error loading .env file: %s", err)

	}
	return nil
}

// NewConfig is a constructor for a gonfik using directory and fileName
func NewConfig(configDir string, fileName string) (Gonfik, error) {
	if fileName == "" {
		fileName = getFileName()
	}
	if configDir == "" {
		configDir = getConfigDir()
	}
	return loadConfig(filepath.Join(configDir, fileName))
}

// GlobalConfig is an initialization function for configuration of the current project
// if there will not be multiple projects etc. it is better to use this
// since there is an instance of gonfik stored and used instead of creating a new one everytime
func GlobalConfig() (Gonfik, error) {
	if konfik == nil {
		config, err := NewConfig("", "")
		if err != nil {
			return nil, err
		}
		konfik = config.(*gonfik)
	}
	return konfik, nil
}

// Config is a getter function for a configuration value using its key
// key can be dot seperated for nested keys e.g. gonfik.filepath
func (c *gonfik) Config(keyPath string) (string, bool) {
	keys := strings.Split(keyPath, ".")
	currentObj := konfik.config
	configVal := ""
	for _, key := range keys {
		value, ok := currentObj[key]
		if !ok {
			return "", false
		}

		if nextData, isMap := value.(map[string]interface{}); isMap {
			currentObj = nextData
		} else {
			return fmt.Sprintf("%v", value), true
		}
	}
	return configVal, true
}

func loadConfig(fileName string) (*gonfik, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting executable directory: %s", err)
	}
	jsonFile, err := os.Open(filepath.Join(currentDir, fileName))
	if err != nil {
		return nil, fmt.Errorf("error creating gonfik: %s", err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading Json File: %s", err)
	}
	konfik = &gonfik{make(LoaderConfig),
		fileName}
	err = json.Unmarshal(byteValue, &((*konfik).config))
	if err != nil {
		return nil, err
	}
	return konfik, nil
}

func getConfigDir() string {
	dir, ok := getConfigFromEnv("GONFIK_DIR")
	if !ok {
		return "/config"
	}
	return dir
}

func getFileName() string {
	isProd, _ := getConfigFromEnv("GONFIK_IS_PROD")
	if isProd == "true" {
		fileName, _ := getConfigFromEnv("GONFIK_PROD_FILE")
		if fileName != "" {
			return fileName
		}
	} else {
		fileName, _ := getConfigFromEnv("GONFIK_DEV_FILE")
		if fileName != "" {
			return fileName
		}
	}
	return "application.json"
}

func getConfigFromEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
