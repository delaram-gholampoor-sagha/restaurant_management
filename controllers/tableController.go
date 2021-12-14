package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Delaram-Gholampoor-Sagha/restaurant_management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := tableCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while lisitng table items"})
			return
		}

		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allTables)

	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		tableId := c.Param("table_id")

		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the tables"})
			return
		}
		c.JSON(http.StatusOK, table)

	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//before we send anything to the database , we validate our struct based on the model that we already created
		validationErr := validate.Struct(table)

		if validationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
			return
		}
		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, insertErr := tableCollection.InsertOne(ctx, table)
		if insertErr != nil {
			msg := fmt.Sprintf("Table item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		tableId := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			log.Fatal(err)
		}

		var updatedObj primitive.D

		if table.Number_of_guests != nil {
			updatedObj = append(updatedObj, bson.E{"number_of_guests": table.Number_of_guests})
		}

		if table.Table_number != nil {
			updatedObj = append(updatedObj, bson.E{"table_number": table.Table_number})
		}

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true

		opt := options.UpdateOptions(
			Upsert : &upsert,
		)

		filter := bson.M{"table_id": tableId}

		result, err := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set": updatedObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("table item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}
