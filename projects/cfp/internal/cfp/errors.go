package cfp

import (
	"errors"
	"fmt"
)

type ErrorReason struct {
	Reason  string
	Summary string
}

func (e ErrorReason) Error() string {
	return e.Summary
}

type Error struct {
	Reason ErrorReason
	Err    error
}

func (e *Error) Error() string {
	if e.Reason.Error() == "" {
		return e.Err.Error()
	}

	return fmt.Sprintf("%s: %s", e.Reason.Error(), e.Err.Error())
}

func (e *Error) Is(target error) bool {
	if e.Reason == target {
		return true
	}
	return errors.Is(e.Err, target)
}

var (
	ErrCreateRequest = ErrorReason{Reason: "InvalidRequest", Summary: "invalid request"}
	ErrMakeRequest   = ErrorReason{Reason: "RequestFailed", Summary: "request failed"}
	ErrCreateSpeaker = ErrorReason{Reason: "CreateSpeakerFailed", Summary: "error creating speaker"}
	ErrUpdateSpeaker = ErrorReason{Reason: "UpdateSpeakerFailed", Summary: "error updating speaker"}
	ErrFetchSpeaker  = ErrorReason{Reason: "FetchSpeakerFailed", Summary: "error fetching speaker"}
	ErrDeleteSpeaker = ErrorReason{Reason: "DeleteSpeakerFailed", Summary: "error deleting speaker"}
)
