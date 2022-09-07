package handlers

import (
	"github.com/Marcel-MD/dining-hall/dto"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Distribution(c *gin.Context) {
	var data dto.Distribution

	if err := c.ShouldBindJSON(&data); err != nil {
		log.Err(err).Msg("Error binding JSON")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int("order_id", data.OrderId).Int("table_id", data.TableId).Int("waiter_id", data.WaiterId).Msg("Serving order")
	c.JSON(200, gin.H{"message": "Order served"})
}
