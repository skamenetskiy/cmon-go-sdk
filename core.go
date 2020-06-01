package cmon

import (
	"fmt"
	"time"
)

const (
	ModuleAuth     = "auth"
	ModuleClusters = "clusters"

	RequestStatusOk             = "Ok"             // The request was successfully processed.
	RequestStatusInvalidRequest = "InvalidRequest" // Something was fundamentally wrong with the request.
	RequestStatusObjectNotFound = "ObjectNotFound" // The referenced object (e.g. the cluster) was not found.
	RequestStatusTryAgain       = "TryAgain"       // The request can not at the moment processed.
	RequestStatusUnknownError   = "UnknownError"   // The exact error could not be identified.
	RequestStatusAccessDenied   = "AccessDenied"   // The authenticated user has insufficient rights.
	RequestStatusAuthRequired   = "AuthRequired"   // The client has to Authenticate first.
)

func NewError(t, m string) error {
	return &Error{t, m}
}

func NewErrorFromResponseData(d *WithResponseData) error {
	return &Error{d.RequestStatus, d.ErrorString}
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return fmt.Sprintf("%s: %s",
		err.Type,
		err.Message)
}

type WithOperation struct {
	Operation string `json:"operation"`
}

type WithClusterID struct {
	ClusterID uint64 `json:"cluster_id"`
}

type WithClassName struct {
	ClassName string `json:"class_name"`
}

type WithResponseData struct {
	RequestID        uint64    `json:"request_id"`
	RequestCreated   time.Time `json:"request_created"`
	RequestProcessed time.Time `json:"request_processed"`
	RequestStatus    string    `json:"request_status"`
	ErrorString      string    `json:"error_string"`
}

type WithTotal struct {
	Total int64 `json:"total"`
}
