package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/scottrigby/cfp-api/pkg/types"
	"github.com/scottrigby/cfp-api/pkg/utils"
)

const proposalsDataPath = "data/proposals/"

// CreateProposal creates a new file with json data about a given Proposal.
func CreateProposal(w http.ResponseWriter, r *http.Request) {
	var proposal types.Proposal
	json.NewDecoder(r.Body).Decode(&proposal)

	if err := validateProposal(&proposal); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if utils.Exists(proposal.ID, proposalsDataPath) {
		utils.Error(w, fmt.Sprintf("proposal with ID '%s' already exists", proposal.ID), http.StatusBadRequest)
		return
	}

	writeProposal(w, r, &proposal)
}

// GetProposal returns the data for a Proposal given the Proposal's ID.
func GetProposalById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		utils.Error(w, "proposal ID must be specified", http.StatusBadRequest)
		return
	}
	b, err := os.ReadFile(fmt.Sprintf("%s%s.json", proposalsDataPath, utils.MakeFileName(id)))

	switch {
	case len(b) == 0:
		utils.Error(w, fmt.Sprintf("could not find proposal with ID '%s'", id), http.StatusNotFound)
		return
	case err != nil:
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var proposal types.Proposal
	if err := json.Unmarshal(b, &proposal); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposal)
}

// GetProposals returns a list with the data of all the files of Proposals
// contained in "data/proposals/".
func GetProposals(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(proposalsDataPath)
	if err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var proposalList []types.Proposal
	for _, file := range files {
		b, err := os.ReadFile(fmt.Sprintf("%s%s", proposalsDataPath, file.Name()))

		switch {
		case len(b) == 0:
			continue
		case err != nil:
			utils.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var proposal types.Proposal
		if err := json.Unmarshal(b, &proposal); err != nil {
			utils.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		proposalList = append(proposalList, proposal)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposalList)
}

// UpdateProposal checks that a file for a Proposal exists given its ID
// then replaces the data for that Proposal by overwriting it.
func UpdateProposal(w http.ResponseWriter, r *http.Request) {
	var proposal types.Proposal
	json.NewDecoder(r.Body).Decode(&proposal)

	id := mux.Vars(r)["id"]
	if id == "" {
		utils.Error(w, "proposal ID must be specified", http.StatusBadRequest)
		return
	}

	if utils.MakeFileName(proposal.ID) != id {
		utils.Error(w, fmt.Sprintf("proposal ID '%s' used as query param does not match ID in request body '%s'", id, proposal.ID), http.StatusBadRequest)
		return
	}

	if proposal.ID == "" {
		proposal.ID = id
	}

	if err := validateProposal(&proposal); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.Exists(proposal.ID, proposalsDataPath) {
		utils.Error(w, fmt.Sprintf("proposal with ID '%s' was not found", proposal.ID), http.StatusBadRequest)
		return
	}

	writeProposal(w, r, &proposal)
}

// DeleteProposal deletes the file with data for a Speaker given their ID.
func DeleteProposal(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		utils.Error(w, "proposal ID must be specified", http.StatusBadRequest)
		return
	}

	if !utils.Exists(id, proposalsDataPath) {
		utils.Error(w, fmt.Sprintf("proposal with ID '%s' was not found", id), http.StatusNotFound)
		return
	}

	if err := os.Remove(fmt.Sprintf("%s%s.json", proposalsDataPath, utils.MakeFileName(id))); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validateProposal(p *types.Proposal) error {
	if p.Type != types.SessionPresentationType && p.Type != types.LightningTalkType {
		return fmt.Errorf("could not validate proposal's talk type; got: %s; want %s or %s", p.Type, types.SessionPresentationType, types.LightningTalkType)
	}

	if p.Submission.Status != types.Draft && p.Submission.Status != types.Final {
		return fmt.Errorf("could not validate proposal's submission status; got: %s; want %s or %s", p.Type, types.Draft, types.Final)
	}

	switch {
	case p.ID == "":
		return fmt.Errorf("proposal ID must be specified")
	case p.SpeakerID == "":
		return fmt.Errorf("speaker ID must be specified")
	}

	_, err := getSpeaker(p.SpeakerID)
	if err != nil {
		return fmt.Errorf("failed to get speaker: %v", err)
	}

	if p.Submission.Status == types.Final {
		switch {
		case p.Title == "":
			return fmt.Errorf("title must be specified")
		case p.Abstract == "":
			return fmt.Errorf("abstract must be specified")
		case p.SpeakerID == "":
			return fmt.Errorf("title must be specified")
		}
	}

	return nil
}

func writeProposal(w http.ResponseWriter, r *http.Request, proposal *types.Proposal) {
	proposal.Submission.LastUpdate = time.Now()

	_ = os.MkdirAll(proposalsDataPath, 0755)

	content, _ := json.MarshalIndent(proposal, "", " ")
	_ = os.MkdirAll(proposalsDataPath, 0755)
	_ = os.WriteFile(fmt.Sprintf("%s%s.json", proposalsDataPath, utils.MakeFileName(proposal.ID)), content, 0644)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposal)
}
