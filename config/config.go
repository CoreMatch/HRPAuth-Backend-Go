package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Version   string
	Site      SiteConfig
	Server    ServerRuntimeConfig
	Callback  CallbackConfig
	Frontend  FrontendConfig
	KeyGen    KeyGenConfig
	Database  DatabaseConfig
	Memcache  MemcacheConfig
	SMTP      SMTPConfig
	Yggdrasil YggdrasilConfig
}

type ServerRuntimeConfig struct {
	Port       string
	CORSOrigin string
}

type SiteConfig struct {
	Name           string
	Implementation string
	Version        string
}

type CallbackConfig struct {
	URL string
}

type FrontendConfig struct {
	URL string
}

type KeyGenConfig struct {
	Enable int
}

type DatabaseConfig struct {
	Host     string
	DBName   string
	User     string
	Password string
	Charset  string
}

type MemcacheConfig struct {
	Host       string
	Port       int
	Prefix     string
	CodeTTL    int
	StorageDir string
}

type SMTPConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Encryption string
	FromEmail  string
	FromName   string
}

type YggdrasilConfig struct {
	Server       ServerConfig
	Security     SecurityConfig
	FeatureFlags FeatureFlagsConfig
}

type ServerConfig struct {
	Name               string
	Implementation     string
	Version            string
	Links              LinksConfig
	SkinDomains        []string
	SignaturePublicKey string
}

type LinksConfig struct {
	Homepage string
	Register string
}

type SecurityConfig struct {
	TokenExpiryDays      int
	SessionExpirySeconds int
	PasswordCost         int
}

type FeatureFlagsConfig struct {
	NonEmailLogin            bool
	LegacySkinAPI            bool
	NoMojangNamespace        bool
	EnableMojangAntiFeatures bool
	EnableProfileKey         bool
	UsernameCheck            bool
}

const ConfigFileName = "config.yaml"
const ConfigFileDir = "./"
const ConfigVersion = "1.0"

var AppConfig *Config

func Load() {
	configPath := filepath.Join(ConfigFileDir, ConfigFileName)

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file %s: %v", configPath, err)
	}

	// Parse YAML
	var yamlConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Validate config version
	configVersion := getString(yamlConfig, "version")
	if configVersion == "" {
		log.Fatalf("Config file is missing version field")
	}
	if configVersion != ConfigVersion {
		log.Printf("Warning: Config file version %s does not match expected version %s", configVersion, ConfigVersion)
	}

	// Map YAML to Config struct
	AppConfig = &Config{
		Version:   getString(yamlConfig, "version"),
		Site:      parseSiteConfig(yamlConfig),
		Server:    parseServerRuntimeConfig(yamlConfig),
		Callback:  parseCallbackConfig(yamlConfig),
		Frontend:  parseFrontendConfig(yamlConfig),
		KeyGen:    parseKeyGenConfig(yamlConfig),
		Database:  parseDatabaseConfig(yamlConfig),
		Memcache:  parseMemcacheConfig(yamlConfig),
		SMTP:      parseSMTPConfig(yamlConfig),
		Yggdrasil: parseYggdrasilConfig(yamlConfig),
	}

	log.Println("Configuration loaded successfully")
}

func parseSiteConfig(config map[string]interface{}) SiteConfig {
	site, _ := config["site"].(map[string]interface{})
	return SiteConfig{
		Name:           getString(site, "name"),
		Implementation: getString(site, "implementation"),
		Version:        getString(site, "version"),
	}
}

func parseServerRuntimeConfig(config map[string]interface{}) ServerRuntimeConfig {
	server, _ := config["server"].(map[string]interface{})
	return ServerRuntimeConfig{
		Port:       getString(server, "port"),
		CORSOrigin: getString(server, "cors_origin"),
	}
}

func parseCallbackConfig(config map[string]interface{}) CallbackConfig {
	callback, _ := config["callback"].(map[string]interface{})
	return CallbackConfig{
		URL: getString(callback, "url"),
	}
}

func parseFrontendConfig(config map[string]interface{}) FrontendConfig {
	frontend, _ := config["frontend"].(map[string]interface{})
	return FrontendConfig{
		URL: getString(frontend, "url"),
	}
}

func parseKeyGenConfig(config map[string]interface{}) KeyGenConfig {
	keygen, _ := config["keygen"].(map[string]interface{})
	return KeyGenConfig{
		Enable: getInt(keygen, "enable"),
	}
}

func parseDatabaseConfig(config map[string]interface{}) DatabaseConfig {
	db, _ := config["database"].(map[string]interface{})
	return DatabaseConfig{
		Host:     getString(db, "host"),
		DBName:   getString(db, "db_name"),
		User:     getString(db, "user"),
		Password: getString(db, "password"),
		Charset:  getString(db, "charset"),
	}
}

func parseMemcacheConfig(config map[string]interface{}) MemcacheConfig {
	memcache, _ := config["memcache"].(map[string]interface{})
	return MemcacheConfig{
		Host:       getString(memcache, "host"),
		Port:       getInt(memcache, "port"),
		Prefix:     getString(memcache, "prefix"),
		CodeTTL:    getInt(memcache, "code_ttl"),
		StorageDir: getString(memcache, "storage_dir"),
	}
}

func parseSMTPConfig(config map[string]interface{}) SMTPConfig {
	smtp, _ := config["smtp"].(map[string]interface{})
	return SMTPConfig{
		Host:       getString(smtp, "host"),
		Port:       getInt(smtp, "port"),
		Username:   getString(smtp, "username"),
		Password:   getString(smtp, "password"),
		Encryption: getString(smtp, "encryption"),
		FromEmail:  getString(smtp, "from_email"),
		FromName:   getString(smtp, "from_name"),
	}
}

func parseYggdrasilConfig(config map[string]interface{}) YggdrasilConfig {
	yggdrasil, _ := config["yggdrasil"].(map[string]interface{})
	return YggdrasilConfig{
		Server:       parseServerConfig(yggdrasil),
		Security:     parseSecurityConfig(yggdrasil),
		FeatureFlags: parseFeatureFlagsConfig(yggdrasil),
	}
}

func parseServerConfig(config map[string]interface{}) ServerConfig {
	server, _ := config["server"].(map[string]interface{})
	links, _ := server["links"].(map[string]interface{})
	skinDomains, _ := server["skin_domains"].([]interface{})

	var skinDomainsStr []string
	for _, domain := range skinDomains {
		if str, ok := domain.(string); ok {
			skinDomainsStr = append(skinDomainsStr, str)
		}
	}

	return ServerConfig{
		Name:               getString(server, "name"),
		Implementation:     getString(server, "implementation"),
		Version:            getString(server, "version"),
		SignaturePublicKey: getString(server, "signature_public_key"),
		Links: LinksConfig{
			Homepage: getString(links, "homepage"),
			Register: getString(links, "register"),
		},
		SkinDomains: skinDomainsStr,
	}
}

func parseSecurityConfig(config map[string]interface{}) SecurityConfig {
	security, _ := config["security"].(map[string]interface{})
	return SecurityConfig{
		TokenExpiryDays:      getInt(security, "token_expiry_days"),
		SessionExpirySeconds: getInt(security, "session_expiry_seconds"),
		PasswordCost:         getInt(security, "password_cost"),
	}
}

func parseFeatureFlagsConfig(config map[string]interface{}) FeatureFlagsConfig {
	featureFlags, _ := config["feature_flags"].(map[string]interface{})
	return FeatureFlagsConfig{
		NonEmailLogin:            getBool(featureFlags, "non_email_login"),
		LegacySkinAPI:            getBool(featureFlags, "legacy_skin_api"),
		NoMojangNamespace:        getBool(featureFlags, "no_mojang_namespace"),
		EnableMojangAntiFeatures: getBool(featureFlags, "enable_mojang_anti_features"),
		EnableProfileKey:         getBool(featureFlags, "enable_profile_key"),
		UsernameCheck:            getBool(featureFlags, "username_check"),
	}
}

func getString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if value, ok := m[key].(string); ok {
		return value
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if m == nil {
		return 0
	}
	if value, ok := m[key].(int); ok {
		return value
	}
	if value, ok := m[key].(float64); ok {
		return int(value)
	}
	return 0
}

func getBool(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	if value, ok := m[key].(bool); ok {
		return value
	}
	return false
}
