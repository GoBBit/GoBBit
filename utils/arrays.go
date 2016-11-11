package utils

import (
	"gopkg.in/mgo.v2/bson"
)

// Functions to search elements in arrays, etc..
func IndexOf(array []bson.ObjectId, e bson.ObjectId) (int){

	for i, ele := range array{
		if ele == e{
			return i
		}
	}
	return -1

}

