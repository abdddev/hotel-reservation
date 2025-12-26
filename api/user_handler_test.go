package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdddev/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	tdb := setup()
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store.User)
	app.Post("/", userHandler.HandlerPostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if user.ID.IsZero() {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}
