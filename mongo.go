package samplelib

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Item struct {
	Name  string `json:"name" bson:"name"`
	Total int    `json:"total" bson:"total"`
	Id    string `json:"id" bson:"id"`
}

func AddOne(mdb *mongo.Database, i *Item) (*mongo.InsertOneResult, error) {
	collection := mdb.Collection("item")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, i)
	return result, err
}

func GetOne(mdb *mongo.Database, i *Item, filter interface{}) error {
	//Will automatically create a collection if not available
	collection := mdb.Collection("item")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(i)
	return err
}

func Get(mdb *mongo.Database, filter interface{}) []*Item {
	collection := mdb.Collection("item")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var result []*Item
	for cur.Next(ctx) {
		item := &Item{}
		er := cur.Decode(item)
		if er != nil {
			log.Fatal(er)
		}
		result = append(result, item)
	}
	return result
}

func Update(mdb *mongo.Database, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mdb.Collection("item")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.UpdateMany(ctx, filter, update)
	return result, err
}

func RemoveOne(mdb *mongo.Database, filter interface{}) (*mongo.DeleteResult, error) {
	collection := mdb.Collection("item")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}
