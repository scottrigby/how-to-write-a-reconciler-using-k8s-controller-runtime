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
)

const (
	// CreateFailedReason indicates that the create failed.
	CreateFailedReason string = "CreateFailed"

	// CreateFailedReason indicates that the update failed.
	UpdateFailedReason string = "UpdateFailed"
)
