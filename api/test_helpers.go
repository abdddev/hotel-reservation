package api

import (
	"context"
	"log"
	"testing"

	"github.com/abdddev/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.TDBURI))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   db.NewMongoHotelStore(client),
			Room:    db.NewMongoRoomStore(client),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
