package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	HTTPAddr    string `envconfig:"HTTP_ADDR" default:":8080"`
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
	AdminToken  string `envconfig:"ADMIN_TOKEN" required:"true"`
	CORSOrigin  string `envconfig:"CORS_ORIGIN" default:"http://localhost:5173"`
}

func Load() (Config, error) { var c Config; err := envconfig.Process("APP", &c); return c, err }
