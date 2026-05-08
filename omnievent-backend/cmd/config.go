package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

const (
	EnvPrefix = "OMNIEVENT"
)

var cfg *Config

type Config struct {
	Global   GlobalConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Security SecurityConfig
	User     UserConfig
}

type GlobalConfig struct {
	AppName     string
	Mode        string
	Version     string
	WorkingPath string
}

type ServerConfig struct {
	Protocol        string
	HttpAddr        string
	HttpPort        int
	Domain          string
	RootURL         string
	StaticRootPath  string
	EnableGZip      bool
	LogRequest      bool
	RequestIDHeader string
}

type DatabaseConfig struct {
	Host              string
	Port              string
	User              string
	Password          string
	DBName            string
	SSLMode           string
	MaxIdleConns      int
	MaxOpenConns      int
	ConnMaxLifetime   int
	EnableQueryLog    bool
	AutoUpdateDatabase bool
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

type JWTConfig struct {
	Secret          string
	ExpireHours     int
	RefreshDays     int
	MaxTokensPerUser int
}

type SecurityConfig struct {
	SecretKey           string
	TokenExpireSeconds  int
	MaxFailuresPerIpMin int
	EnableRateLimiting  bool
}

type UserConfig struct {
	EnableRegister bool
	MinUsernameLen int
	MaxUsernameLen int
	MinPasswordLen int
	MaxPasswordLen int
}

func InitConfig() error {
	// Load from INI file if exists, then environment variables
	cfg = loadConfiguration()
	return nil
}

func GetConfig() *Config {
	return cfg
}

func GetConfigPath() string {
	paths := []string{
		"conf/omnievent.ini",
		"conf/omnievent.conf",
		"/etc/omnievent/omnievent.ini",
		os.Getenv("OMNIEVENT_CONFIG"),
	}

	for _, path := range paths {
		if path != "" {
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				return path
			}
		}
	}
	return ""
}

func loadConfiguration() *Config {
	configPath := GetConfigPath()

	var file *ini.File
	if configPath != "" {
		var err error
		file, err = ini.Load(configPath)
		if err != nil {
			fmt.Printf("Warning: failed to load config file %s: %v\n", configPath, err)
		}
	}

	cfg := &Config{}

	// Global config
	cfg.Global = GlobalConfig{
		AppName:     getConfigString(file, "global", "app_name", "OmniEvent"),
		Mode:        getConfigString(file, "global", "mode", "release"),
		Version:     "1.0.0",
		WorkingPath: getWorkingPath(),
	}

	// Server config
	cfg.Server = ServerConfig{
		Protocol:        getConfigString(file, "server", "protocol", "http"),
		HttpAddr:        getConfigString(file, "server", "http_addr", "0.0.0.0"),
		HttpPort:        getConfigInt(file, "server", "http_port", 8080),
		Domain:          getConfigString(file, "server", "domain", "localhost"),
		RootURL:         getConfigString(file, "server", "root_url", ""),
		StaticRootPath:  getConfigString(file, "server", "static_root_path", ""),
		EnableGZip:      getConfigBool(file, "server", "enable_gzip", false),
		LogRequest:      getConfigBool(file, "server", "log_request", true),
		RequestIDHeader: getConfigString(file, "server", "request_id_header", "X-Request-ID"),
	}

	// Database config
	cfg.Database = DatabaseConfig{
		Host:              getConfigString(file, "database", "host", "localhost"),
		Port:              getConfigString(file, "database", "port", "5432"),
		User:              getConfigString(file, "database", "user", "postgres"),
		Password:          getConfigString(file, "database", "password", "postgres"),
		DBName:            getConfigString(file, "database", "dbname", "omnievent"),
		SSLMode:           getConfigString(file, "database", "sslmode", "disable"),
		MaxIdleConns:      getConfigInt(file, "database", "max_idle_conns", 10),
		MaxOpenConns:      getConfigInt(file, "database", "max_open_conns", 100),
		ConnMaxLifetime:   getConfigInt(file, "database", "conn_max_lifetime", 3600),
		EnableQueryLog:    getConfigBool(file, "database", "log_query", false),
		AutoUpdateDatabase: getConfigBool(file, "database", "auto_update_database", true),
	}

	// Redis config
	cfg.Redis = RedisConfig{
		Host:     getConfigString(file, "redis", "host", "localhost"),
		Port:     getConfigString(file, "redis", "port", "6379"),
		Password: getConfigString(file, "redis", "password", ""),
		DB:       getConfigInt(file, "redis", "db", 0),
		PoolSize: getConfigInt(file, "redis", "pool_size", 10),
	}

	// JWT config
	cfg.JWT = JWTConfig{
		Secret:          getConfigString(file, "jwt", "secret", ""),
		ExpireHours:     getConfigInt(file, "jwt", "expire_hours", 24),
		RefreshDays:     getConfigInt(file, "jwt", "refresh_days", 7),
		MaxTokensPerUser: getConfigInt(file, "jwt", "max_tokens_per_user", 5),
	}

	// Security config
	cfg.Security = SecurityConfig{
		SecretKey:           getConfigString(file, "security", "secret_key", ""),
		TokenExpireSeconds:  getConfigInt(file, "security", "token_expire_seconds", 86400),
		MaxFailuresPerIpMin: getConfigInt(file, "security", "max_failures_per_ip_per_minute", 60),
		EnableRateLimiting:  getConfigBool(file, "security", "enable_rate_limiting", true),
	}

	// User config
	cfg.User = UserConfig{
		EnableRegister:  getConfigBool(file, "user", "enable_register", true),
		MinUsernameLen: getConfigInt(file, "user", "min_username_length", 3),
		MaxUsernameLen: getConfigInt(file, "user", "max_username_length", 32),
		MinPasswordLen: getConfigInt(file, "user", "min_password_length", 6),
		MaxPasswordLen: getConfigInt(file, "user", "max_password_length", 128),
	}

	// Apply environment variable overrides
	applyEnvOverrides(cfg)

	// Set defaults for empty secrets
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = generateSecretKey(32)
	}
	if cfg.Security.SecretKey == "" {
		cfg.Security.SecretKey = generateSecretKey(32)
	}

	return cfg
}

func applyEnvOverrides(cfg *Config) {
	// JWT Secret
	if secret := os.Getenv(fmt.Sprintf("%s_JWT_SECRET", EnvPrefix)); secret != "" {
		cfg.JWT.Secret = secret
	}
	// JWT Expire Hours
	if hours := os.Getenv(fmt.Sprintf("%s_JWT_EXPIRE_HOURS", EnvPrefix)); hours != "" {
		if h, err := strconv.Atoi(hours); err == nil {
			cfg.JWT.ExpireHours = h
		}
	}
	// Database
	if host := os.Getenv(fmt.Sprintf("%s_DB_HOST", EnvPrefix)); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv(fmt.Sprintf("%s_DB_PORT", EnvPrefix)); port != "" {
		cfg.Database.Port = port
	}
	if user := os.Getenv(fmt.Sprintf("%s_DB_USER", EnvPrefix)); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv(fmt.Sprintf("%s_DB_PASSWORD", EnvPrefix)); password != "" {
		cfg.Database.Password = password
	}
	if dbname := os.Getenv(fmt.Sprintf("%s_DB_NAME", EnvPrefix)); dbname != "" {
		cfg.Database.DBName = dbname
	}
	// Redis
	if host := os.Getenv(fmt.Sprintf("%s_REDIS_HOST", EnvPrefix)); host != "" {
		cfg.Redis.Host = host
	}
	if port := os.Getenv(fmt.Sprintf("%s_REDIS_PORT", EnvPrefix)); port != "" {
		cfg.Redis.Port = port
	}
	// Server
	if port := os.Getenv(fmt.Sprintf("%s_SERVER_PORT", EnvPrefix)); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.HttpPort = p
		}
	}
}

func getConfigString(file *ini.File, section, key, defaultVal string) string {
	envKey := fmt.Sprintf("%s_%s_%s", EnvPrefix, strings.ToUpper(section), strings.ToUpper(key))
	if val := os.Getenv(envKey); val != "" {
		return val
	}

	if file != nil {
		if sec, err := file.GetSection(section); err == nil {
			val := sec.Key(key).String()
			if val != "" {
				return val
			}
		}
	}
	return defaultVal
}

func getConfigInt(file *ini.File, section, key string, defaultVal int) int {
	envKey := fmt.Sprintf("%s_%s_%s", EnvPrefix, strings.ToUpper(section), strings.ToUpper(key))
	if val := os.Getenv(envKey); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}

	if file != nil {
		if sec, err := file.GetSection(section); err == nil {
			if val, err := sec.Key(key).Int(); err == nil {
				return val
			}
		}
	}
	return defaultVal
}

func getConfigBool(file *ini.File, section, key string, defaultVal bool) bool {
	envKey := fmt.Sprintf("%s_%s_%s", EnvPrefix, strings.ToUpper(section), strings.ToUpper(key))
	if val := os.Getenv(envKey); val != "" {
		return strings.ToLower(val) == "true" || val == "1"
	}

	if file != nil {
		if sec, err := file.GetSection(section); err == nil {
			if val, err := sec.Key(key).Bool(); err == nil {
				return val
			}
		}
	}
	return defaultVal
}

func getWorkingPath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

func generateSecretKey(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(time.Nanosecond)
	}
	return string(result)
}
