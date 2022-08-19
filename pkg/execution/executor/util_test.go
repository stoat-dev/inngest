package executor

import (
	"context"
	"testing"
	"time"

	"github.com/inngest/inngest/inngest"
	"github.com/inngest/inngest/pkg/execution/state"
	"github.com/inngest/inngest/pkg/execution/state/inmemory"
	"github.com/stretchr/testify/require"
)

func TestParseWait(t *testing.T) {
	ctx := context.Background()

	state := inmemory.NewStateInstance(
		inngest.Workflow{},
		state.Identifier{},
		state.Metadata{},
		map[string]any{
			"data": time.Now().Format(time.RFC3339),
		},
		nil,
		nil,
	)

	tests := []struct {
		wait     string
		duration time.Duration
		err      error
	}{
		{
			wait:     "1h30m",
			duration: 90 * time.Minute,
			err:      nil,
		},
		// expression as duration
		{
			wait:     "duration('1h')",
			duration: time.Hour,
			err:      nil,
		},
		// event data as expression
		{
			wait:     "date(event.data) + duration('45m30s')",
			duration: (time.Minute * 45) + (time.Second * 30),
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.wait, func(t *testing.T) {
			duration, err := ParseWait(ctx, test.wait, state)
			require.NoError(t, err)
			require.WithinDuration(
				t,
				time.Now().Add(test.duration),
				time.Now().Add(duration),
				time.Second,
			)
		})
	}

}