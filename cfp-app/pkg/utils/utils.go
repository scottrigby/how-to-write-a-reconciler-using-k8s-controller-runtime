package utils

import (
	"net/http"
	"strings"
)

func MakeFileName(ID string) string {
	return strings.Replace(ID, "/", "-", 1)

}

func Error(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}
