package backoff

import "time"

// Retry the function f until it does not return error or BackOff stops.
// f is guaranteed to be run at least once.
// It is the caller's responsibility to reset b after Retry returns.
//
// Retry sleeps the goroutine for the duration returned by BackOff after a
// failed operation returns.
//
// Usage:
// 	operation := func() error {
// 		// An operation that may fail
// 	}
//
// 	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
// 	if err != nil {
// 		// Operation has failed.
// 	}
//
// 	// Operation is successfull.
//
func Retry(f func() error, b BackOff) error { return RetryNotify(f, b, nil) }

type ErrorHandler func(wait time.Duration, err error)

type RetryStrategy struct {
	b        BackOff
	notify   ErrorHandler
	attempts int
}

// RetryN Like Retry but will retry only up to N times
func RetryN(n int, b BackOff, notify ErrorHandler, f func() error) error {
	var err error
	var next time.Duration

	b.Reset()
	for i := 0; i < n; i++ {
		if err = f(); err == nil {
			return nil
		}

		if next = b.NextBackOff(); next == Stop {
			return err
		}

		if notify != nil {
			notify(next, err)
		}

		time.Sleep(next)
	}

	return err
}

// RetryNotify calls notify function with the error and wait duration for each failed attempt before sleep.
func RetryNotify(f func() error, b BackOff, notify func(err error, wait time.Duration)) error {
	var err error
	var next time.Duration

	b.Reset()
	for {
		if err = f(); err == nil {
			return nil
		}

		if next = b.NextBackOff(); next == Stop {
			return err
		}

		if notify != nil {
			notify(err, next)
		}

		time.Sleep(next)
	}
}
