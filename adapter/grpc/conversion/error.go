package conversion

import (
	"context"
	"errors"

	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseHandler[D any] struct {
	err         error
	resp        D
	code        codes.Code
	defaultCode codes.Code
}

func NewResponseHandler[D any](ctx context.Context, resp D, err error) *ResponseHandler[D] {
	if timeoutErr := context.Cause(ctx); timeoutErr != nil && errors.Is(timeoutErr, usecase.ErrServerTimeout) {
		err = usecase.ErrServerTimeout.Hide(err, "timeout")
	}

	return (&ResponseHandler[D]{
		err:  err,
		resp: resp,
		code: codes.Unknown,
	}).WithDefaultCode(codes.OK).
		Map(codes.DeadlineExceeded, usecase.ErrServerTimeout)
}

func (h *ResponseHandler[D]) WithDefaultCode(code codes.Code) *ResponseHandler[D] {
	h.defaultCode = code
	return h
}

func (h *ResponseHandler[D]) Map(code codes.Code, errs ...error) *ResponseHandler[D] {
	if code == codes.Unknown {
		panic("do not use code unknown")
	}

	if h.err == nil || h.code != codes.Unknown {
		return h
	}

	if len(errs) == 0 {
		h.code = code
	} else {
		for _, err := range errs {
			if errors.Is(h.err, err) {
				h.code = code
				break
			}
		}
	}

	return h
}

func (h *ResponseHandler[D]) Finalize(ctx context.Context) (D, error) {
	h.Map(codes.Internal)

	if h.code == codes.Unknown {
		h.code = h.defaultCode
	}

	var defaultResp D
	if h.err != nil {
		var richError xerror.RichError
		if errors.As(h.err, &richError) {
			if richError.Detail() != nil {
				attrs := []any{"err", richError.Detail()}
				attrs = append(attrs, richError.Attributes()...)
				if errors.Is(h.err, usecase.ErrServer) {
					xcontext.Logger(ctx).Warn(richError.Event(), attrs...)
				} else {
					xcontext.Logger(ctx).Debug(richError.Event(), attrs...)
				}
			}

			h.err = richError.Reduce()
		} else {
			xcontext.Logger(ctx).Critical("internal-error", "err", h.err)
			h.err = errors.New("unexpected_server_error: an unexpected error occured")
		}

		return defaultResp, status.Errorf(h.code, h.err.Error())
	}

	return h.resp, nil
}
