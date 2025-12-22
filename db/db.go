package db

const DBNAME = "hotel-reservation"
const TDBNAME = "hotel-reservation-test"
const DBURI = "mongodb://localhost:27017/?directConnection=true"

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
