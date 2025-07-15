package js2go

import (
	"fmt"
	"syscall/js"
)

// RegisterPromise sets a function on the DOM that returns a promise that get's either resolved
// with the return value of the function, or rejected with the error returned by the function.
func RegisterPromise(name string, fn func(js.Value, []js.Value) (any, error)) {
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		return js.Global().Get("Promise").New(js.FuncOf(func(_ js.Value, promiseArgs []js.Value) any {
			resolve := promiseArgs[0]
			reject := promiseArgs[1]

			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(js.ValueOf(fmt.Sprintf("panic: %v", r)))
					}
				}()

				result, err := fn(this, args)
				if err != nil {
					reject.Invoke(js.ValueOf(err.Error()))
					return
				}

				resolve.Invoke(js.ValueOf(result))
			}()

			return nil
		}))
	}))
}
