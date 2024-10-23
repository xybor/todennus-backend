package response

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/xybor/todennus-backend/adapter/rest/standard"
	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xerror"
	"github.com/xybor/x/xhttp"
)

type ResponseHandler struct {
	err         error
	f           any
	resp        any
	code        int
	defaultCode int
}

func NewResponseHandler(ctx context.Context, f any, resp any, err error) *ResponseHandler {
	if timeoutErr := context.Cause(ctx); timeoutErr != nil && errors.Is(timeoutErr, usecase.ErrServerTimeout) {
		err = usecase.ErrServerTimeout.Hide(err, "timeout")
	}

	return (&ResponseHandler{
		err: err,
		f:   f, resp: resp,
		code: -1,
	}).WithDefaultCode(http.StatusOK).
		Map(http.StatusGatewayTimeout, usecase.ErrServerTimeout)
	// Map(http.StatusGatewayTimeout, usecase.ErrServerTimeout)
	// Map(http.StatusBadRequest, usecase.ErrRequestInvalid).
	// Map(http.StatusUnauthorized, usecase.ErrUnauthenticated).
	// Map(http.StatusForbidden, usecase.ErrForbidden).
}

func (h *ResponseHandler) WithDefaultCode(code int) *ResponseHandler {
	h.defaultCode = code
	return h
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
		h.code = h.defaultCode
	}

	var resp any
	if h.err != nil {
		resp = standard.NewErrorResponse(ctx, h.err)
	} else if h.resp != nil {
		if h.f == nil {
			resp = standard.NewResponse(h.resp)
		} else {
			fValue := reflect.ValueOf(h.f)
			resp = standard.NewResponse(fValue.Call([]reflect.Value{reflect.ValueOf(h.resp)})[0].Interface())
		}
	}

	Write(ctx, w, h.code, resp)
}

func (h *ResponseHandler) WriteHTTPResponseWithoutWrap(ctx context.Context, w http.ResponseWriter) {
	h.Map(http.StatusInternalServerError)

	if h.code == -1 {
		h.code = h.defaultCode
	}

	var resp any
	if h.err != nil {
		errResp := standard.NewErrorResponse(ctx, h.err)
		errResp.Status = ""
		resp = errResp
	} else if h.resp != nil {
		if h.f == nil {
			resp = h.resp
		} else {
			fValue := reflect.ValueOf(h.f)
			resp = fValue.Call([]reflect.Value{reflect.ValueOf(h.resp)})[0].Interface()
		}
	}

	Write(ctx, w, h.code, resp)
}

func (h *ResponseHandler) Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, code int) {
	h.Map(http.StatusInternalServerError)

	if h.code == -1 {
		h.code = code
	}

	var redirect string
	if h.err != nil {
		Write(ctx, w, h.code, standard.NewErrorResponse(ctx, h.err))
		return
	} else if h.resp != nil {
		if h.f == nil {
			redirect = h.resp.(string)
		} else {
			fValue := reflect.ValueOf(h.f)
			redirect = fValue.Call([]reflect.Value{reflect.ValueOf(h.resp)})[0].Interface().(string)
		}
	}

	Redirect(ctx, w, r, redirect, h.code)
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		panic("do not pass a nil error here")
	}

	var code int
	response := &standard.Response{}
	switch {
	case xerror.Is(err, xhttp.ErrHTTPBadRequest, usecase.ErrRequestInvalid):
		code = http.StatusBadRequest
		response = standard.NewErrorResponseWithMessage(ctx, "invalid_request", err.Error())
	default:
		code = http.StatusInternalServerError
		response = standard.NewUnexpectedErrorResponse(ctx)

		xcontext.Logger(ctx).Debug("failed-to-parse-data", "err", err)
	}

	Write(ctx, w, code, response)
}

func WriteError(ctx context.Context, w http.ResponseWriter, code int, err error) {
	Write(ctx, w, code, standard.NewErrorResponse(ctx, err))
}

func Write(ctx context.Context, w http.ResponseWriter, code int, resp any) {
	xcontext.SessionManager(ctx).Save(w, xcontext.Session(ctx))
	if err := xhttp.WriteResponseJSON(w, code, resp); err != nil {
		xcontext.Logger(ctx).Critical("failed to write response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, url string, code int) {
	xcontext.SessionManager(ctx).Save(w, xcontext.Session(ctx))
	http.Redirect(w, r, url, code)
}
