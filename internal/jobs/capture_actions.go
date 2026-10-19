package jobs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/notify"
	"github.com/Cimmerial/tps/internal/ui"
)

// WatchCapture polls logPath until pattern matches, then runs configured actions.
func WatchCapture(opts CaptureOptions) {
	if opts.Timeout <= 0 {
		opts.Timeout = defaultCaptureTimeout
	}
	pattern := config.ResolveCapturePattern(opts.Pattern)
	actions := opts.Actions
	if len(actions) == 0 {
		actions = config.DefaultCaptureActions()
	}

	deadline := time.Now().Add(opts.Timeout)
	var match string
	ui.BrailleWait("watching output for link…", func() {
		for time.Now().Before(deadline) {
			data, err := os.ReadFile(opts.LogPath)
			if err == nil && len(data) > 0 {
				m, err := config.ExtractCapture(string(data), pattern)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("capture: invalid pattern: "+err.Error()))
					return
				}
				if m != "" {
					match = m
					return
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	})

	if match == "" {
		fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("capture: no match yet (still running in terminal)"))
		return
	}

	var msgs []string
	for _, action := range actions {
		switch action {
		case config.CapturePrint:
			ui.PrintCaptureResult(match, "")
		case config.CaptureClipboard:
			if err := ui.CopyToClipboard(match); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("capture: clipboard: "+err.Error()))
			} else {
				msgs = append(msgs, ui.Green("copied to clipboard"))
			}
		case config.CaptureSMS:
			if err := notify.SendSMS(opts.SMSTo, match, opts.SMSShortcut); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("capture: sms: "+err.Error()))
			} else {
				msgs = append(msgs, ui.Green("sent sms"))
			}
		}
	}
	if len(msgs) > 0 {
		fmt.Println(strings.Join(msgs, " · "))
	}
}
