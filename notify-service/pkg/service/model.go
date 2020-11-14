package service

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//MessageNotification struct schema
type MessageNotification struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Message   string             `json:"message,omitempty" bson:"message,omitempty" binding:"required"`
	Method    string             `json:"method,omitempty" bson:"method,omitempty" binding:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty" binding:"required"`
}
