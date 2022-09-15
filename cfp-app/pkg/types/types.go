package types

import (
	"time"

	"github.com/blang/semver/v4"
)

const (
	SessionPresentationType = "SessionPresentation"
	PanelType               = "PanelDiscussion"
	LightiningTalkType      = "LigntningTalk"
	KeynoteType             = "Keynote"
)

// Speaker represents a speaker who is submitting a proposal.
type Speaker struct {
	ID        int64 `json:"id"`
	Name      string
	Bio       string
	Email     string
	Timestamp time.Time
}

// Proposal represents an instance of a proposed talk that is submitted to a CFP.
type Proposal struct {
	ID        int64
	Title     string
	Abstract  string
	Type      string
	Speakers  []Speaker
	Status    string
	Save      semver.Version
	Submitted bool
	Timestamp time.Time
}
