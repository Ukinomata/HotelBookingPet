package user

import (
	"HotelBooking/internal/models/infoUser"
	"HotelBooking/pkg/helper"
	"HotelBooking/pkg/logging"
	"database/sql"
	"encoding/base64"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
	"os"
)

type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

func (u *User) SignUp(db *sql.DB, logger logging.Logger) {
	data := `INSERT INTO users(username, hash_password) VALUES($1,$2)`

	if _, err := db.Exec(data, u.Username, helper.Hashing(u.Password)); err != nil {
		logger.Info(err)
		return
	}

	data = `SELECT id FROM users WHERE username = $1`

	if err := db.QueryRow(data, u.Username).Scan(&u.UserID); err != nil {
		logger.Info(err)
		return
	}

	file, err := os.Open("/Users/artemlukmanov/GolandProjects/HotelBooking/pkg/pages/channels4_profile.jpg")
	if err != nil {
		logger.Info(err)
		return
	}

	photoData, err := io.ReadAll(file)
	if err != nil {
		logger.Info(err)
		return
	}

	data = `INSERT INTO info_about_user(user_id,name,last_name,dob,photo) VALUES ($1,$2,$3,$4,$5)`
	if _, err = db.Exec(data, u.UserID, "default", "default", "1990-01-01", photoData); err != nil {
		logger.Info(err)
		return
	}
}

func (u *User) CorrectData(db *sql.DB, logger logging.Logger) bool {
	usr2 := &User{}

	data := `SELECT username,hash_password FROM users WHERE username = $1`

	if err := db.QueryRow(data, u.Username).Scan(&usr2.Username, &usr2.Password); err != nil {
		logger.Info(err)
		return false
	}

	if helper.Hashing(u.Password) == usr2.Password {
		data = `SELECT id FROM users WHERE username = $1`
		db.QueryRow(data, u.Username).Scan(&u.UserID)
		return true
	}
	return false
}

func (u *User) FillProfile(db *sql.DB, logger logging.Logger) {
	data := `SELECT username,status FROM users WHERE id = $1`

	if err := db.QueryRow(data, u.UserID).Scan(&u.Username, &u.Status); err != nil {
		logger.Info(err)
		return
	}
}

func (u *User) FillProfileByUsername(db *sql.DB, logger logging.Logger) {
	data := `SELECT id,status FROM users WHERE username = $1`

	if err := db.QueryRow(data, u.Username).Scan(&u.UserID, &u.Status); err != nil {
		logger.Info(err)
		return
	}
}

func (u *User) SetNewStatus(db *sql.DB, logger logging.Logger) {
	data := `UPDATE users SET status = $1 WHERE id = $2`

	if _, err := db.Exec(data, u.Status, u.UserID); err != nil {
		logger.Info(err)
		return
	}
}

func (u *User) GetInfoAboutUser(db *sql.DB, logger logging.Logger) (info infoUser.InfoUser) {
	data := `SELECT user_id,name,last_name,dob,photo FROM info_about_user WHERE user_id = $1`

	if err := db.QueryRow(data, u.UserID).Scan(&info.Id, &info.Name, &info.LastName, &info.DOB, &info.Photo); err != nil {
		logger.Info(err)
	}
	info.PhotoBase64 = base64.StdEncoding.EncodeToString(info.Photo)
	info.Status = u.Status
	return info
}

func GetAllUsers(db *sql.DB, logger logging.Logger) (usrs []User) {
	data := `SELECT id,username,status FROM users WHERE status != 'god'`

	query, err := db.Query(data)
	if err != nil {
		logger.Info(err)
		return nil
	}

	defer query.Close()

	for query.Next() {
		var usr User
		err = query.Scan(&usr.UserID, &usr.Username, &usr.Status)
		if err != nil {
			logger.Info(err)
			return nil
		}

		usrs = append(usrs, usr)
	}
	if query.Err() != nil {
		logger.Info(err)
		return nil
	}

	return usrs
}

func ValidSessionOrNot(usr *User, store *sessions.CookieStore, logger logging.Logger, w http.ResponseWriter, r *http.Request, db *sql.DB) bool {
	session, err := store.Get(r, "ukinos")
	if err != nil {
		logger.Info(err)
		return false
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return false
	}

	usr.UserID = userID
	usr.FillProfile(db, logger)
	return true
}
