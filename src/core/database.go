package core

import (
	"bytes"
	"context"
	"fmt"
	pk "github.com/Tnze/go-mc/net/packet"
	"go.mongodb.org/mongo-driver/bson"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		if ctx.Err() != nil {
			panic(fmt.Errorf("database connection timed out: %v", ctx.Err()))
		}
	}()

	if client, err := mongo.Connect(ctx, options.Client().ApplyURI(url)); err == nil {
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

func (d *Database) Write(ip string, data pk.FieldEncoder) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	buffer := new(bytes.Buffer)
	_, err = data.WriteTo(buffer)

	entry := Entry{ip, buffer.Bytes()}

	_, err = d.Database("copingheimer").Collection("servers").InsertOne(ctx, entry)
	return
}

func (d *Database) Read(ip string, data pk.FieldDecoder) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var entry Entry

	if err = d.Database("copingheimer").Collection("servers").FindOne(ctx, bson.M{"ip": ip}).Decode(&entry); err != nil {
		return
	}

	buffer := bytes.NewBuffer(entry.Data)
	_, err = data.ReadFrom(buffer)
	return
}

func (d *Database) GetAll() (entries []Entry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := d.Database("copingheimer").Collection("servers").Find(ctx, bson.M{})
	if err != nil {
		return
	}

	if err = cursor.All(ctx, &entries); err != nil {
		return
	}

	return
}
