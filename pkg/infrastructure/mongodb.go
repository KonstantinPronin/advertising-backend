package infrastructure

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Database struct {
	client *mongo.Client
	name   string
}

func (db *Database) GetConnection(opts ...*options.DatabaseOptions) *mongo.Database {
	return db.client.Database(db.name)
}

func (db *Database) Connect() error {
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Minute)
	defer cancel1()

	if err := db.client.Connect(ctx1); err != nil {
		return err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Minute)
	defer cancel2()

	if err := db.client.Ping(ctx2, nil); err != nil {
		return err
	}

	return nil
}

func (db *Database) Disconnect() error {
	return db.client.Disconnect(context.Background())
}

func InitDatabase(path string) (*Database, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("uri")))
	if err != nil {
		return nil, err
	}

	db := &Database{
		client: client,
		name:   viper.GetString("database"),
	}

	err = db.Connect()
	if err != nil {
		return nil, err
	}

	return db, nil
}
