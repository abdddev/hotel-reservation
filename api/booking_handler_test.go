package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abdddev/hotel-reservation/db/fixtures"
	"github.com/abdddev/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	db := setup()
	defer db.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "Jimmy", "water cooler", false)
		user           = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.Store, "bar hotel", "loc", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 2)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, till, from)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		route          = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
	}

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup()
	defer db.teardown(t)

	var (
		adminUser = fixtures.AddUser(db.Store, "admin", "admin", true)
		user      = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(db.Store, "bar hotel", "loc", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from      = time.Now()
		till      = time.Now().AddDate(0, 0, 2)
		booking   = fixtures.AddBooking(db.Store, user.ID, room.ID, till, from)
		app       = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin     = app.Group("/", JWTAuthentication(db.User), AdminAuth)
	)

	_ = booking
	bookingHandler := NewBookingHandler(db.Store)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", have.ID, booking.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", have.UserID, booking.UserID)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized but got %d", resp.StatusCode)
	}
}
