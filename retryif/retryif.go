package retryif

import (
	"errors"
	"reflect"
)

type Config struct {
	NumRetries  int
	ShouldRetry func(error) bool
}

// Retry executes the provided function with retry logic
func Retry(fn interface{}, config Config) (interface{}, error) {
	if config.NumRetries < 1 {
		return nil, errors.New("number of retries must be greater than 0")
	}

	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		return nil, errors.New("fn must be a function")
	}

	var ok bool
	var err error
	var result []reflect.Value

	for i := 0; i < config.NumRetries; i++ {

		err = nil // reset error each time, we try to execute

		result = fnVal.Call(nil)
		if len(result) != 2 {
			return nil, errors.New("fn must return exactly two values: (result, error)")
		}

		if result[1].Interface() == nil {
			break
		}

		err, ok = result[1].Interface().(error)
		if !ok {
			return nil, errors.New("second return value of fn must be an error")
		}

		if !config.ShouldRetry(err) {
			break
		}
	}

	return result[0].Interface(), err
}
