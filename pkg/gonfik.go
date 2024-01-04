package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type LoaderConfig map[string]interface{}

type Gonfik struct {
	config   LoaderConfig
	fileName string
}

var konfik *Gonfik

func NewConfig() (*Gonfik, error) {
	return NewConfig2(getFileName())
}

func NewConfig2(fileName string) (*Gonfik, error) {
	return NewConfig3(getConfigDir(), fileName)
}

// this not-letting-overloading is bs.
func NewConfig3(configDir string, fileName string) (*Gonfik, error) {
	return loadConfig(configDir + fileName)
}

func GlobalConfig() (*Gonfik, error) {
	if konfik == nil {
		config, err := NewConfig()
		if err == nil {
			return nil, err
		}
		konfik = config
	}
	return konfik, nil
}

func loadConfig(fileName string) (*Gonfik, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating Gonfik: %s", err))
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading Json File: %s", err))
	}
	json.Unmarshal(byteValue, &konfik)
	return konfik, nil
}

func (c *Gonfik) Config(keyPath string) (string, bool) {
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

func getConfigDir() string {
	dir := getConfigFromEnv("GONFIK_DIR")
	if dir == "" {
		return "/config"
	}
	return dir
}

func getFileName() string {
	isProd := getConfigFromEnv("GONFIK_IS_PROD")
	if isProd == "true" || isProd == "" {
		fileName := getConfigFromEnv("GONFIK_PROD_FILE")
		if fileName != "" {
			return fileName
		}
	} else {
		fileName := getConfigFromEnv("GONFIK_DEV_FILE")
		if fileName != "" {
			return fileName
		}
	}
	return "application.json"
}

func loadDotEnv() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.New(fmt.Sprintf("Error getting current working directory: %s", err))
	}

	envFiles, err := filepath.Glob(filepath.Join("gonfik.env"))
	if err != nil {
		return errors.New(fmt.Sprintf("Error listing .env files: %s", err))
	}
	configFiles, err := filepath.Glob(filepath.Join(currentDir+"/config", "*.env"))
	if err != nil {
		return errors.New(fmt.Sprintf("Error listing .env files in config dir: %s", err))
	}

	envFiles = append(envFiles, configFiles...)
	// Load .env file
	err = godotenv.Load(envFiles...)
	if err != nil {
		return errors.New(fmt.Sprintf("Error loading .env file: %s", err))

	}
	return nil
}

func getConfigFromEnv(key string) string {
	return os.Getenv(key)
}
