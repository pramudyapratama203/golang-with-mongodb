package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"` // menggunakkan _id
	Title string `bson:"title" json:"title"`
	Author string `bson:"author" json:"author"`
	ISBN string `bson:"isbn" json:"isbn"`
	Year int `bson:"year" json:"year"`
}
