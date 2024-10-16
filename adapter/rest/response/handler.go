package response

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x"
	"github.com/xybor/x/errorx"
	"github.com/xybor/x/logging"
	"github.com/xybor/x/xcontext"
)

type ErrorResponse struct {
	ErrMsg any `json:"error"`
}

type ResponseHandler struct {
	err  error
	resp any
	code int
}

func NewResponseHandler(val any, err error) *ResponseHandler {
	return (&ResponseHandler{resp: val, err: err, code: -1}).
		Map(http.StatusUnauthorized, usecase.ErrUnauthorized).
		Map(http.StatusForbidden, usecase.ErrForbidden)
}

func (h *ResponseHandler) Map(code int, errs ...error) *ResponseHandler {
	if h.err == nil || h.code != -1 {
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

func (h *ResponseHandler) WriteHTTPResponse(ctx context.Context, w http.ResponseWriter) {
	h.Map(http.StatusInternalServerError)

	if h.code == -1 {
		h.code = http.StatusOK
	}

	var resp any = h.resp

	if h.err != nil {
		var serviceErr errorx.ServiceError
		switch {
		case errors.As(h.err, &serviceErr):
			resp = ErrorResponse{ErrMsg: serviceErr.Message}
		default:
			resp = ErrorResponse{ErrMsg: http.StatusText(h.code)}
		}

		logging.LogError(xcontext.Logger(ctx), h.err)
	}

	Write(ctx, w, h.code, resp)
}

func HandleParseError(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		panic("do not pass a nil error here")
	}

	var code int
	response := ErrorResponse{}
	if errors.Is(err, x.ErrHTTPBadRequest) {
		code = http.StatusBadRequest
		response.ErrMsg = err.Error()
	} else {
		xcontext.Logger(ctx).Debug("failed to parse data", "err", err.Error())
		code = http.StatusInternalServerError
		response.ErrMsg = "Internal Server Error"
	}

	Write(ctx, w, code, response)
}

func WriteErrorMsg(ctx context.Context, w http.ResponseWriter, code int, msg string, a ...any) {
	response := ErrorResponse{ErrMsg: fmt.Sprintf(msg, a...)}
	Write(ctx, w, code, response)
}

func Write(ctx context.Context, w http.ResponseWriter, code int, obj any) {
	if err := x.WriteHTTPResponseJSON(w, code, obj); err != nil {
		xcontext.Logger(ctx).Critical("failed to write response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
