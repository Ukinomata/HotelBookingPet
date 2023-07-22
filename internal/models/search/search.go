package search

import (
	"HotelBooking/pkg/logging"
	"database/sql"
	"fmt"
)

type Hotel struct {
	Id          int
	HotelName   string
	Address     string
	FullAddress string
}

func (h *Hotel) FillInfo(db *sql.DB, logger logging.Logger) {
	data := `SELECT hotel_name,address FROM hotels WHERE id = $1`

	if err := db.QueryRow(data, h.Id).Scan(&h.HotelName, &h.Address); err != nil {
		logger.Info(err)
		return
	}

	data = `SELECT countries.country_name, cities.city_name
FROM hotels
JOIN cities ON hotels.city_id = cities.id
JOIN countries ON cities.country_id = countries.id
WHERE hotels.id = $1`

	var (
		country = ""
		city    = ""
	)

	if err := db.QueryRow(data, h.Id).Scan(&country, &city); err != nil {
		logger.Info(err)
		return
	}

	h.FullAddress = fmt.Sprintf("%s ,%s, %s", country, city, h.Address)
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
