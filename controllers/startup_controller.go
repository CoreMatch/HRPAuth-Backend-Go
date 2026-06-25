package controllers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = "config.yaml"
const ConfigFileDir = "./"

type StartupController struct{}

func NewStartupController() *StartupController {
	return &StartupController{}
}

func (sc *StartupController) InitializeConfig() error {
	configPath := filepath.Join(ConfigFileDir, ConfigFileName)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found at %s, creating default config...", configPath)
		return sc.createDefaultConfig(configPath)
	}

	log.Printf("Config file found at %s", configPath)
	return nil
}

func (sc *StartupController) createDefaultConfig(path string) error {
	defaultConfig := map[string]interface{}{
		"version": "1.0",
		"site": map[string]interface{}{
			"name":           "HRPAuth",
			"implementation": "HRPAuth zggdrasil-api service",
			"version":        "62526",
		},
		"server": map[string]interface{}{
			"port":        ":2778",
			"cors_origin": "https://auth.samuelcheston.com",
		},
		"callback": map[string]interface{}{
			"url": "https://backend.auth.samuelcheston.com/",
		},
		"frontend": map[string]interface{}{
			"url": "https://auth.samuelcheston.com/",
		},
		"keygen": map[string]interface{}{
			"enable": 0,
		},
		"database": map[string]interface{}{
			"host":     "192.168.1.124",
			"db_name":  "hrpa",
			"user":     "hrpa",
			"password": "hrpa",
			"charset":  "utf8mb4",
		},
		"memcache": map[string]interface{}{
			"host":        "127.0.0.1",
			"port":        11211,
			"prefix":      "hrpauth_",
			"code_ttl":    600,
			"storage_dir": "./cache/verification_codes",
		},
		"smtp": map[string]interface{}{
			"host":       "127.0.0.1",
			"port":       25,
			"username":   "",
			"password":   "",
			"encryption": "tls",
			"from_email": "no-reply@samuelcheston.com",
			"from_name":  "HRPAuth",
		},
		"yggdrasil": map[string]interface{}{
			"server": map[string]interface{}{
				"name":                 "HRPAuth",
				"implementation":       "HRPAuth zggdrasil-api service",
				"version":              "5526",
				"signature_public_key": "",
				"links": map[string]interface{}{
					"homepage": "",
					"register": "",
				},
				"skin_domains": []string{},
			},
			"security": map[string]interface{}{
				"token_expiry_days":      15,
				"session_expiry_seconds": 30,
				"password_cost":          10,
			},
			"feature_flags": map[string]interface{}{
				"non_email_login":             true,
				"legacy_skin_api":             true,
				"no_mojang_namespace":         false,
				"enable_mojang_anti_features": false,
				"enable_profile_key":          false,
				"username_check":              true,
			},
		},
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	log.Printf("Default config file created at %s", path)
	log.Printf("Please edit the configuration file and restart the application")
	return nil
}
