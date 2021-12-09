package routes

import (
	controller "github.com/Delaram-Gholampoor-Sagha/restaurant_management/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods", controller.GetFoods())
	incomingRoutes.GET("/foods/:food_id", controller.GetFoods())
	// to create a food item
	incomingRoutes.POST("/foods", controller.CreateFood())
	// to update a food item
	incomingRoutes.PATCH("/foods/:food_id", controller.UpdaeFood())
}
