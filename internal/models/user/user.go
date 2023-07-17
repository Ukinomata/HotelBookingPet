package user

import (
	"HotelBooking/pkg/helper"
	"HotelBooking/pkg/logging"
	"database/sql"
	"github.com/gorilla/sessions"
	"net/http"
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

func (u *User) CompareStatus(db *sql.DB, logger logging.Logger) bool {
	return false
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
