package retryif

import (
	"errors"
	"testing"
)

func TestRetrySuccess(t *testing.T) {
	fn := func() (interface{}, error) {
		return "success", nil
	}

	config := Config{
		NumRetries: 3,
		ShouldRetry: func(err error) bool {
			return true
		},
	}

	result, err := Retry(fn, config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "success" {
		t.Fatalf("expected result 'success', got %v", result)
	}
}

func TestRetryFailAndSucceed(t *testing.T) {
	attempts := 0
	fn := func() (interface{}, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("temporary failure")
		}
		return "success", nil
	}

	config := Config{
		NumRetries: 5,
		ShouldRetry: func(err error) bool {
			return err.Error() == "temporary failure"
		},
	}

	result, err := Retry(fn, config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "success" {
		t.Fatalf("expected result 'success', got %v", result)
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %v", attempts)
	}
}

func TestRetryExceedsRetries(t *testing.T) {
	fn := func() (interface{}, error) {
		return nil, errors.New("persistent failure")
	}

	config := Config{
		NumRetries: 3,
		ShouldRetry: func(err error) bool {
			return err.Error() == "persistent failure"
		},
	}

	result, err := Retry(fn, config)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if err.Error() != "persistent failure" {
		t.Fatalf("expected error 'persistent failure', got %v", err)
	}

	if result != nil {
		t.Fatalf("expected result nil, got %v", result)
	}
}

func TestRetryShouldNotRetry(t *testing.T) {
	fn := func() (interface{}, error) {
		return nil, errors.New("do not retry")
	}

	config := Config{
		NumRetries: 3,
		ShouldRetry: func(err error) bool {
			return false
		},
	}

	result, err := Retry(fn, config)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if err.Error() != "do not retry" {
		t.Fatalf("expected error 'do not retry', got %v", err)
	}

	if result != nil {
		t.Fatalf("expected result nil, got %v", result)
	}
}

func TestRetryInvalidFunction(t *testing.T) {
	fn := "not a function"

	config := Config{
		NumRetries: 3,
		ShouldRetry: func(err error) bool {
			return true
		},
	}

	_, err := Retry(fn, config)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestRetryFunctionWrongSignature(t *testing.T) {
	fn := func() {}

	config := Config{
		NumRetries: 3,
		ShouldRetry: func(err error) bool {
			return true
		},
	}

	_, err := Retry(fn, config)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestRetryNilError(t *testing.T) {
	fn := func() (interface{}, error) {
		return "success", nil
	}

	config := Config{
		NumRetries: 3,
		ShouldRetry: func(err error) bool {
			return err != nil
		},
	}

	result, err := Retry(fn, config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "success" {
		t.Fatalf("expected result 'success', got %v", result)
	}
}

func TestRetryEarlySuccess(t *testing.T) {
	attempts := 0
	fn := func() (interface{}, error) {
		attempts++
		if attempts == 2 {
			return "success", nil
		}
		return nil, errors.New("temporary failure")
	}

	config := Config{
		NumRetries: 5,
		ShouldRetry: func(err error) bool {
			return err.Error() == "temporary failure"
		},
	}

	result, err := Retry(fn, config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "success" {
		t.Fatalf("expected result 'success', got %v", result)
	}

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %v", attempts)
	}
}
