package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/scottrigby/cfp/pkg/types"
	"github.com/scottrigby/cfp/pkg/utils"
)

const speakerDataPath = "data/speakers/"

// CreateSpeaker creates a new file with json data about a given Speaker.
func CreateSpeaker(w http.ResponseWriter, r *http.Request) {
	var speaker types.Speaker
	json.NewDecoder(r.Body).Decode(&speaker)

	if err := validateSpeaker(&speaker); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if speakerExists(speaker.ID) {
		utils.Error(w, fmt.Sprintf("speaker with ID '%s' already exists", speaker.ID), http.StatusBadRequest)
		return
	}

	writeSpeaker(w, r, &speaker)
}

// GetSpeakerById returns the data for a Speaker given the Speaker's ID.
func GetSpeakerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	if id == "" {
		utils.Error(w, "speaker ID must be specified", http.StatusBadRequest)
		return
	}

	b, err := os.ReadFile(fmt.Sprintf("%s%s.json", speakerDataPath, utils.MakeFileName(id)))

	switch {
	case len(b) == 0:
		utils.Error(w, fmt.Sprintf("failed to find Speaker with ID '%s'", id), http.StatusNotFound)
		return
	case err != nil:
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var speaker types.Speaker
	if err := json.Unmarshal(b, &speaker); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(speaker)
}

// GetSpeakerById returns a list with the data of all the files of Speakers
// contained in "data/speakers/"".
func GetSpeakers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	files, err := os.ReadDir(speakerDataPath)
	if err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var speakerList []types.Speaker
	for _, file := range files {
		b, err := os.ReadFile(fmt.Sprintf("%s%s", speakerDataPath, file.Name()))

		switch {
		case len(b) == 0:
			continue
		case err != nil:
			utils.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var speaker types.Speaker
		if err := json.Unmarshal(b, &speaker); err != nil {
			utils.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		speakerList = append(speakerList, speaker)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(speakerList)
}

// UpdateSpeaker checks that a file for a Speaker exists given their ID
// then replaces the json data for that Speaker by overwriting it.
func UpdateSpeaker(w http.ResponseWriter, r *http.Request) {
	var speaker types.Speaker
	json.NewDecoder(r.Body).Decode(&speaker)
	id := mux.Vars(r)["id"]

	if utils.MakeFileName(speaker.ID) != id {
		utils.Error(w, fmt.Sprintf("ID '%s' used as query param does not match ID in request body '%s'", id, speaker.ID), http.StatusBadRequest)
		return
	}

	if speaker.ID == "" {
		speaker.ID = id
	}

	if err := validateSpeaker(&speaker); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !speakerExists(speaker.ID) {
		utils.Error(w, fmt.Sprintf("speaker with ID '%s' does not exist", id), http.StatusNotFound)
		return
	}

	writeSpeaker(w, r, &speaker)
}

// DeleteSpeaker deletes the file with data for a Speaker given their ID.
func DeleteSpeaker(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if !speakerExists(id) {
		utils.Error(w, fmt.Sprintf("speaker with ID '%s' does not exist", id), http.StatusNotFound)
		return
	}

	if err := os.Remove(fmt.Sprintf("%s%s.json", speakerDataPath, utils.MakeFileName(id))); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// writeSpeaker writes or overwrites a Speaker.
func writeSpeaker(w http.ResponseWriter, r *http.Request, speaker *types.Speaker) {
	content, _ := json.MarshalIndent(speaker, "", " ")

	speaker.Timestamp = time.Now()

	_ = os.WriteFile(fmt.Sprintf("%s%v.json", speakerDataPath, utils.MakeFileName(speaker.ID)), content, 0644)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(speaker)
}

func speakerExists(id string) bool {
	b, _ := os.ReadFile(fmt.Sprintf("%s%s.json", speakerDataPath, utils.MakeFileName(id)))
	return len(b) > 0
}

func validateSpeaker(speaker *types.Speaker) error {
	if speaker.ID == "" || speaker.Name == "" || speaker.Email == "" {
		return fmt.Errorf("speaker ID, Name and Email must be provided")
	}

	return nil
}
