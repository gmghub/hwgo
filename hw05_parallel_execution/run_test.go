package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("invalid or empty values of arguments", func(t *testing.T) {
		tasks := make([]Task, 0)
		var err error
		err = Run(tasks, 10, 10)
		require.Nilf(t, err, "not nil on empty task list")

		err = Run(tasks, 0, 10)
		require.Truef(t, errors.Is(err, ErrInvalidArgument), "not error on zero N")

		err = Run(tasks, -1, 10)
		require.Truef(t, errors.Is(err, ErrInvalidArgument), "not error on negative N")
	})

	t.Run("test zero or negative M", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Nilf(t, err, "should be nil on M == 0")

		workersCount = 10
		maxErrorsCount = -1
		runTasksCount = 0
		err = Run(tasks, workersCount, maxErrorsCount)
		require.Nilf(t, err, "should be nil on M < 0")
	})

	t.Run("test concurrency execution with no sleep", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime int64

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				now := time.Now()
				atomic.AddInt32(&runTasksCount, 1)
				elapsed := int64(time.Since(now))
				atomic.AddInt64(&sumTime, elapsed)
				return nil
			})
		}

		workersCount := 10
		maxErrorsCount := 1

		runtime.GOMAXPROCS(1)
		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTimeSerially := int64(time.Since(start))
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")

		runTasksCount = 0
		runtime.GOMAXPROCS(4)
		start = time.Now()
		err = Run(tasks, workersCount, maxErrorsCount)
		elapsedTimeConcurrently := int64(time.Since(start))
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")

		require.Less(t, elapsedTimeConcurrently, elapsedTimeSerially/2, "tasks were run sequentially?")
	})
}
