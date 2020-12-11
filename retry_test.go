package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/jkgx/logrus"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/tj/assert"
)

func TestRetry(t *testing.T) {
	t.Run("case=fails after timeout", func(t *testing.T) {
		l, _ := test.NewNullLogger()
		logger := logrus.New("", "", logrus.UseLogger(l))

		randomErr := fmt.Errorf("some error")

		err := Retry(logger, 100*time.Millisecond, 100*time.Millisecond, func() error {
			return randomErr
		})

		assert.Equal(t, err, randomErr)
	})

	t.Run("case=logs error when failing", func(t *testing.T) {
		l, hook := test.NewNullLogger()
		logger := logrus.New("", "", logrus.UseLogger(l))

		const errPattern = "error %d"

		var i int
		err := Retry(logger, 100*time.Millisecond, 200*time.Millisecond, func() error {
			defer func() { i++ }()
			return fmt.Errorf(errPattern, i)
		})

		assert.Equal(t, fmt.Errorf(errPattern, 1), err)
		assert.Len(t, hook.AllEntries(), 2)
		assert.Equal(t, hook.LastEntry().Data["error"], map[string]interface{}{"message": fmt.Errorf(errPattern, 1).Error()})
	})
}
