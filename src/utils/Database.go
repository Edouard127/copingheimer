package utils

import (
	"context"
	"edouard127/copingheimer/src/intf"
	"fmt"
	"github.com/boltdb/bolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InitDatabase(cfg *intf.Config) (interface{}, error) {
	switch cfg.Database {
	case "bolt":
		if Database, err := bolt.Open("copingheimer.db", 0600, nil); err != nil {
			return nil, err
		} else {
			if err := Database.Update(func(tx *bolt.Tx) error {
				if _, err := tx.CreateBucketIfNotExists([]byte("servers")); err != nil {
					return fmt.Errorf("create bucket: %s", err)
				} else {
					return nil
				}
			}); err != nil {
				return nil, err
			}
			return Database, nil
		}
	case "mongodb":
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DatabaseURL)); err != nil {
			return nil, err
		} else {
			// Create the collection if it doesn't exist
			if err := client.Database("copingheimer").CreateCollection(ctx, "servers"); err != nil {
				return nil, err
			} else {
				return client, nil
			}
		}
	}
	return nil, nil
}

func Write(db interface{}, ip string, data []byte) error {
	switch db.(type) {
	case *bolt.DB:
		if err := db.(*bolt.DB).Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("servers"))
			err := b.Put([]byte(ip), data)
			return err
		}); err != nil {
			return err
		}
	case *mongo.Client:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := db.(*mongo.Client).Database("copingheimer").Collection("servers").InsertOne(ctx, bson.D{
			{"ip", ip},
			{"data", data},
		}); err != nil {
			return err
		}
	}
	return nil
}

func Read(db interface{}, ip string) ([]byte, error) {
	switch db.(type) {
	case *bolt.DB:
		var data []byte
		if err := db.(*bolt.DB).View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("servers"))
			data = b.Get([]byte(ip))
			return nil
		}); err != nil {
			return nil, err
		}
		return data, nil
	case *mongo.Client:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var result bson.M
		if err := db.(*mongo.Client).Database("copingheimer").Collection("servers").FindOne(ctx, bson.M{"ip": ip}).Decode(&result); err != nil {
			return nil, err
		}
		return result["data"].([]byte), nil
	}
	return nil, nil
}

func Find(db interface{}, search string) ([]string, error) {
	switch db.(type) {
	case *bolt.DB:
		var ips []string
		if err := db.(*bolt.DB).View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("servers"))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if string(v) == search {
					ips = append(ips, string(k))
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return ips, nil
	case *mongo.Client:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var results []*bson.M
		if cursor, err := db.(*mongo.Client).Database("copingheimer").Collection("servers").Find(ctx, bson.M{search: bson.M{"$exists": true}}); err != nil {
			return nil, err
		} else {
			if err := cursor.All(ctx, &results); err != nil {
				return nil, err
			}
			var ips []string
			for _, result := range results {
				ips = append(ips, (*result)["ip"].(string))
			}
			return ips, nil
		}
	}
	return nil, nil
}

func GetAll(db interface{}) (map[string][]byte, error) {
	switch db.(type) {
	case *bolt.DB:
		var data map[string][]byte
		if err := db.(*bolt.DB).View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("servers"))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				data[string(k)] = v
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return data, nil
	case *mongo.Client:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var results []*bson.M
		if cursor, err := db.(*mongo.Client).Database("copingheimer").Collection("servers").Find(ctx, bson.M{}); err != nil {
			return nil, err
		} else {
			if err := cursor.All(ctx, &results); err != nil {
				return nil, err
			} else {
				data := make(map[string][]byte)
				for _, result := range results {
					data[(*result)["ip"].(string)] = (*result)["data"].([]byte)
				}
				return data, nil
			}
		}
	}
	return nil, nil
}
