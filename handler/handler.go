package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/fabmation-gmbh/briefkasten-go/handler/apiv1"
	"github.com/fabmation-gmbh/briefkasten-go/handler/ftracer"
	"github.com/fabmation-gmbh/briefkasten-go/handler/middleware"
	"github.com/fabmation-gmbh/briefkasten-go/internal/config"
	"github.com/fabmation-gmbh/briefkasten-go/internal/log"
	"github.com/fabmation-gmbh/briefkasten-go/models"
)

var (
	app    *fiber.App
	tracer trace.Tracer
)

// StartServer will create a new mux router and listen on the configured port.
func StartServer() error {
	tracer = otel.Tracer("server")
	log.Info("Starting server...")

	// establish database connection
	log.Info("Connect to database")
	models.Connect()

	log.Debug("Initializing router")

	decoder := json.Unmarshal
	encoder := json.Marshal

	// configure the fiber behavior
	appSettings := fiber.Config{
		// StrictRouting is very important to be able to use such routes:
		//      * /customer/:id
		//      * /customer/-/:iid
		StrictRouting:           true,
		CaseSensitive:           true,
		ErrorHandler:            ErrorHandler,
		EnableTrustedProxyCheck: true,
		Prefork:                 false,
		TrustedProxies:          []string{"127.0.0.1", "127.0.0.6"},
		JSONEncoder:             encoder,
		JSONDecoder:             decoder,
	}

	app = fiber.New(appSettings)

	registerMiddlewares(app)

	// add all sub-handlers
	log.Debug("Add all sub-handlers for path prefixes")

	// ========== API ==========
	apiv1.AddApiV1(app.Group("/api/v1"))
	// app.Get("/health", apiv1.HealthHandler(rdb))
	// middleware.RegisterAnonymousRoute("/health")

	log.Info("Started successfully!")

	return app.Listen(config.C.General.Listen)
}

const defaultFormat = "${time} [${request_id}] ${method} ${path} - ${ip} (${forwarded_for}) - ${status} - ${latency}\n"

// registerMiddlewares registers all middlewares.
func registerMiddlewares(app *fiber.App) {
	app.Use(ftracer.New(ftracer.Config{
		Tracer: tracer,
	}))

	// the Request ID middleware must be executed before all the other middlewares.
	app.Use(middleware.NewRequestID)

	app.Use(logger.New(logger.Config{
		Format:     defaultFormat,
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Europe/Berlin",
	}))

	if config.C.General.EnableCompression {
		app.Use(compress.New(compress.Config{
			Level: compress.Level(config.C.General.CompressionLevel),
		}))
	}

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
}

func untilError(f ...func() error) error {
	for _, fc := range f {
		if err := fc(); err != nil {
			return err
		}
	}

	return nil
}
