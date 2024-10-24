package conversion

import (
	"context"
	"errors"

	"github.com/xybor/todennus-backend/usecase"
	"github.com/xybor/x/xcontext"
	"github.com/xybor/x/xerror"
)

func ReduceError(ctx context.Context, err error) error {
	var richError xerror.RichError
	if errors.As(err, &richError) {
		if richError.Detail() != nil {
			attrs := []any{"err", richError.Detail()}
			attrs = append(attrs, richError.Attributes()...)
			if errors.Is(err, usecase.ErrServer) {
				xcontext.Logger(ctx).Warn(richError.Event(), attrs...)
			} else {
				xcontext.Logger(ctx).Debug(richError.Event(), attrs...)
			}
		}

		return richError.Reduce()
	}

	xcontext.Logger(ctx).Critical("internal-error", "err", err)
	return errors.New("unexpected_server_error")
}
