package gorm_provider

type QueryFieldFilter[T any] struct {
	Defined bool
	Val     T
}

func NewQueryFieldFilter[T any](val T) QueryFieldFilter[T] {
	return QueryFieldFilter[T]{Defined: true, Val: val}
}

type QuerySliceFieldFilter[T any] struct {
	Defined bool
	Val     []T
}

func NewQuerySliceFieldFilter[T any](val []T) QuerySliceFieldFilter[T] {
	return QuerySliceFieldFilter[T]{Defined: true, Val: val}
}

func (o *QuerySliceFieldFilter[T]) Append(items ...T) {
	o.Defined = true
	o.Val = append(o.Val, items...)
}

func (o *QuerySliceFieldFilter[T]) Concat(items []T) {
	o.Defined = true
	o.Val = append(o.Val, items...)
}
