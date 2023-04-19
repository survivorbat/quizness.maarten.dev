package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server"
	_ "github.com/survivorbat/qq.maarten.dev/server/routes" // Import for swaggo
	"github.com/toorop/gin-logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"log"
	"os"
)

const ServiceName = "Quizness-Server"

// @title						QQ
// @BasePath					/
// @securityDefinitions.apikey	JWT
// @in							header
// @name						Authorization
func main() {
	_ = godotenv.Load()

	cancel, err := configureTracing(os.Getenv("TRACING_ENDPOINT"))
	if err != nil {
		logrus.Fatal(err)
	}
	defer cancel()

	router := gin.New()
	router.Use(ginlogrus.Logger(logrus.New()), gin.Recovery(), otelgin.Middleware(ServiceName))

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{os.Getenv("CORS_ALLOW_ORIGIN")},
		AllowMethods:  []string{"GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Authorization"},
		ExposeHeaders: []string{"Content-Length", "Token"},
	}))

	instance, err := server.NewServer(os.Getenv("DB_CONNECTION_STRING"), os.Getenv("JWT_SECRET"), os.Getenv("AUTH_CLIENT_ID"), os.Getenv("AUTH_CLIENT_SECRET"), os.Getenv("AUTH_REDIRECT_URL"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err := instance.Configure(router); err != nil {
		log.Fatalln(err.Error())
	}

	log.Fatalln(router.Run("0.0.0.0:8000"))
}

func configureTracing(url string) (func(), error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		logrus.Error(err)
		return func() {}, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(ServiceName))),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New())

	// Cleanly shutdown and flush telemetry when the application exits.
	shutdown := func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logrus.Fatal(err)
		}
	}

	logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	)))

	return shutdown, nil
}
