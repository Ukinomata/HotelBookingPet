package booking

import (
	"HotelBooking/pkg/logging"
	"database/sql"
)

type Reservations struct {
	UserID   int
	Reserves []Booking
}

func (r *Reservations) GetReservations(db *sql.DB, logger logging.Logger) {
	data := `SELECT enter_date,out_date,peoples,hotel_id FROM booking WHERE user_id = $1`

	query, err := db.Query(data, r.UserID)
	if err != nil {
		logger.Info(err)
		return
	}

	for query.Next() {
		b := Booking{}

		err = query.Scan(&b.Enter, &b.Out, &b.Peoples, &b.HotelID)
		if err != nil {
			logger.Info(err)
			return
		}
		b.HotelInfo.Id = b.HotelID
		b.HotelInfo.FillInfo(db, logger)
		r.Reserves = append(r.Reserves, b)
	}
	if query.Err() != nil {
		logger.Info(query.Err())
		return
	}
}
