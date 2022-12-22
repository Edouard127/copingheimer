package intf

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Entry struct {
	IP   string
	Data []byte
}

type Database struct {
	*mongo.Client
}

func InitDatabase(url string) *Database {
	var db Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if client, err := mongo.Connect(ctx, options.Client().ApplyURI(url)); err != nil {
		panic(err)
	} else {
		// Create the collection if it doesn't exist
		if client.Database("copingheimer").Collection("servers").FindOne(ctx, bson.M{}); err == nil {
			db = Database{client}
		} else {
			if err := client.Database("copingheimer").CreateCollection(ctx, "servers"); err != nil {
				panic(err)
			} else {
				db = Database{client}
			}
		}
	}
	return &db
}

func (d *Database) Write(ip string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := d.Database("copingheimer").Collection("servers").InsertOne(
		ctx,
		bson.M{
			"ip":   ip,
			"data": data,
		}); err != nil {
		return err
	}
	return nil
}

func (d *Database) Read(ip string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result Entry
	if err := d.Database("copingheimer").Collection("servers").FindOne(ctx, bson.M{"ip": ip}).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
