package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type (
	User struct {
		ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Email    string        `json:"email" bson:"email"`
		Password string        `json:"password,omitempty" bson:"password"`
		CreateAt time.Time     `json:"createdAt" bson:"createdAt"`
		Token    string        `json:"token,omitempty" bson:"-"`
	}

	UserDetail struct {
		ID        bson.ObjectId `json:"id" bson:"_id"`
		UserID    int           `json:"userId" bson:"userId"`
		Firstname string        `json:"firstname" bson:"firstname"`
		Lastname  string        `json:"lastname" bson:"lastname"`
		Age       int           `json:"age" bson:"age"`
	}
)
