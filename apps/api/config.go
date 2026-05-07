package main

import (
	"errors"
	"os"
	"strings"
)

type config struct {
	Port        string
	DatabaseURL string
	CORSOrigins []string
	Production  bool
}

func loadConfig() (*config, error) {
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return nil, errors.New("DATABASE_URL not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	rawOrigins := os.Getenv("CORS_ORIGINS")
	if rawOrigins == "" {
		rawOrigins = "http://localhost:3000"
	}
	var origins []string
	for _, o := range strings.Split(rawOrigins, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}
	prod := os.Getenv("NODE_ENV") == "production" || os.Getenv("APP_ENV") == "production"
	return &config{
		Port:        port,
		DatabaseURL: dburl,
		CORSOrigins: origins,
		Production:  prod,
	}, nil
}
