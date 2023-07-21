package search

import (
	"HotelBooking/pkg/logging"
	"database/sql"
)

type Hotel struct {
	Id        int
	HotelName string
	Address   string
}

func (h *Hotel) FillInfo(db *sql.DB, logger logging.Logger) {
	data := `SELECT hotel_name,address FROM hotels WHERE id = $1`

	if err := db.QueryRow(data, h.Id).Scan(&h.HotelName, &h.Address); err != nil {
		logger.Info(err)
		return
	}
}

type Search struct {
	Id      int
	Request string
	Hotels  []Hotel
}

func (s *Search) GetResults(db *sql.DB, logger logging.Logger) {
	data := `SELECT id,hotel_name,address FROM hotels WHERE city_id = (SELECT id FROM cities WHERE city_name = $1)`

	query, err := db.Query(data, s.Request)
	if err != nil {
		logger.Info(err)
		return
	}

	for query.Next() {
		var htl Hotel
		err = query.Scan(&htl.Id, &htl.HotelName, &htl.Address)
		if err != nil {
			logger.Info(err)
			return
		}
		s.Hotels = append(s.Hotels, htl)
	}
	if query.Err() != nil {
		logger.Info(err)
		return
	}
}
