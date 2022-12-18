package utils

import (
	"bytes"
	"context"
	"edouard127/copingheimer/src/intf"
	"fmt"
	"github.com/boltdb/bolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Entry struct {
	IP   string
	Data []byte
}

func InitDatabase(cfg *intf.Arguments) (interface{}, error) {
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
		fmt.Println("Connecting to MongoDB...", cfg.DatabaseURL)
		if client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DatabaseURL)); err != nil {
			return nil, err
		} else {
			// Create the collection if it doesn't exist
			if client.Database("copingheimer").Collection("servers").FindOne(ctx, bson.M{}); err == nil {
				return client, nil
			} else {
				if err := client.Database("copingheimer").CreateCollection(ctx, "servers"); err != nil {
					return nil, err
				} else {
					return client, nil
				}
			}
		}
	}
	return nil, nil
}

func Write(db interface{}, ip string, data []byte) error {
	// Check if the IP is in the database
	entries, _ := Find(db, ip)
	if len(entries) > 1 {
		for _, entry := range entries {
			if bytes.Equal(entry, data) {
				if err := Delete(db, ip); err != nil {
					return err
				}
			}
		}
	} else if len(entries) == 1 {
		if bytes.Equal(entries[0], data) {
			return nil
		}
	}
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

func Find(db interface{}, search string) ([][]byte, error) {
	var ips [][]byte
	switch db.(type) {
	case *bolt.DB:
		if err := db.(*bolt.DB).View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("servers"))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if string(v) == search {
					ips = append(ips, k)
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
			for _, result := range results {
				ips = append(ips, (*result)["ip"].([]byte))
			}
			return ips, nil
		}
	}
	return nil, nil
}

func Delete(db interface{}, ip string) error {
	switch db.(type) {
	case *bolt.DB:
		if err := db.(*bolt.DB).Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("servers"))
			err := b.Delete([]byte(ip))
			return err
		}); err != nil {
			return err
		}
	case *mongo.Client:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := db.(*mongo.Client).Database("copingheimer").Collection("servers").DeleteOne(ctx, bson.M{"ip": ip}); err != nil {
			return err
		}
	}
	return nil
}

func Drop(db interface{}) error {
	switch db.(type) {
	case *bolt.DB:
		if err := db.(*bolt.DB).Update(func(tx *bolt.Tx) error {
			err := tx.DeleteBucket([]byte("servers"))
			return err
		}); err != nil {
			return err
		}
	case *mongo.Client:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := db.(*mongo.Client).Database("copingheimer").Collection("servers").DeleteMany(ctx, bson.M{}); err != nil {
			return err
		}
	}
	return nil
}

func CleanDuplicates(db interface{}) (int, error) {
	if all, err := GetAll(db); err != nil {
		return 0, err
	} else {
		if err := Drop(db); err != nil {
			return 0, err
		}
		for _, entry := range all {
			if err := Write(db, entry.IP, entry.Data); err != nil {
				return 0, err
			}
		}
		return len(all), nil
	}
}

func GetAll(db interface{}) ([]Entry, error) {
	all := make([]Entry, 0)
	switch db.(type) {
	case *bolt.DB:
		if err := db.(*bolt.DB).View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("servers"))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				all = append(all, Entry{IP: string(k), Data: v})
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return all, nil
	case *mongo.Client:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		results := make([]*bson.M, 0)
		if cursor, err := db.(*mongo.Client).Database("copingheimer").Collection("servers").Find(ctx, bson.M{}); err != nil {
			return nil, err
		} else {
			if err := cursor.All(ctx, &results); err != nil {
				return nil, err
			} else {
				for _, result := range results {
					if binData, ok := (*result)["data"].(primitive.Binary); ok {
						all = append(all, Entry{IP: (*result)["ip"].(string), Data: binData.Data})
					}
				}
				return all, nil
			}
		}
	}
	return nil, nil
}
