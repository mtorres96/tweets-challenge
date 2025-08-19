package main

import (
	"log"
	"os"

	docs "tweetschallenge/docs"
	"tweetschallenge/internal/bootstrap"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if m := os.Getenv("GIN_MODE"); m != "" {
		gin.SetMode(m)
	}
	docs.SwaggerInfo.BasePath = "/"

	r, shutdown, err := bootstrap.BuildHTTPServer()
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown()

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("listening on %s", addr) // opcional para verificar el puerto
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
