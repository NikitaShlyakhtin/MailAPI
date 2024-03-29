package main

import (
	"flag"
	"mailapi/internal/jsonlog"
	"mailapi/internal/mailer"
	"sync"

	"github.com/gin-gonic/gin"
)

type config struct {
	port    int
	env     string
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	authKey string
}

type application struct {
	config config
	logger *jsonlog.Logger
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "mode", "debug", "Mode (debug|release)")

	flag.StringVar(&cfg.authKey, "auth-key", "", "Authentication key")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", false, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 0, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "", "SMTP sender")

	flag.Parse()

	if cfg.env == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := jsonlog.New(gin.DefaultWriter, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
