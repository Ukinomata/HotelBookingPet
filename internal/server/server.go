package server

import (
	"HotelBooking/internal/handlers"
	"HotelBooking/internal/models/infoUser"
	"HotelBooking/internal/models/search"
	"HotelBooking/internal/models/user"
	"HotelBooking/pkg/helper"
	"HotelBooking/pkg/logging"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
)

const (
	userStatus      = "user"
	superUserStatus = "superUser"
	godStatus       = "god"
)

type handler struct {
	logger logging.Logger
	store  *sessions.CookieStore
	db     *sql.DB
}

func NewHandler(logger logging.Logger, store *sessions.CookieStore, db *sql.DB) handlers.Handler {
	return &handler{
		logger: logger,
		store:  store,
		db:     db,
	}
}

func (h *handler) Register() {
	//тут будут хэндлеры
	http.HandleFunc("/signup", h.signupHandler)
	http.HandleFunc("/login", h.loginHandler)
	http.HandleFunc("/profile", h.profileHandler)
	http.HandleFunc("/logout", h.logoutHandler)
	http.HandleFunc("/profile/booking", h.bookingHandler)
	http.HandleFunc("/profile/info", h.infoHandler)
	http.HandleFunc("/profile/correctinfo", h.correctInfoHandler)
}

func (h *handler) signupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		usr := &user.User{
			Username: "",
			Password: "",
		}

		all, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Info(err)
			return
		}

		err = json.Unmarshal(all, &usr)
		if err != nil {
			h.logger.Info(err)
			return
		}

		usr.SignUp(h.db, h.logger)

		return
	default:
		helper.LoadPage(w, "signup", nil)
	}
}

func (h *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		usr := &user.User{}

		all, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Info(err)
			return
		}

		err = json.Unmarshal(all, &usr)
		if err != nil {
			h.logger.Info(err)
			return
		}

		data := usr.CorrectData(h.db, h.logger)
		if data != true {
			return
		}

		session, err := h.store.Get(r, "ukinos")
		if err != nil {
			h.logger.Info(err)
			return
		}

		session.Values["userID"] = usr.UserID
		err = session.Save(r, w)

		if err != nil {
			h.logger.Info(err)
			return
		}
		return
	default:
		helper.LoadPage(w, "login", nil)
	}

}

func (h *handler) profileHandler(w http.ResponseWriter, r *http.Request) {
	usr := &user.User{}
	not := user.ValidSessionOrNot(usr, h.store, h.logger, w, r, h.db)
	if not != true {
		return
	}

	switch usr.Status {
	case godStatus:
		helper.LoadPage(w, "profileGod", usr)
	case superUserStatus:
		helper.LoadPage(w, "profileSuperUser", usr)
	case userStatus:
		helper.LoadPage(w, "profileUser", usr)
	}
}

func (h *handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "ukinos")
	if err != nil {
		h.logger.Info(err)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *handler) bookingHandler(w http.ResponseWriter, r *http.Request) {
	usr := &user.User{}

	not := user.ValidSessionOrNot(usr, h.store, h.logger, w, r, h.db)
	if not != true {
		return
	}
	fmt.Println(usr)
	switch r.Method {
	case http.MethodPost:
		srch := &search.Search{}

		all, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Info(err)
			return
		}

		err = json.Unmarshal(all, &srch)

		if err != nil {
			h.logger.Info(err)
			return
		}

		fmt.Println(srch)
		srch.GetResults(h.db, h.logger)
		srch.GetHotels(h.db, h.logger)
		fmt.Println(srch.Hotels)
		//todo cделать так чтобы после нажатия на кнопку поиск отображали все отели в введенном городе
	default:
		helper.LoadPage(w, "booking", nil)
		return
	}
}

func (h *handler) infoHandler(w http.ResponseWriter, r *http.Request) {
	usr := &user.User{}

	not := user.ValidSessionOrNot(usr, h.store, h.logger, w, r, h.db)
	if not != true {
		return
	}

	info := usr.GetInfoAboutUser(h.db, h.logger)
	switch r.Method {
	default:
		helper.LoadPage(w, "infoUser", info)
	}
}

func (h *handler) correctInfoHandler(w http.ResponseWriter, r *http.Request) {
	usr := &user.User{}
	not := user.ValidSessionOrNot(usr, h.store, h.logger, w, r, h.db)
	if not != true {
		return
	}

	switch r.Method {
	case http.MethodPost:
		inf := &infoUser.InfoUser{Id: usr.UserID}

		file, _, err := r.FormFile("photo")
		if err != nil {
			h.logger.Info(err)
			return
		}

		defer file.Close()

		photoData, err := io.ReadAll(file)

		if err != nil {
			h.logger.Info(err)
			return
		}
		inf.Name = r.FormValue("name")
		inf.LastName = r.FormValue("lastname")
		inf.DOB = r.FormValue("dob")
		inf.Photo = photoData

		inf.CorrectInfo(h.db, h.logger)
		return
	default:
		helper.LoadPage(w, "correctInfoUser", nil)
	}
}

//todo you stop here при нажатии на кнопку отправить данные изменяются а url остается тем же хуй пойми с чем связано, будем фиксить
