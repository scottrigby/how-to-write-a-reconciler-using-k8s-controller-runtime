package types

import (
	"time"
)

const (
	SessionPresentationType = "SessionPresentation"
	PanelType               = "PanelDiscussion"
	LightiningTalkType      = "LigntningTalk"
	KeynoteType             = "Keynote"
)

// Speaker represents a speaker who is submitting a proposal.
type Speaker struct {
	ID        int64
	Name      string
	Bio       string
	Email     string
	Timestamp time.Time
}

// Proposal represents an instance of a proposed talk that is submitted to a CFP.
type Proposal struct {
	ID                int64
	Title             string
	Abstract          string
	Type              string
	Speakers          []Speaker
	Final             bool
	ApplicationStatus ApplicationStatus
}

// ApplicationStatus represents the status of a Proposal.
type ApplicationStatus struct {
	Timestamp time.Time
	Submitted bool
}
