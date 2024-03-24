package logger

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		err := InitLogger("info")
		require.NoError(t, err)
		UseLogger().Info("test message")
	})
}
