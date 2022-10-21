package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/cfp/api/v1"
)

type SpeakerChangePredicate struct {
	predicate.Funcs
}

func (SpeakerChangePredicate) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}

	oldSource, ok := e.ObjectOld.(*talksv1.Speaker)
	if !ok {
		return false
	}

	newSource, ok := e.ObjectNew.(*talksv1.Speaker)
	if !ok {
		return false
	}

	// take action if the speaker ID is created or updated
	if oldSource.Status.ID == "" && newSource.Status.ID != "" {
		return true
	}

	if oldSource.Status.ID != "" && newSource.Status.ID != "" && oldSource.Status.ID != newSource.Status.ID {
		return true
	}

	return false
}

func (SpeakerChangePredicate) Create(e event.CreateEvent) bool {
	return false
}

func (SpeakerChangePredicate) Delete(e event.DeleteEvent) bool {
	return true
}
