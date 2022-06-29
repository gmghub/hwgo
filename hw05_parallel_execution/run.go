package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidArgument     = errors.New("invalid argument")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m <= 0 ignore errors and return nil.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrInvalidArgument
	}
	if len(tasks) < 1 {
		return nil
	}
	var tasksIndex int
	var running int32

	taskchan := make(chan Task)
	errchan := make(chan error)
	errors := make([]error, 0)

	for i := 1; i <= n; i++ {
		atomic.AddInt32(&running, 1)
		go func(tc <-chan Task, ec chan<- error) {
			defer func() {
				atomic.AddInt32(&running, -1)
			}()
			for t := range tc {
				if err := t(); err != nil {
					ec <- err
				}
			}
		}(taskchan, errchan)
	}

	for {
		if tasksIndex >= len(tasks) {
			break
		}
		if m > 0 && len(errors) >= m {
			break
		}
		select {
		case e := <-errchan:
			errors = append(errors, e)
		case taskchan <- tasks[tasksIndex]:
			tasksIndex++
		default:
		}
	}
	close(taskchan)

	for {
		if atomic.LoadInt32(&running) < 1 {
			break
		}
		select {
		case e := <-errchan:
			errors = append(errors, e)
		default:
		}
	}

	if m > 0 && len(errors) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
