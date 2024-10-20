package standard

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xerror"
)

type ResponseStatus string

const (
	ResponseStatusSuccess = "success"
	ResponseStatusError   = "error"
)

type Metadata struct {
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id"`
}

func NewMetadata(ctx context.Context) *Metadata {
	return &Metadata{
		Timestamp: time.Now(),
		RequestID: xcontext.RequestID(ctx),
	}
}

type Response struct {
	Status           ResponseStatus `json:"status,omitempty"`
	Data             any            `json:"data,omitempty"`
	Error            string         `json:"error,omitempty"`
	ErrorDescription string         `json:"error_description,omitempty"`
	Metadata         *Metadata      `json:"metadata,omitempty"`
}

func NewResponse(data any) *Response {
	return &Response{
		Status: ResponseStatusSuccess,
		Data:   data,
	}
}

func NewErrorResponse(ctx context.Context, err error) *Response {
	serviceErr := xerror.ServiceError{}
	if errors.As(err, &serviceErr) {
		return NewErrorResponseWithMessage(ctx, serviceErr.Error(), serviceErr.Message())
	}

	xcontext.Logger(ctx).Critical("internal-error", "err", err)
	return NewUnexpectedErrorResponse(ctx)
}

func NewUnexpectedErrorResponse(ctx context.Context) *Response {
	return NewErrorResponseWithMessage(ctx,
		"server_error",
		"an unexpected error occurred, please contact to admin if you see this error",
	)
}

func NewErrorResponseWithMessage(ctx context.Context, err string, description string) *Response {
	response := &Response{Status: ResponseStatusError, Metadata: NewMetadata(ctx)}

	response.Error = err
	response.ErrorDescription = description

	return response
}

func SetQuery(ctx context.Context, q url.Values, err error) {
	errResp := NewErrorResponse(ctx, err)

	q.Set("error", errResp.Error)
	q.Set("error_description", errResp.ErrorDescription)
	q.Set("timestamp", errResp.Metadata.Timestamp.Format(TimeLayout))
	q.Set("request_id", errResp.Metadata.RequestID)
}
