package progress

import (
	"context"
	"fmt"
	"github.com/zcubbs/hotpot/pkg/x/style"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Task func() error

func RunTask(task Task, spinner bool) error {
	var err error
	if spinner {
		err = RunTaskWithSpinner(task)
	} else {
		err = task()
	}

	if err != nil {
		style.PrintError("failed")
	} else {
		style.PrintSuccess("completed")
	}
	return err
}

func RunTaskWithSpinner(task Task) error {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start the spinner in a separate goroutine
	wg.Add(1)
	go func() {
		spinner(ctx, 50*time.Millisecond)
		wg.Done()
	}()

	// Run the task
	err := task()

	// Once the task is done, stop the spinner and wait for it to clean up
	cancel()
	wg.Wait()

	return err
}

func spinner(ctx context.Context, delay time.Duration) {
	frames := []string{"\\", "|", "!", "/", "-"}

	if err := tPut("civis"); err != nil {
		fmt.Print("\033[?25l")
	}

	for {
		for _, frame := range frames {
			select {
			case <-ctx.Done():
				if err := tPut("cvvis"); err != nil {
					fmt.Print("\033[?25h")
				}
				return
			default:
				fmt.Print(frame, "\010")
				time.Sleep(delay)
			}
		}
	}
}

func tPut(arg string) error {
	cmd := exec.Command("tput", arg)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
