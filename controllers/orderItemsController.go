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

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result , err :=orderItemCollection.UpdateOne(
			ctx ,
			filter , 
			bson.D{
				{"$set" , updatedObj},
			},
			&opt,
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

func ItemsByOrder(id string ) (OrderItems []primitive.M , err error) {
	var  ctx , cancel = context.WithTimeout(context.Background() , 100*time.Second )
    // if you want the get the items in a particulare order you want to lookup in the collection  which is order_item collection but uou have to match with the order_id , so you pass the user_id and match it using the match operator that will give us all the different orderitems for that particular order_id 

	// matchstage is a stage where you are able to match a particulare record using a particulare key from your database
	matchStage := bson.D{{"$match" , bson.D{{"orer_id" , id }}}}
	lookupFoodStage := bson.D{{"$lookup" , bson.D{{"from" , "food"} , {"localField" , "food_id"} , {"foreignField" , "food_id"} , {"as" , "food"}}}} 
	// it takes an array and unwind it so mongodb can do operations on it 
	// if we se this => "preserveNullAndEmptyArrays" tofalse it will remove all the empty arrays 
	unwindFoodStage:= bson.D{{"$unwind" , bson.D{{"$path" , "$food"} , {"preserveNullAndEmptyArrays" , true}}}}

	lookupOrderStage := bson.D{{"$lookup" , bson.D{{"from" , "order"} , {"$localField" , "order_id"} , {"foreignField" , "order_id"} , {"as" , "order"}}}}

	// you are giving it a path and say this is the field from my prevoius stage .. in my prevoius stage i have represnted as order.. i want to take it as a path and unwind it 
	unwindOrderStage := bson.D{{"$unwind" , bson.D{{"path" , "$order"} , {"preserveNullAndEmptyArrays" , true}}}}



	lookupTableStage := bson.D{{"$lookup" , bson.D{{"from" , "table"} , {"localField" , "order.table_id"} , {"foreignField"  ,  "table_id"} , {"as" , "table"}}}}


	unwindTableStage := bson.D{{"$unwind" , bson.D{{"$path" , "$table"} , {"presreveNullAndEmptyArrays" , true}}}}



    // projectStage is basicalyy to manage the fields that you are returning to the frintend
	// controls whatever goes to the next stage
	projectStage := bson.D{
		{
			"$project" , bson.D{
				// it means that i dont wnat id to go to the next stage
				{"id" , 0 },
				// i want the amount field to go , what that field be ? it will refereing food.price
				// i want the food.price item and capture that into amount  
				{"amount" , "$food.price"},
				// i want to send total_count to the front end thats why i said 1 
				{"total_count" , 1},
				{"food_name" , "$food.name"},
				{"food_image" , "$food.food_image"},
				{"table_number" , "$table.table_number"},
				{"table_id" , "$table.table_id"},
				{"order_id" , "$order.order_id"},
				{"price" , "$food.price"},
				{"quantity" , 1},
			},
 
		},

     
	
	}

	   // it basically groups all the data based on the particulare criteria or format  
	   groupStage :=bson.D{{"$group" , bson.D{{"_id" , bson.D{{"order_id" , "$order_id"} , {"table_id" , "$table_id"} , {"table_number" , "$table_number"}}} , {"payment_due" , bson.D{{"$sum" , "$amount"}}}, {"total_count" , bson.D{
		"$sum" , 1,
	}}}, 
	{"order_items" ,
	bson.D{{"$push" , "$$ROOT"}} ,
},
	}} 



	projectStage2 := bson.D{
		{"$project" , bson.D{
			{"id" , 0},
			{"payment_due" , 1},
			{"total_count" , 1},
			{"table_number" , "$_id.table_number"},
			{"order_items" , 1 },
		}},
	}

	result , err := orderItemCollection.Aggregate(
		ctx ,
		mongo.Pipeline{
			matchStage,
			lookupFoodStage,
			unwindFoodStage,
			lookupOrderStage,
			unwindOrderStage,
			lookupTableStage,
			unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,	

		},
	)

	 if err != nil {
      panic(err) 
	 }


	 result.All(ctx , &OrderItems); err !=nil {
		 panic(err)
	 }

	 defer cancel()

	 return OrderItems ,err


	    

	




 
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
