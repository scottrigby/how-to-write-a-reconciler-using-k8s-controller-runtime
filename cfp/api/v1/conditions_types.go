package v1

const Finalizer = "finalizers.talks.kubecon.na"

const (
	// CreateFailedCondition indicates a transient or persistent creation failure
	// of an external resource.
	// This is a "negative polarity" or "abnormal-true" type, and is only
	// present on the resource if it is True.
	CreateFailedCondition string = "CreateFailed"

	// UpdateFailedCondition indicates a transient or persistent update failure
	// of an external resource.
	// This is a "negative polarity" or "abnormal-true" type, and is only
	// present on the resource if it is True.
	UpdateFailedCondition string = "CreateFailed"

	// FetchFailedCondition indicates a transient or persistent fetch failure
	// of an external resource.
	// This is a "negative polarity" or "abnormal-true" type, and is only
	// present on the resource if it is True.
	FetchFailedCondition string = "FetchFailed"
)

const (
	// CreateFailedReason indicates that the create failed.
	CreateFailedReason string = "CreateFailed"

	// UpdateFailedReason indicates that the update failed.
	UpdateFailedReason string = "UpdateFailed"

	// FetchFailedReason indicates that the fetch failed.
	FetchFailedReason string = "FetchFailed"
)
