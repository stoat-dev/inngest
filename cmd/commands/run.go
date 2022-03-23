package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"cuelang.org/go/cue"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/inngest/event-schemas/pkg/fakedata"
	"github.com/inngest/inngestctl/pkg/cli"
	"github.com/inngest/inngestctl/pkg/function"
	"github.com/spf13/cobra"
)

var runSeed int64

func NewCmdRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Run a serverless function locally",
		Example: "inngestctl run",
		Run:     doRun,
	}

	cmd.Flags().Int64Var(&runSeed, "seed", 0, "Sets the seed for deterministically generating random events")
	return cmd
}

func doRun(cmd *cobra.Command, args []string) {
	fn, err := function.Load(".")
	if err != nil {
		fmt.Println("\n" + cli.RenderError("No inngest.json or inngest.cue file found in your current directory") + "\n")
		os.Exit(1)
		return
	}

	err = runFunction(cmd.Context(), *fn)
	if err != nil {
		os.Exit(1)
	}
}

// runFunction builds the function's images and runs the function.
func runFunction(ctx context.Context, fn function.Function) error {
	if runSeed <= 0 {
		rand.Seed(time.Now().UnixNano())
		runSeed = rand.Int63n(1_000_000)
	}

	evt, err := event(ctx, fn)
	if err != nil {
		return err
	}

	actions, err := fn.Actions()
	if err != nil {
		return err
	}
	if len(actions) != 1 {
		return fmt.Errorf("running step-functions locally is not yet supported")
	}

	// Build the image.
	ui, err := cli.NewRunUI(ctx, cli.RunUIOpts{
		Action: actions[0],
		Event:  evt,
		Seed:   runSeed,
	})
	if err != nil {
		return err
	}
	if err := tea.NewProgram(ui).Start(); err != nil {
		return err
	}
	// So we can exit with a non-zero code.
	return ui.Error()
}

// event retrieves the event for use within testing the function.  It first checks stdin
// to see if we're passed an event, or resorts to generating a fake event based off of
// the function's event type.
func event(ctx context.Context, fn function.Function) (map[string]interface{}, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// Read stdin
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		evt := scanner.Bytes()

		data := map[string]interface{}{}
		err := json.Unmarshal(evt, &data)
		return data, err
	}

	return fakeEvent(ctx, fn)
}

func fakeEvent(ctx context.Context, fn function.Function) (map[string]interface{}, error) {
	evtTriggers := []function.Trigger{}
	for _, t := range fn.Triggers {
		if t.EventTrigger != nil {
			evtTriggers = append(evtTriggers, t)
		}
	}

	i := rand.Intn(len(evtTriggers))
	if evtTriggers[i].EventTrigger.Definition == nil {
		return nil, nil
	}

	def, err := evtTriggers[i].EventTrigger.Definition.Cue()
	if err != nil {
		return nil, err
	}

	r := &cue.Runtime{}
	inst, err := r.Compile(".", def)
	if err != nil {
		return nil, err
	}

	fakedata.DefaultOptions.Rand = rand.New(rand.NewSource(runSeed))

	val, err := fakedata.Fake(ctx, inst.Value())
	if err != nil {
		return nil, err
	}

	mapped := map[string]interface{}{}
	err = val.Decode(&mapped)

	return mapped, err
}
