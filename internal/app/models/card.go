package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Card represents a flashcard with spaced repetition data
// Contains question/answer pair and repetition tracking
type Card struct {
	Id              primitive.ObjectID `json:"id" bson:"_id"`                               // Unique card identifier
	UserId          primitive.ObjectID `json:"user_id" bson:"user_id"`                      // Reference to user who owns this card
	CollectionId    primitive.ObjectID `json:"collection_id" bson:"collection_id"`          // Reference to collection containing this card
	Question        string             `json:"question" bson:"question" binding:"required"` // Front side of the flashcard
	Answer          string             `json:"answer" bson:"answer" binding:"required"`     // Back side of the flashcard
	RepetitionStage int                `json:"repetition_stage" bson:"repetition_stage"`    // Current stage in spaced repetition algorithm
	NextRepeatTime  time.Time          `json:"next_repeat_time" bson:"next_repeat_time"`    // When this card should be reviewed next
}

// UpdateCardInput represents data for updating a card
// Uses pointers to distinguish between zero values and unset fields
type UpdateCardInput struct {
	Question        *string    `json:"question" bson:"question"`                 // New question text (optional)
	Answer          *string    `json:"answer" bson:"answer"`                     // New answer text (optional)
	RepetitionStage *int       `json:"repetition_stage" bson:"repetition_stage"` // New repetition stage (optional)
	NextRepeatTime  *time.Time `json:"next_repeat_time" bson:"next_repeat_time"` // New next repetition time (optional)
}

// Validate ensures at least one field is provided for update
// Prevents empty updates that don't change anything
func (i UpdateCardInput) Validate() error {
	if i.Question == nil && i.Answer == nil && i.RepetitionStage == nil && i.NextRepeatTime == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
