package resource

import (
	"context"
	"reflect"

	"github.com/xybor/todennus-backend/pkg/scope"
	"github.com/xybor/todennus-backend/pkg/xcontext"
)

type Filterer[T any] struct {
	ctx      context.Context
	obj      *T
	filtered bool
	target   T
}

func Set[T any](ctx context.Context, obj *T, target T) *Filterer[T] {
	return &Filterer[T]{ctx: ctx, obj: obj, target: target}
}

func Filter[T any](ctx context.Context, obj *T) *Filterer[T] {
	var t T
	return Set(ctx, obj, t)
}

func (f *Filterer[T]) When(cond bool) *Filterer[T] {
	if !f.filtered && cond {
		f.setzero()
	}
	return f
}

func (f *Filterer[T]) WhenNot(cond bool) *Filterer[T] {
	return f.When(!cond)
}

func (f *Filterer[T]) WhenNotContainsScope(target scope.Scope) *Filterer[T] {
	return f.WhenNot(xcontext.Scope(f.ctx).Contains(target))
}

func (f *Filterer[T]) WhenRequestUserNot(userID int64) *Filterer[T] {
	return f.WhenNot(xcontext.RequestUserID(f.ctx) == userID)
}

func (f *Filterer[T]) setzero() {
	reflect.ValueOf(f.obj).Elem().Set(reflect.ValueOf(f.target))
}
