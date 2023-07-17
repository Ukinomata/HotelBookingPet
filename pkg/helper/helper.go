package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
)

func Hashing(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}

func LoadPage(w http.ResponseWriter, str string, data any) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("/Users/artemlukmanov/GolandProjects/HotelBooking/pkg/pages/%s.html", str))
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
