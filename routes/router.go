package routes

import (
	"dairyfortune/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, world!"})
	})

	// 🃏 Draw card
	r.GET("/draw", controllers.DrawCard)

	r.GET("/cards", controllers.GetCards)
	r.PATCH("/cards", controllers.UpdateCard)
	r.POST("/cards", controllers.CreateCard)
	r.DELETE("/cards", controllers.DeleteCard)

	return r
}
