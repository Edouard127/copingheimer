package intf

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Entry struct {
	IP   string `bson:"ip"`
	Data []byte `bson:"data"`
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

func (d *Database) Write(ip string, data *StatusResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	players := make([]bson.M, 0)
	for _, player := range data.Players.Sample {
		players = append(players, bson.M{
			"name": player.Name,
			"id":   player.ID,
		})
	}
	if _, err := d.Database("copingheimer").Collection("servers").InsertOne(
		ctx,
		bson.M{
			"ip": ip,
			"data": bson.M{
				"version": bson.M{
					"name":     data.Version.Name,
					"protocol": data.Version.Protocol,
				},
				"players": bson.M{
					"max":    data.Players.Max,
					"online": data.Players.Online,
					"sample": players,
				},
				"favicon":             data.Favicon,
				"preview_chat":        data.PreviewChat,
				"enforce_secure_chat": data.EnforceSecureChat,
			},
		}); err != nil {
		return err
	}
	return nil
}

type reader struct {
	IP   string           `bson:"ip"`
	Data primitive.Binary `bson:"data"`
} // This is a temporary struct to convert the old database to the new one

func (d *Database) ConvertLegacy() error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	if cursor, err := d.Database("copingheimer").Collection("servers").Find(ctx, bson.M{}); err != nil {
		return err
	} else {
		for cursor.Next(ctx) {
			var entry reader
			if err := cursor.Decode(&entry); err != nil {
				fmt.Println("Error decoding entry:", err)
				continue
			}
			var data StatusResponse
			if err := data.Put(entry.Data.Data); err != nil {
				return err
			}
			// Delete the old entry
			if _, err := d.Database("copingheimer").Collection("servers").DeleteMany(ctx, bson.M{"ip": entry.IP}); err != nil {
				fmt.Println("Error deleting old entry:", err)
			}
			// Write the new entry
			if err := d.Write(entry.IP, &data); err != nil {
				fmt.Println("Error writing new entry:", err)
				return err
			}
		}
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
