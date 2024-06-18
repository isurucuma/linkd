package link

import (
	"log/slog"
	"math"
	"testing"
)

func TestMain(m *testing.M) {
	// silence the logger by setting the level to a higher
	slog.SetLogLoggerLevel(math.MaxInt)
	m.Run()
}
