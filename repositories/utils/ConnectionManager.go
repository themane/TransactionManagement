package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	connectTimeoutSecs = 30
)

func GetMongoConnection(mongoURL string) (*mongo.Client, context.Context, context.CancelFunc) {
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeoutSecs*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL).SetMonitor(cmdMonitor))
	if err != nil {
		log.Fatal(err)
		defer cancel()
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		defer cancel()
	}
	log.Println("Successfully connected to MongoDB")
	return client, ctx, cancel
}

func Disconnect(client *mongo.Client, ctx context.Context) {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("Connection to MongoDB closed.")
}
