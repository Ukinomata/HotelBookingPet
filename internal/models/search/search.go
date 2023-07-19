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

type Search struct {
	Id      int
	Request string `json:"search"`
	Hotels  []Hotel
}

func (s *Search) GetResults(db *sql.DB, logger logging.Logger) {
	data := `SELECT id FROM cities WHERE city_name = $1`

	if err := db.QueryRow(data, s.Request).Scan(&s.Id); err != nil {
		logger.Info(err)
		return
	}
}

func (s *Search) GetHotels(db *sql.DB, logger logging.Logger) {
	data := `SELECT id,hotel_name,address FROM hotels WHERE city_id = $1`

	query, err := db.Query(data, s.Id)
	if err != nil {
		logger.Info(err)
		return
	}

	defer query.Close()

	for query.Next() {
		htl := Hotel{}

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
