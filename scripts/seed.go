package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/abdddev/hotel-reservation/api"
	"github.com/abdddev/hotel-reservation/db"
	"github.com/abdddev/hotel-reservation/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL_TEST")
		dbName        = os.Getenv("MONGO_DB_NAME")
	)

	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(dbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Hotel:   db.NewMongoHotelStore(client),
		Room:    db.NewMongoRoomStore(client),
		Booking: db.NewMongoBookingStore(client),
	}
	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Printf("user -> %s\n", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Printf("admin -> %s\n", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Printf("booking -> %s\n", booking.ID)

	for i := 1; i < 10001; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("random hotel location %d", i)

		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
