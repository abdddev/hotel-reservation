package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/abdddev/hotel-reservation/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.Store, "bar hotel", "loc", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)

	from := time.Now()
	till := time.Now().AddDate(0, 0, 2)
	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, till, from)

	fmt.Println("booking:", booking)
}
