package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout = 10 * time.Second
	queryTimeout      = 5 * time.Second
)

type Database struct {
	Client *mongo.Client
	Name   string
}

func NewDatabase(uri, dbName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &Database{Client: client, Name: dbName}, nil
}

func (db *Database) Close() {
	if db.Client != nil {
		if err := db.Client.Disconnect(context.Background()); err != nil {
			fmt.Println("error while disconnecting from MongoDB:", err)
		}
	}
}

func (db *Database) GetTripsByUserId(userID string) ([]Trip, error) {
	coll := db.Client.Database(db.Name).Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{"client_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get trips by user ID: %w", err)
	}

	var trips []Trip
	if err := cursor.All(ctx, &trips); err != nil {
		return nil, fmt.Errorf("failed to decode trips: %w", err)
	}

	fmt.Println("got trips from the database")
	fmt.Println(trips)

	return trips, nil
}

func (db *Database) GetTripByTripId(tripID string) (*Trip, error) {
	coll := db.Client.Database(db.Name).Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	currID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert trip ID: %w", err)
	}

	var trip Trip
	if err := coll.FindOne(ctx, bson.M{"_id": currID}).Decode(&trip); err != nil {
		return nil, fmt.Errorf("failed to get trip by ID: %w", err)
	}

	fmt.Println("got trip from the database")
	fmt.Println(trip)

	return &trip, nil
}

func (db *Database) CancelTrip(tripID string) error {
	coll := db.Client.Database(db.Name).Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	currID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return fmt.Errorf("failed to convert trip ID: %w", err)
	}

	if _, err := coll.DeleteOne(ctx, bson.M{"_id": currID}); err != nil {
		return fmt.Errorf("failed to cancel trip: %w", err)
	}

	fmt.Println("cancelled trip")

	return nil
}

func (db *Database) CreateTrip(trip *Trip) error {
	coll := db.Client.Database(db.Name).Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if _, err := coll.InsertOne(ctx, trip); err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}

	fmt.Println("trip created")

	return nil
}
