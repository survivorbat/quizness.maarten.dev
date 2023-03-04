package main

import (
	"github.com/gin-gonic/gin"
	"github.com/survivorbat/qq.maarten.dev/server"
	"log"
	"os"
)

func main() {
	router := gin.Default()

	instance, err := server.NewServer(os.Getenv("DB_CONNECTION_STRING"), os.Getenv("JWT_SECRET"), os.Getenv("AUTH_CREDENTIALS_PATH"), []byte(os.Getenv("AUTH_SECRET")), os.Getenv("AUTH_REDIRECT_URL"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err := instance.Configure(router); err != nil {
		log.Fatalln(err.Error())
	}

	log.Fatalln(router.Run("0.0.0.0:8000"))
}
