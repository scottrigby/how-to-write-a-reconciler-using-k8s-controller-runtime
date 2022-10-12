package v1

const Finalizer = "finalizers.talks.kubecon.na"

const (
	// FetchFailedCondition indicates a transient or persistent fetch failure
	// of an external resource.
	// This is a "negative polarity" or "abnormal-true" type, and is only
	// present on the resource if it is True.
	FetchFailedCondition string = "FetchFailed"
)

const (
	// FetchFailedReason indicates that the fetch failed.
	FetchFailedReason string = "FetchFailed"
)
