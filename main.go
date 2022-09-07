package main

import (
	"fmt"
	"os"

	"github.com/Marcel-MD/dining-hall/domain"
	"github.com/Marcel-MD/dining-hall/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	menu := domain.GetMenu()
	fmt.Println(menu)

	r := gin.Default()

	r.POST("/distribution", handlers.Distribution)

	r.Run()
}
