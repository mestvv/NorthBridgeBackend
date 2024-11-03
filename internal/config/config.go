package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string      `yaml:"env" env-required:"true"`
	HttpServer HttpServer  `yaml:"http_server" env-required:"true"`
	Database   Database    `yaml:"database" env-required:"true"`
	Limiter    Limiter     `yaml:"limiter" env-required:"true"`
	Auth       AuthConfig  `yaml:"auth" env-required:"true"`
	SMTP       SMTPConfig  `yaml:"smtp" env-required:"true"`
	Email      EmailConfig `yaml:"email" env-required:"true"`
}

type HttpServer struct {
	Port           string        `yaml:"port" env-default:"8080"`
	Timeout        time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout    time.Duration `yaml:"iddle_timeout" env-default:"60s"`
	SwaggerEnabled bool          `yaml:"swagger_enabled" env-default:"false"`
}

type Database struct {
	Net                string        `yaml:"net" env-default:"tcp"`
	Server             string        `yaml:"server" env-required:"true"`
	DBName             string        `yaml:"db_name" env-required:"true"`
	User               string        `yaml:"user" env:"mysql_user" env-required:"true"`
	Password           string        `yaml:"password" env:"mysql_password" env-required:"true"`
	TimeZone           string        `yaml:"time_zone"`
	Timeout            time.Duration `yaml:"timeout" env-default:"2s"`
	MaxIdleConnections int           `yaml:"max_idle_connections" env-default:"40"`
	MaxOpenConnections int           `yaml:"max_open_connections" env-default:"40"`
}

type Limiter struct {
	RPS   int           `yaml:"rps" env-default:"10"`
	Burst int           `yaml:"burst" env-default:"20"`
	TTL   time.Duration `yaml:"ttl" env-default:"10m"`
}

type AuthConfig struct {
	JWT                    JWTConfig `yaml:"jwt" env-required:"true"`
	PasswordSalt           string    `yaml:"password_salt" env-required:"true"`
	VerificationCodeLength int       `yaml:"verification_code_length" env-default:"6"`
}

type JWTConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env-default:"1m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-default:"240h"`
	SigningKey      string        `yaml:"signing_key" env-required:"true"`
}

type SMTPConfig struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
	From string `yaml:"from" env-required:"true"`
	Pass string `yaml:"pass" env-required:"true"`
}

type EmailConfig struct {
	Enabled   bool           `yaml:"enabled" env-default:"false"`
	Templates EmailTemplates `yaml:"templates" env-required:"true"`
}

type EmailTemplates struct {
	Verification string `yaml:"verification" env-required:"true"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
	}

	return &cfg
}
