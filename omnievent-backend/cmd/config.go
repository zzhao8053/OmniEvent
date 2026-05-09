package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"omnievent-backend/pkg/storage"
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
	Uuid     UuidConfig
	SMTP     SMTPConfig
	Storage  StorageConfig
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	Type                 string
	LocalFileSystemPath  string
	MinIOEndpoint        string
	MinIOLocation        string
	MinIOAccessKeyID     string
	MinIOSecretAccessKey string
	MinIOUseSSL          bool
	MinIOSkipTLSVerify   bool
	MinIOBucket          string
	MinIORootPath        string
	WebDAVURL            string
	WebDAVUsername       string
	WebDAVPassword       string
	WebDAVRootPath       string
	WebDAVRequestTimeout int
	WebDAVProxy          string
	WebDAVSkipTLSVerify  bool
}

type SMTPConfig struct {
	EnableSMTP    bool
	Host          string
	Port          int
	User          string
	Password      string
	FromAddress   string
	SkipTLSVerify bool
}

type UuidConfig struct {
	ServerId uint8
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

	// Uuid config
	cfg.Uuid = UuidConfig{
		ServerId: uint8(getConfigInt(file, "uuid", "server_id", 1)),
	}

	// SMTP config
	cfg.SMTP = SMTPConfig{
		EnableSMTP:    getConfigBool(file, "smtp", "enable_smtp", false),
		Host:          getConfigString(file, "smtp", "host", ""),
		Port:          getConfigInt(file, "smtp", "port", 587),
		User:          getConfigString(file, "smtp", "user", ""),
		Password:      getConfigString(file, "smtp", "password", ""),
		FromAddress:   getConfigString(file, "smtp", "from_address", ""),
		SkipTLSVerify: getConfigBool(file, "smtp", "skip_tls_verify", false),
	}

	// Storage config
	cfg.Storage = StorageConfig{
		Type:                  getConfigString(file, "storage", "type", "local_filesystem"),
		LocalFileSystemPath:   getConfigString(file, "storage", "local_filesystem_path", ""),
		MinIOEndpoint:         getConfigString(file, "storage", "minio_endpoint", ""),
		MinIOLocation:         getConfigString(file, "storage", "minio_location", ""),
		MinIOAccessKeyID:      getConfigString(file, "storage", "minio_access_key_id", ""),
		MinIOSecretAccessKey:  getConfigString(file, "storage", "minio_secret_access_key", ""),
		MinIOUseSSL:           getConfigBool(file, "storage", "minio_use_ssl", false),
		MinIOSkipTLSVerify:    getConfigBool(file, "storage", "minio_skip_tls_verify", false),
		MinIOBucket:           getConfigString(file, "storage", "minio_bucket", ""),
		MinIORootPath:         getConfigString(file, "storage", "minio_root_path", ""),
		WebDAVURL:             getConfigString(file, "storage", "webdav_url", ""),
		WebDAVUsername:        getConfigString(file, "storage", "webdav_username", ""),
		WebDAVPassword:        getConfigString(file, "storage", "webdav_password", ""),
		WebDAVRootPath:        getConfigString(file, "storage", "webdav_root_path", ""),
		WebDAVRequestTimeout:  getConfigInt(file, "storage", "webdav_request_timeout", 10000),
		WebDAVProxy:           getConfigString(file, "storage", "webdav_proxy", "system"),
		WebDAVSkipTLSVerify:   getConfigBool(file, "storage", "webdav_skip_tls_verify", false),
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

	// Initialize storage
	storageConfig := &storage.StorageConfig{
		Type:                cfg.Storage.Type,
		LocalFileSystemPath: cfg.Storage.LocalFileSystemPath,
		MinIOEndpoint:       cfg.Storage.MinIOEndpoint,
		MinIOLocation:        cfg.Storage.MinIOLocation,
		MinIOAccessKeyID:     cfg.Storage.MinIOAccessKeyID,
		MinIOSecretAccessKey: cfg.Storage.MinIOSecretAccessKey,
		MinIOUseSSL:          cfg.Storage.MinIOUseSSL,
		MinIOSkipTLSVerify:   cfg.Storage.MinIOSkipTLSVerify,
		MinIOBucket:          cfg.Storage.MinIOBucket,
		MinIORootPath:        cfg.Storage.MinIORootPath,
		WebDAVURL:            cfg.Storage.WebDAVURL,
		WebDAVUsername:       cfg.Storage.WebDAVUsername,
		WebDAVPassword:       cfg.Storage.WebDAVPassword,
		WebDAVRootPath:       cfg.Storage.WebDAVRootPath,
		WebDAVRequestTimeout: cfg.Storage.WebDAVRequestTimeout,
		WebDAVProxy:          cfg.Storage.WebDAVProxy,
		WebDAVSkipTLSVerify:  cfg.Storage.WebDAVSkipTLSVerify,
	}
	storage.SetStorageConfig(storageConfig)
	if err := storage.InitializeStorageContainer(); err != nil {
		fmt.Printf("Error: failed to initialize storage: %v\n", err)
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
	// Uuid
	if serverId := os.Getenv(fmt.Sprintf("%s_UUID_SERVER_ID", EnvPrefix)); serverId != "" {
		if id, err := strconv.Atoi(serverId); err == nil {
			cfg.Uuid.ServerId = uint8(id)
		}
	}
	// SMTP
	if smtpHost := os.Getenv(fmt.Sprintf("%s_SMTP_HOST", EnvPrefix)); smtpHost != "" {
		cfg.SMTP.Host = smtpHost
	}
	if smtpPort := os.Getenv(fmt.Sprintf("%s_SMTP_PORT", EnvPrefix)); smtpPort != "" {
		if p, err := strconv.Atoi(smtpPort); err == nil {
			cfg.SMTP.Port = p
		}
	}
	if smtpUser := os.Getenv(fmt.Sprintf("%s_SMTP_USER", EnvPrefix)); smtpUser != "" {
		cfg.SMTP.User = smtpUser
	}
	if smtpPassword := os.Getenv(fmt.Sprintf("%s_SMTP_PASSWORD", EnvPrefix)); smtpPassword != "" {
		cfg.SMTP.Password = smtpPassword
	}
	if smtpFrom := os.Getenv(fmt.Sprintf("%s_SMTP_FROM_ADDRESS", EnvPrefix)); smtpFrom != "" {
		cfg.SMTP.FromAddress = smtpFrom
	}
	if smtpEnable := os.Getenv(fmt.Sprintf("%s_SMTP_ENABLE", EnvPrefix)); smtpEnable != "" {
		cfg.SMTP.EnableSMTP = strings.ToLower(smtpEnable) == "true" || smtpEnable == "1"
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
