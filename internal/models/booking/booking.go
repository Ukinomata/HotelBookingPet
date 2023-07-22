package booking

import (
	"HotelBooking/internal/models/search"
	"HotelBooking/pkg/logging"
	"database/sql"
)

type Booking struct {
	UserID    int
	HotelID   int `json:"hotel_id"`
	HotelInfo search.Hotel
	Enter     string `json:"enter"`
	Out       string `json:"out"`
	Peoples   int    `json:"peoples"`
}

func (b *Booking) BookHotel(db *sql.DB, logger logging.Logger) {
	data := `INSERT INTO booking(user_id, enter_date, out_date, hotel_id,peoples) VALUES ($1,$2,$3,$4,$5)`

	if _, err := db.Exec(data, b.UserID, b.Enter, b.Out, b.HotelID, b.Peoples); err != nil {
		logger.Info(err)
		return
	}
}

func (b *Booking) DeleteReservation(db *sql.DB, logger logging.Logger) {
	data := `DELETE FROM booking WHERE user_id = $1 and hotel_id = $2`

	if _, err := db.Exec(data, b.UserID, b.HotelID); err != nil {
		logger.Info(err)
		return
	}
}
