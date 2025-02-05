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

// SpinnerIndicator defines the type for spinner indicators.
type SpinnerIndicator int

const (
	// SpinnerDotsIndicator shows a loading animation with dots.
	SpinnerDotsIndicator SpinnerIndicator = iota
	// SpinnerDotsIndicator shows static dots.
	SpinnerStaticDotsIndicator SpinnerIndicator = iota
	// SpinnerTimerIndicator displays a timer alongside the loading message.
	SpinnerTimerIndicator
)

type SpinnerOptions struct {
	Context       context.Context
	Output        io.Writer
	Indicator     SpinnerIndicator
	Frames        []string
	FrameInterval time.Duration
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
	isCI := os.Getenv("CI") == "true"

	if options.Context == nil {
		options.Context = context.Background()
	}
	if options.Output == nil {
		options.Output = os.Stdout
	}
	isUnicodeSupported := isunicodesupported.IsUnicodeSupported()
	if options.Frames == nil {
		if isUnicodeSupported {
			options.Frames = []string{"◒", "◐", "◓", "◑"}
		} else {
			options.Frames = []string{"•", "o", "O", "0"}
		}
	}
	if options.FrameInterval == 0 {
		if isUnicodeSupported {
			options.FrameInterval = 80 * time.Millisecond
		} else {
			options.FrameInterval = 120 * time.Millisecond
		}
	}
	if isCI {
		options.Indicator = SpinnerStaticDotsIndicator
	}

	var ctx context.Context
	var ticker *time.Ticker
	var startTime time.Time

	var stopSpinner func()

	var message, prevMessage string

	var frameIndex int

	ticker = time.NewTicker(options.FrameInterval)

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

			ctx, stopSpinner = context.WithCancel(options.Context)
			ticker.Reset(options.FrameInterval)

			frameIndex = 0
			startTime = time.Now()

			message = trimMessageDots(msg)

			go func() {
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
						frame := picocolors.Magenta(options.Frames[frameIndex])
						duration := time.Since(startTime)
						formattedFrame := formatSpinnerFrame(options.Indicator, frame, message, duration)
						write(formattedFrame)
						if frameIndex+1 < len(options.Frames) {
							frameIndex++
						} else {
							frameIndex = 0
						}
					}
				}
			}()
		},
		Message: func(msg string) {
			message = trimMessageDots(msg)
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
				message = trimMessageDots(msg)
			}
			write(sisteransi.ShowCursor())
			write(fmt.Sprintf("%s %s\n", step, message))
		},
	}
}

func trimMessageDots(msg string) string {
	dotsRegex := regexp.MustCompile(`\.+$`)
	return dotsRegex.ReplaceAllString(msg, "")
}

func formatSpinnerFrame(indicator SpinnerIndicator, frame string, message string, duration time.Duration) string {
	switch indicator {
	case SpinnerDotsIndicator:
		dotsCounter := int(duration.Seconds()) % 4
		return fmt.Sprintf("%s %s%s", frame, message, strings.Repeat(".", dotsCounter))
	case SpinnerStaticDotsIndicator:
		return fmt.Sprintf("%s %s...", frame, message)
	case SpinnerTimerIndicator:
		min := int(duration.Minutes())
		secs := int(duration.Seconds()) - (min * 60)
		if min > 0 {
			return fmt.Sprintf("%s %s [%dm %ds]", frame, message, min, secs)
		}
		return fmt.Sprintf("%s %s [%ds]", frame, message, secs)
	default:
		return fmt.Sprintf("%s %s", frame, message)
	}
}
