package main

import (
	"os"

	"github.com/Delaram-Gholampoor-Sagha/restaurant_management/database"
	"github.com/Delaram-Gholampoor-Sagha/restaurant_management/middleware"
	"github.com/Delaram-Gholampoor-Sagha/restaurant_management/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	//first we have to check if the user is authenticated or not then would be ready to use these routes
	router.Use(middleware.Authentication())

	// here we are initilizing our routes
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)
}
