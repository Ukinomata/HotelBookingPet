package infoUser

import (
	"HotelBooking/pkg/logging"
	"database/sql"
)

type InfoUser struct {
	Id          int
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
	DOB         string `json:"dob"`
	Photo       []byte `json:"photo"`
	PhotoBase64 string
	Status      string
}

func (info *InfoUser) CorrectInfo(db *sql.DB, logger logging.Logger) {

	data := `UPDATE info_about_user SET name = $1, last_name = $2,dob = $3,photo = $4 WHERE user_id = $5`

	if _, err := db.Exec(data, info.Name, info.LastName, info.DOB, info.Photo, info.Id); err != nil {
		logger.Info(err)
		return
	}
}
