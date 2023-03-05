package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/survivorbat/qq.maarten.dev/server"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000", "https://qq.maarten.dev"},
		AllowMethods:  []string{"GET", "PUT", "PATCH", "Delete"},
		ExposeHeaders: []string{"Content-Length", "token"},
	}))

	instance, err := server.NewServer(os.Getenv("DB_CONNECTION_STRING"), os.Getenv("JWT_SECRET"), os.Getenv("AUTH_CLIENT_ID"), os.Getenv("AUTH_SECRET"), os.Getenv("AUTH_REDIRECT_URL"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err := instance.Configure(router); err != nil {
		log.Fatalln(err.Error())
	}

	log.Fatalln(router.Run("0.0.0.0:8000"))
}
