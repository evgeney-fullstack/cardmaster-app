package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection represents a group of flashcards
// Users can organize their cards into collections
type Collection struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`                       // Unique collection identifier
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`              // Reference to user who owns this collection
	Name   string             `json:"name" bson:"name" binding:"required"` // Collection display name
}

// UpdateCollectionInput represents data for updating a collection
// Uses pointers to distinguish between zero values and unset fields
type UpdateCollectionInput struct {
	Name *string `json:"name" bson:"name"` // Pointer allows detecting if field was provided
}

// Validate ensures at least one field is provided for update
// Prevents empty updates that don't change anything
func (i UpdateCollectionInput) Validate() error {
	if i.Name == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
