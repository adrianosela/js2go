package argparse

import (
	"syscall/js"

	"github.com/adrianosela/js2go"
)

type Arg[T any] struct {
	raw   js.Value
	value *T
	err   error
}

func (a *Arg[T]) Raw() js.Value      { return a.raw }
func (a *Arg[T]) Value() (*T, error) { return a.value, a.err }

func Parse[T any](arg js.Value) Arg[T] {
	var value T
	if err := js2go.Decode(arg, &value); err != nil {
		return Arg[T]{
			raw: arg,
			err: err,
		}
	}
	return Arg[T]{
		raw:   arg,
		value: &value,
	}
}

func ParseMultiple[T any](args ...js.Value) []Arg[T] {
	var values []Arg[T]
	for _, arg := range args {
		values = append(values, Parse[T](arg))
	}
	return values
}
