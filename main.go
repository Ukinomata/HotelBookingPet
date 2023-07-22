package main

import (
	"HotelBooking/internal/models/booking"
	"HotelBooking/pkg/database"
	"HotelBooking/pkg/logging"
	"fmt"
)

func main() {
	db := database.ConnectToDB()
	logger := logging.GetLogger()

	b := booking.Reservations{
		UserID:   1,
		Reserves: nil,
	}

	b.GetReservations(db, logger)

	fmt.Println(b)
}
