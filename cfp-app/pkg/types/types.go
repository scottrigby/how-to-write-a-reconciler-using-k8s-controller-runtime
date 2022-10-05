package types

import (
	"time"
)

const (
	SessionPresentationType = "SessionPresentation"
	// PanelType               = "PanelDiscussion"
	LightningTalkType = "LigntningTalk"
	// KeynoteType             = "Keynote"
)

// Speaker represents a speaker who is submitting a proposal.
type Speaker struct {
	ID        string
	Name      string
	Bio       string
	Email     string
	Timestamp time.Time
}

// Proposal represents an instance of a proposed talk that is submitted to a CFP.
type Proposal struct {
	ID               string
	Title            string
	Abstract         string
	Type             string
	Speaker          *Speaker
	Final            bool
	SubmissionStatus SubmissionStatus
}

const (
	Draft = "draft"
	Final = "final"
)

// Submission represents the status of a Proposal created by the user.
type Submission struct {
	LastUpdate time.Time
	Status     string
}
