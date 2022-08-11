package promise_test

import (
	"testing"
	"time"

	"github.com/IfanTsai/go-lib/promise"
	"github.com/stretchr/testify/require"
)

func TestPromise_Then(t *testing.T) {
	t.Parallel()

	var result string

	promise.New(func(resolve, reject promise.Handler) {
		go func() {
			time.Sleep(time.Second)
			resolve("hello")
		}()
	}).Then(func(value interface{}) interface{} {
		return promise.New(func(resolve, reject promise.Handler) {
			go func() {
				time.Sleep(time.Second)
				resolve(value.(string) + " world")
			}()
		})
	}).Then(func(value interface{}) interface{} {
		result = value.(string)

		return nil
	})

	require.Equal(t, "", result)

	time.Sleep(time.Second * 3)
	require.Equal(t, "hello world", result)
}
