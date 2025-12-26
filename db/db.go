package db

const DBNAME = "hotel-reservation"
const TDBNAME = "hotel-reservation-test"
const DBURI = "mongodb://localhost:27017/?directConnection=true"
const TDBURI = "mongodb://localhost:27017/?directConnection=true"

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
