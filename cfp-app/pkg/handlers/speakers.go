package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/scottrigby/cfp/pkg/types"
)

func CreateSpeaker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var speaker types.Speaker
	json.NewDecoder(r.Body).Decode(&speaker)
	file, _ := json.MarshalIndent(speaker, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("data/speakers/%v.json", speaker.ID), file, 0644)
	json.NewEncoder(w).Encode(speaker)
}

func GetSpeakers(w http.ResponseWriter, r *http.Request) {}

func GetSpeakerById(w http.ResponseWriter, r *http.Request) {}

func UpdateSpeaker(w http.ResponseWriter, r *http.Request) {}

func DeleteSpeaker(w http.ResponseWriter, r *http.Request) {}
