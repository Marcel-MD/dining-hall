package main

import (
	"os"
	"strconv"

	"github.com/Marcel-MD/dining-hall/domain"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	timeUnit, nrOfTables, nrOfWaiters := config()

	domain.SetTimeUnit(timeUnit)
	menu := domain.GetMenu()
	newOrderChan := make(chan domain.Order)
	ratingChan := make(chan int)
	tablesChans := make([]chan domain.Order, nrOfTables)
	waitersChans := make([]chan domain.Distribution, nrOfWaiters)

	for i := 0; i < nrOfTables; i++ {
		table := domain.NewTable(i, menu, newOrderChan, ratingChan)
		tablesChans = append(tablesChans, table.ReceiveChan)
		go table.Run()
	}

	for i := 0; i < nrOfWaiters; i++ {
		waiter := domain.NewWaiter(i, newOrderChan, tablesChans)
		waitersChans = append(waitersChans, waiter.DistributionChan)
		go waiter.Run()
	}

	go rating(ratingChan)

	r := gin.Default()
	r.POST("/distribution", func(c *gin.Context) {
		var distribution domain.Distribution

		if err := c.ShouldBindJSON(&distribution); err != nil {
			log.Err(err).Msg("Error binding JSON")
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		waiterId := distribution.WaiterId
		waitersChans[waiterId] <- distribution
		c.JSON(200, gin.H{"message": "Order served"})
	})
	r.Run()
}

func rating(ratingChan <-chan int) {
	nrOfRatings := 0
	totalRating := 0

	for {
		rating := <-ratingChan
		nrOfRatings++
		totalRating += rating
		log.Info().Int("rating", rating).Float64("avg_rating", float64(totalRating)/float64(nrOfRatings)).Msg("Received rating")
	}
}

func config() (timeUnit, nrOfTables, nrOfWaiters int) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	timeUnit, err = strconv.Atoi(os.Getenv("TIME_UNIT"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading TIME_UNIT")
	}

	nrOfTables, err = strconv.Atoi(os.Getenv("TABLES"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading TABLES")
	}

	nrOfWaiters, err = strconv.Atoi(os.Getenv("WAITERS"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading WAITERS")
	}

	return
}
