package api

import (
	"github.com/abdddev/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandlerGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("room")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandlerGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(hotel)
}

type ResourseResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

func (h *HotelHandler) HandlerGetHotels(c *fiber.Ctx) error {
	var pagination db.Pagination
	if err := c.QueryParser(&pagination); err != nil {
		return err
	}

	rating := c.QueryInt("rating", 0)
	filter := bson.M{}
	if rating > 0 {
		filter["rating"] = rating
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &pagination)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}

	resp := ResourseResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(pagination.Page),
	}
	return c.JSON(resp)
}
