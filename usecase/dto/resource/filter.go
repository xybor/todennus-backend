package resource

import (
	"context"
	"reflect"

	"github.com/xybor/todennus-backend/pkg/scope"
	"github.com/xybor/todennus-backend/pkg/xcontext"
)

type Filterer struct {
	ctx      context.Context
	obj      any
	filtered bool
}

func Filter(ctx context.Context, obj any) *Filterer {
	return &Filterer{ctx: ctx, obj: obj}
}

func (f *Filterer) When(target bool) *Filterer {
	if !f.filtered && target {
		f.setzero()
	}
	return f
}

func (f *Filterer) WhenNot(target bool) *Filterer {
	return f.When(!target)
}

func (f *Filterer) WhenNotContainsScope(target scope.Scope) *Filterer {
	return f.WhenNot(xcontext.Scope(f.ctx).Contains(target))
}

func (f *Filterer) WhenRequestUserNot(userID int64) *Filterer {
	return f.WhenNot(xcontext.RequestUserID(f.ctx) == userID)
}

func (f *Filterer) IfNotAdmin() *Filterer {
	return f.WhenNot(xcontext.IsAdmin(f.ctx))
}

func (f *Filterer) setzero() {
	reflect.ValueOf(f.obj).Elem().SetZero()
}
