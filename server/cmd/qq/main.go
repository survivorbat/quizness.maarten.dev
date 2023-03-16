package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server"
	_ "github.com/survivorbat/qq.maarten.dev/server/routes" // Import for swaggo
	"github.com/toorop/gin-logrus"
	"log"
	"os"
)

//	@title						QQ
//	@BasePath					/
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
func main() {
	_ = godotenv.Load()

	router := gin.New()
	router.Use(ginlogrus.Logger(logrus.New()), gin.Recovery())

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
