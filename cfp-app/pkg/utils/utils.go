package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func MakeFileName(ID string) string {
	return strings.Replace(ID, "/", "-", 1)

}

func Error(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func Exists(id string, path string) bool {
	b, _ := os.ReadFile(fmt.Sprintf("%s%s.json", path, MakeFileName(id)))
	return len(b) > 0
}
