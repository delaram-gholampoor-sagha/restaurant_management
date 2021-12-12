package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Delaram-Gholampoor-Sagha/restaurant_management/database"
	"github.com/Delaram-Gholampoor-Sagha/restaurant_management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(context.TODO(), bson.M{})

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})
			return
		}

		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

// one order has multiple order items , if you want you can have one particular item using this function if you want
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "eeror occured while listing ordered item"})
		}
		c.JSON(http.StatusOK, orderItem)

	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemsToBeInserted := []interface{}{}

		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreated(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = ToFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)

		}

		insertedOrderedItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			log.Fatal(err)
		}

		defer cancel()
		c.JSON(http.StatusOK, insertedOrderedItems)

	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")

		filter := bson.M{"order_item_id": orderItem}

		var updatedObj primitive.D

		if orderItem.Unit_price != nil {
			updatedObj = append(updatedObj, bson.E{"unit_price", *&orderItem.Unit_price})
		}

		if orderItem.Quantity != nil {
			updatedObj = append(updatedObj, bson.E{"quantity", *orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updatedObj = append(updatedObj, bson.E{"food_id", *orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updatedObj = append(updatedObj, bson.E{"updated_at", orderItem.Updated_at})

		upsert := true

		opt := options.UpdateOptions(
			Upsert: &upsert,
		)
		result , err :=orderItemCollection.UpdateOne(
			ctx ,
			filter , 
			bson.D{
				{"$set" , updatedObj},
			},
			&opt
		)

		if err != nil {
			msg := "order item update failed"
			c.JSON(http.StatusInternalServerError , gin.H{"error" : msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK , result)

	}
}

func ItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItem, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items by order ID "})
			return
		}
		c.JSON(http.StatusOK, allOrderItem)
	}
}
