package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lnb/HRPAuth-Backend-Go/config"

	"gopkg.in/yaml.v3"
)

type StartupController struct{}

func NewStartupController() *StartupController {
	return &StartupController{}
}

func (sc *StartupController) InitializeConfig() error {
	configPath := filepath.Join(config.ConfigFileDir, config.ConfigFileName)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found at %s, creating default config...", configPath)
		return sc.createDefaultConfig(configPath)
	}

	log.Printf("Config file found at %s", configPath)
	return nil
}

func (sc *StartupController) createDefaultConfig(path string) error {
	cfgDir := filepath.Dir(path)

	publicKeyPath := filepath.Join(cfgDir, "public_key.pem")
	privateKeyPath := filepath.Join(cfgDir, "private_key.pem")

	if err := sc.generateKeyPair(publicKeyPath, privateKeyPath); err != nil {
		log.Printf("Warning: Failed to generate RSA key pair: %v", err)
		log.Printf("Falling back to pseudo-random keys...")
		if err := sc.generatePseudoKeys(publicKeyPath, privateKeyPath); err != nil {
			log.Printf("Warning: Failed to generate pseudo keys: %v", err)
		}
	}

	defaultConfig := map[string]interface{}{
		"version": config.ConfigVersion,
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
				"name":                       "HRPAuth",
				"implementation":             "HRPAuth zggdrasil-api service",
				"version":                    "5526",
				"signature_public_key_path":  publicKeyPath,
				"signature_private_key_path": privateKeyPath,
				"textures_storage":           "./",
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
	log.Printf("Key pair generated at %s and %s", publicKeyPath, privateKeyPath)
	log.Printf("Please edit the configuration file and restart the application")
	return nil
}

func (sc *StartupController) generateKeyPair(publicKeyPath, privateKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate RSA private key: %v", err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key file: %v", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("failed to write public key file: %v", err)
	}

	return nil
}

func (sc *StartupController) generatePseudoKeys(publicKeyPath, privateKeyPath string) error {
	publicPseudo := sc.generateRandomString(512)
	privatePseudo := sc.generateRandomString(1024)

	publicKeyContent := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----\n", publicPseudo)
	privateKeyContent := fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----\n", privatePseudo)

	if err := os.WriteFile(publicKeyPath, []byte(publicKeyContent), 0644); err != nil {
		return fmt.Errorf("failed to write pseudo public key file: %v", err)
	}
	if err := os.WriteFile(privateKeyPath, []byte(privateKeyContent), 0600); err != nil {
		return fmt.Errorf("failed to write pseudo private key file: %v", err)
	}

	return nil
}

func (sc *StartupController) generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
