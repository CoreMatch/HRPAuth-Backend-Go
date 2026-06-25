package config

import (
	"os"
	"strconv"
)

type Config struct {
	Site        SiteConfig
	Callback    CallbackConfig
	Frontend    FrontendConfig
	KeyGen      KeyGenConfig
	Database    DatabaseConfig
	Memcache    MemcacheConfig
	SMTP        SMTPConfig
	Yggdrasil   YggdrasilConfig
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
	Server        ServerConfig
	Security      SecurityConfig
	FeatureFlags  FeatureFlagsConfig
}

type ServerConfig struct {
	Name              string
	Implementation    string
	Version           string
	Links             LinksConfig
	SkinDomains       []string
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
	NonEmailLogin         bool
	LegacySkinAPI         bool
	NoMojangNamespace     bool
	EnableMojangAntiFeatures bool
	EnableProfileKey      bool
	UsernameCheck         bool
}

var AppConfig *Config

func Load() {
	AppConfig = &Config{
		Site: SiteConfig{
			Name:           getEnv("SITE_NAME", "HRPAuth"),
			Implementation: getEnv("SITE_IMPLEMENTATION", "HRPAuth zggdrasil-api service"),
			Version:        getEnv("SITE_VERSION", "5526"),
		},
		Callback: CallbackConfig{
			URL: getEnv("CALLBACK_URL", "https://hrpauth.samuelcheston.com/"),
		},
		Frontend: FrontendConfig{
			URL: getEnv("FRONTEND_URL", "https://auth.samuelcheston.com/"),
		},
		KeyGen: KeyGenConfig{
			Enable: getEnvInt("KEYGEN_ENABLE", 0),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			DBName:   getEnv("DB_NAME", "hrpa"),
			User:     getEnv("DB_USER", "hrpa"),
			Password: getEnv("DB_PASS", "hrpa"),
			Charset:  "utf8mb4",
		},
		Memcache: MemcacheConfig{
			Host:       getEnv("MEMCACHE_HOST", "127.0.0.1"),
			Port:       getEnvInt("MEMCACHE_PORT", 11211),
			Prefix:     getEnv("MEMCACHE_PREFIX", "hrpauth_"),
			CodeTTL:    getEnvInt("MEMCACHE_CODE_TTL", 600),
			StorageDir: getEnv("MEMCACHE_STORAGE_DIR", "./cache/verification_codes"),
		},
		SMTP: SMTPConfig{
			Host:       getEnv("SMTP_HOST", "127.0.0.1"),
			Port:       getEnvInt("SMTP_PORT", 25),
			Username:   getEnv("SMTP_USERNAME", ""),
			Password:   getEnv("SMTP_PASSWORD", ""),
			Encryption: getEnv("SMTP_ENCRYPTION", "tls"),
			FromEmail:  getEnv("SMTP_FROM_EMAIL", "no-reply@samuelcheston.com"),
			FromName:   getEnv("SMTP_FROM_NAME", "HRPAuth"),
		},
		Yggdrasil: YggdrasilConfig{
			Server: ServerConfig{
				Name:              getEnv("SITE_NAME", "HRPAuth"),
				Implementation:    getEnv("SITE_IMPLEMENTATION", "HRPAuth zggdrasil-api service"),
				Version:           getEnv("SITE_VERSION", "5526"),
				SignaturePublicKey: "",
			},
			Security: SecurityConfig{
				TokenExpiryDays:      getEnvInt("TOKEN_EXPIRY_DAYS", 15),
				SessionExpirySeconds: getEnvInt("SESSION_EXPIRY_SECONDS", 30),
				PasswordCost:         getEnvInt("PASSWORD_COST", 10),
			},
			FeatureFlags: FeatureFlagsConfig{
				NonEmailLogin:             getEnvBool("FEATURE_NON_EMAIL_LOGIN", true),
				LegacySkinAPI:             getEnvBool("FEATURE_LEGACY_SKIN_API", true),
				NoMojangNamespace:         getEnvBool("FEATURE_NO_MOJANG_NAMESPACE", false),
				EnableMojangAntiFeatures:  getEnvBool("FEATURE_ENABLE_MOJANG_ANTI_FEATURES", false),
				EnableProfileKey:          getEnvBool("FEATURE_ENABLE_PROFILE_KEY", false),
				UsernameCheck:             getEnvBool("FEATURE_USERNAME_CHECK", true),
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
