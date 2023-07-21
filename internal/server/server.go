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
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
	"strconv"
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
	router *mux.Router
}

func NewHandler(logger logging.Logger, store *sessions.CookieStore, db *sql.DB, route *mux.Router) handlers.Handler {
	return &handler{
		logger: logger,
		store:  store,
		db:     db,
		router: route,
	}
}

func (h *handler) Register() {
	//тут будут хэндлеры
	h.router.HandleFunc("/signup", h.signupHandler)
	h.router.HandleFunc("/login", h.loginHandler)
	h.router.HandleFunc("/profile", h.profileHandler)
	h.router.HandleFunc("/logout", h.logoutHandler)
	h.router.HandleFunc("/profile/booking", h.bookingHandler)
	h.router.HandleFunc("/profile/booking/{query}", h.bookingQueryHandler)
	h.router.HandleFunc("/profile/booking/{query}/{id}", h.HotelHandler)
	h.router.HandleFunc("/profile/info", h.infoHandler)
	h.router.HandleFunc("/profile/correctinfo", h.correctInfoHandler)

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
	switch r.Method {
	default:
		helper.LoadPage(w, "booking", nil)
		return
	}
}

func (h *handler) bookingQueryHandler(w http.ResponseWriter, r *http.Request) {
	usr := &user.User{}

	not := user.ValidSessionOrNot(usr, h.store, h.logger, w, r, h.db)
	if not != true {
		return
	}
	switch r.Method {
	default:
		srch := &search.Search{
			Request: mux.Vars(r)["query"],
		}

		srch.GetResults(h.db, h.logger)

		helper.LoadPage(w, "booking.query", srch)
	}
}

func (h *handler) HotelHandler(w http.ResponseWriter, r *http.Request) {
	usr := &user.User{}

	not := user.ValidSessionOrNot(usr, h.store, h.logger, w, r, h.db)
	if not != true {
		return
	}
	switch r.Method {
	default:
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		hotel := search.Hotel{Id: id}
		hotel.FillInfo(h.db, h.logger)
		helper.LoadPage(w, "hotel", hotel)
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

//todo при выборе отеля назначать определенный url. где частью url будет id отеля(dinamicHandler/dinamicURL)
