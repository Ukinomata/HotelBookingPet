package main

import (
	"HotelBooking/internal/server"
	"HotelBooking/pkg/database"
	"HotelBooking/pkg/logging"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	logger := logging.GetLogger()
	store := sessions.NewCookieStore([]byte("SDJsdnjsDk"))
	db := database.ConnectToDB()

	handler := server.NewHandler(logger, store, db, router)
	handler.Register()

	logger.Info("start server")
	logger.Fatal(http.ListenAndServe("localhost:8080", router))
}
