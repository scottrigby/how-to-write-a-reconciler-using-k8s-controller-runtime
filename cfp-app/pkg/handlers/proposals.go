package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/scottrigby/cfp/pkg/types"
	"github.com/scottrigby/cfp/pkg/utils"
)

const proposalsDataPath = "data/proposals/"

func CreateProposal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var proposal types.Proposal
	json.NewDecoder(r.Body).Decode(&proposal)
	file, _ := json.MarshalIndent(proposal, "", " ")

	if err := validateProposal(&proposal); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = os.WriteFile(fmt.Sprintf("%s%s.json", proposalsDataPath, utils.MakeFileName(proposal.ID)), file, 0644)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposal)
}

func GetProposals(w http.ResponseWriter, r *http.Request) {}

func GetProposalById(w http.ResponseWriter, r *http.Request) {}

func UpdateProposal(w http.ResponseWriter, r *http.Request) {}

func DeleteProposal(w http.ResponseWriter, r *http.Request) {}

func validateProposal(p *types.Proposal) error {
	if p.Type != types.SessionPresentationType && p.Type != types.LightningTalkType {
		return fmt.Errorf("could not validate proposal's talk type; got: %s; want %s or %s", p.Type, types.SessionPresentationType, types.LightningTalkType)
	}

	return nil
}
