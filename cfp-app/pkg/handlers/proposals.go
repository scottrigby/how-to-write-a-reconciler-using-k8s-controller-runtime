package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/scottrigby/cfp/pkg/types"
)

func CreateProposal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var proposal types.Proposal
	json.NewDecoder(r.Body).Decode(&proposal)
	file, _ := json.MarshalIndent(proposal, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("data/proposals/%v.json", proposal.ID), file, 0644)
	json.NewEncoder(w).Encode(proposal)
}

func GetProposals(w http.ResponseWriter, r *http.Request) {}

func GetProposalById(w http.ResponseWriter, r *http.Request) {}

func UpdateProposal(w http.ResponseWriter, r *http.Request) {}

func DeleteProposal(w http.ResponseWriter, r *http.Request) {}

func validateTalkType(t string) error {
	switch t {
	case types.SessionPresentationType, types.LightningTalkType:
		return nil
	default:
		return fmt.Errorf("could not validate Talk type")
	}
}
