package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdddev/hotel-reservation/api"
	"github.com/abdddev/hotel-reservation/db"
	"github.com/abdddev/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var (
//	client       *mongo.Client
//	roomStore    db.RoomStore
//	hotelStore   db.HotelStore
//	userStore    db.UserStore
//	bookingStore db.BookingStore
//	ctx          = context.Background()
//)
//
//func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
//	user, err := types.NewUserFromParams(types.CreateUserParams{
//		FirstName: fname,
//		LastName:  lname,
//		Email:     email,
//		Password:  password,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	user.IsAdmin = isAdmin
//	insertedUser, err := userStore.InsertUser(ctx, user)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
//	return insertedUser
//}
//
//func seedRoom(size string, ss bool, price float64, hotelId primitive.ObjectID) *types.Room {
//	room := &types.Room{
//		Size:    size,
//		Seaside: ss,
//		Price:   price,
//		HotelID: hotelId,
//	}
//	insertedRoom, err := roomStore.InsertRoom(ctx, room)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return insertedRoom
//}
//
//func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
//	booking := &types.Booking{
//		UserID:   userID,
//		RoomID:   roomID,
//		FromDate: from,
//		TillDate: till,
//	}
//	resp, err := bookingStore.InsertBooking(ctx, booking)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("booking -> %s\n", resp.ID)
//}
//
//func seedHotel(name, location string, rating int) *types.Hotel {
//	hotel := types.Hotel{
//		Name:     name,
//		Location: location,
//		Rating:   rating,
//		Rooms:    []primitive.ObjectID{},
//	}
//	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return insertedHotel
//}

func main() {
	ctx := context.Background()
	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
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
	fmt.Printf("booling -> %s\n", booking.ID)

	//james := seedUser(false, "james", "foo", "james@foo.com", "supersecurepassword")
	//seedUser(true, "admin", "admin", "admin@foo.com", "admin")
	//
	//hotel1 := seedHotel("Hotel1", "Location1", 3)
	//seedRoom("small", true, 89.99, hotel1.ID)
	//seedRoom("medium", true, 189.99, hotel1.ID)
	//seedRoom("large", false, 289.99, hotel1.ID)
	//
	//hotel2 := seedHotel("Hotel2", "Location2", 4)
	//seedRoom("small", true, 89.99, hotel2.ID)
	//seedRoom("medium", true, 189.99, hotel2.ID)
	//seedRoom("large", false, 289.99, hotel2.ID)
	//
	//hotel3 := seedHotel("Hotel3", "Location3", 1)
	//seedRoom("small", true, 89.99, hotel3.ID)
	//seedRoom("medium", true, 189.99, hotel3.ID)
	//roomForBooking := seedRoom("large", false, 289.99, hotel3.ID)
	//
	//seedBooking(james.ID, roomForBooking.ID, time.Now(), time.Now().AddDate(0, 0, 2))
}

//func init() {
//	ctx := context.Background()
//	var err error
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
//		log.Fatal(err)
//	}
//	hotelStore = db.NewMongoHotelStore(client)
//	roomStore = db.NewMongoRoomStore(client)
//	userStore = db.NewMongoUserStore(client)
//	bookingStore = db.NewMongoBookingStore(client)
//}
