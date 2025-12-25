package main

import (
	"context"
	"flag"
	"log"

	"github.com/abdddev/hotel-reservation/api"
	"github.com/abdddev/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", api.AdminAuth)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthentication)

	// user handlers
	apiv1.Post("/user", userHandler.HandlerPostUser)
	apiv1.Put("/user/:id", userHandler.HandlerPutUser)
	apiv1.Get("/user", userHandler.HandlerGetUsers)
	apiv1.Get("/user/:id", userHandler.HandlerGetUser)
	apiv1.Delete("/user/:id", userHandler.HandlerDeleteUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandlerGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandlerGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandlerGetRooms)

	// room handlers
	apiv1.Get("/room", roomHandler.HandlerGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandlerBookRoom)

	// bookings handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	app.Listen(*listenAddr)
}
