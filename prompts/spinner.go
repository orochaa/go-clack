package prompts

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	isunicodesupported "github.com/Mist3rBru/go-clack/third_party/is-unicode-supported"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
	"github.com/Mist3rBru/go-clack/third_party/sisteransi"
)

type SpinnerOptions struct {
	Context context.Context
	Output  io.Writer
}

type SpinnerController struct {
	// Starts the spinner animation with the provided message
	Start func(msg string)
	// Updates the spinner's displayed message
	Message func(msg string)
	// Stops the spinner animation and displays a final message with a status indicator.
	Stop func(msg string, code int)
}

// Spinner initializes and returns a SpinnerController with the provided options.
func Spinner(options SpinnerOptions) *SpinnerController {
	if options.Context == nil {
		options.Context = context.Background()
	}
	if options.Output == nil {
		options.Output = os.Stdout
	}

	var ticker *time.Ticker

	var message, prevMessage string

	var frames []string
	var frameIndex int
	var frameInterval time.Duration

	const dotsInterval float32 = 0.125
	var dotsTimer float32

	if isunicodesupported.IsUnicodeSupported() {
		frames = []string{"◒", "◐", "◓", "◑"}
		frameInterval = 80 * time.Millisecond
	} else {
		frames = []string{"•", "o", "O", "0"}
		frameInterval = 120 * time.Millisecond
	}

	ctx, stopSpinner := context.WithCancel(options.Context)
	ticker = time.NewTicker(frameInterval)

	isCI := os.Getenv("CI") == "true"

	write := func(str string) {
		options.Output.Write([]byte(str))
	}

	clearPrevMessage := func() {
		if prevMessage == "" {
			return
		}
		if isCI {
			write("\r\n")
		}
		prevLines := utils.SplitLines(prevMessage)
		write(sisteransi.MoveCursor(-len(prevLines)+1, -999))
		write(sisteransi.EraseDown())
	}

	return &SpinnerController{
		Start: func(msg string) {
			write(sisteransi.HideCursor())
			write(picocolors.Gray(symbols.BAR) + "\r\n")

			frameIndex = 0
			dotsTimer = 0
			message = parseMessage(msg)

			go func() {
				ticker.Reset(frameInterval)

				for {
					select {
					case <-ctx.Done():
						ticker.Stop()
					case <-ticker.C:
						if isCI && message == prevMessage {
							continue
						}
						clearPrevMessage()
						prevMessage = message
						frame := picocolors.Magenta(frames[frameIndex])
						var loadingDots string
						if isCI {
							loadingDots = "..."
						} else {
							loadingDots = strings.Repeat(".", min(int(dotsTimer), 3))
						}
						write(fmt.Sprintf("%s %s%s", frame, message, loadingDots))
						if frameIndex+1 < len(frames) {
							frameIndex++
						} else {
							frameIndex = 0
						}
						if int(dotsTimer) < 4 {
							dotsTimer += dotsInterval
						} else {
							dotsTimer = 0
						}
					}
				}
			}()
		},
		Message: func(msg string) {
			message = parseMessage(msg)
		},
		Stop: func(msg string, code int) {
			stopSpinner()
			clearPrevMessage()
			var step string
			switch code {
			case 0:
				step = picocolors.Green(symbols.STEP_SUBMIT)
			case 1:
				step = picocolors.Red(symbols.STEP_CANCEL)
			default:
				step = picocolors.Red(symbols.STEP_ERROR)
			}
			if msg != "" {
				message = parseMessage(msg)
			}
			write(sisteransi.ShowCursor())
			write(fmt.Sprintf("%s %s\n", step, message))
		},
	}
}

func parseMessage(msg string) string {
	dotsRegex := regexp.MustCompile(`\.+$`)
	return dotsRegex.ReplaceAllString(msg, "")
}
