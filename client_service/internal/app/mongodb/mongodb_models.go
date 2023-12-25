package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type Geolocation struct {
	Lat float64 `bson:"lat"`
	Lng float64 `bson:"lng"`
}

type Price struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}

type Trip struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	OfferId  string             `bson:"offer_id"`
	From     Geolocation        `bson:"from"`
	To       Geolocation        `bson:"to"`
	ClientId string             `bson:"client_id"`
	Price    Price              `bson:"price"`
	Status   string             `bson:"status"`
}
