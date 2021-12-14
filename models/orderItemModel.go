package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID         primitive.ObjectID `bson:"_id"`
	Quantity   *string            `json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	Unit_price *float64           `json:"unit_price" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	// if i have the food_id stored in the order_item model that means i can access the entire food object all the details that are stored at particulare food record ...
	// in sql we may use join but here we use lookup ..from ordermodel we lookup to foodmodel using this id
	Food_id       *string `json:"food_id" validate:"required"`
	Order_item_id string  `json:"order_item_id"`
	// similarly we have to loo up order model , so we can access all the details of order model , again using the look up function
	Order_id string `json:"order_id" validate:"required"`
}
