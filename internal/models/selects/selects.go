package selects

import (
	"HotelBooking/pkg/logging"
	"database/sql"
)

type Country struct {
	ID           int `json:"country_id"`
	CountryName  string
	AllCountries []struct {
		ID          int
		CountryName string
	}
	Cities []City
}

func (country *Country) GetCountries(db *sql.DB, logger logging.Logger) {
	data := `SELECT id,country_name FROM countries`

	query, err := db.Query(data)
	if err != nil {
		logger.Info(err)
		return
	}

	defer query.Close()
	for query.Next() {
		cnrt := struct {
			ID          int
			CountryName string
		}{}

		err = query.Scan(&cnrt.ID, &cnrt.CountryName)
		if err != nil {
			logger.Info(err)
			return
		}

		country.AllCountries = append(country.AllCountries, cnrt)
	}
	if query.Err() != nil {
		logger.Info(query.Err())
		return
	}
}

func (country *Country) GetCities(db *sql.DB, logger logging.Logger) {
	data := `SELECT id,country_id,city_name FROM cities`

	query, err := db.Query(data)
	if err != nil {
		return
	}

	defer query.Close()
	for query.Next() {
		city := City{}

		err = query.Scan(&city.ID, &city.CountryID, &city.CityName)
		if err != nil {
			logger.Info(err)
			return
		}

		country.Cities = append(country.Cities, city)
	}
	if query.Err() != nil {
		logger.Info(query.Err())
		return
	}
}

type City struct {
	ID        int
	CountryID int
	CityName  string
}
