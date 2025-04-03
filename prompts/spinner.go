package prompts

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/orochaa/go-clack/core/utils"
	"github.com/orochaa/go-clack/prompts/symbols"
	isunicodesupported "github.com/orochaa/go-clack/third_party/is-unicode-supported"
	"github.com/orochaa/go-clack/third_party/picocolors"
	"github.com/orochaa/go-clack/third_party/sisteransi"
)

// SpinnerIndicator defines the type for spinner indicators.
type SpinnerIndicator int

const (
	// SpinnerDotsIndicator shows a loading animation with dots.
	SpinnerDotsIndicator SpinnerIndicator = iota
	// SpinnerStaticDotsIndicator shows static dots.
	SpinnerStaticDotsIndicator
	// SpinnerTimerIndicator displays a timer alongside the loading message.
	SpinnerTimerIndicator
)

type SpinnerOptions struct {
	Context       context.Context
	Output        io.Writer
	Indicator     SpinnerIndicator
	Frames        []string
	FrameInterval time.Duration
	OnCancel      func()
}

type SpinnerController struct {
	options     SpinnerOptions
	isCI        bool
	IsCancelled bool
	ticker      *time.Ticker
	message     string
	stop        func()
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

	return &SpinnerController{
		options: options,
		isCI:    isCI,
		ticker:  time.NewTicker(options.FrameInterval),
	}
}

func (s *SpinnerController) write(msg string) {
	s.options.Output.Write([]byte(msg))
}

func (s *SpinnerController) trimMessageDots(msg string) string {
	dotsRegex := regexp.MustCompile(`\.+$`)
	return dotsRegex.ReplaceAllString(msg, "")
}

func (s *SpinnerController) clearMessage(msg string) {
	if msg == "" {
		return
	}
	if s.isCI {
		s.write("\r\n")
	}
	prevLines := utils.SplitLines(msg)
	s.write(sisteransi.MoveCursor(-len(prevLines)+1, -999))
	s.write(sisteransi.EraseDown())
}

// Starts the spinner animation with the provided message
func (s *SpinnerController) Start(msg string) {
	s.write(sisteransi.HideCursor())
	s.write(picocolors.Gray(symbols.BAR) + "\r\n")

	s.ticker.Reset(s.options.FrameInterval)

	ctx, stop := context.WithCancel(s.options.Context)
	s.stop = stop
	cancel := make(chan os.Signal, 2)
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)

	frameIndex := 0
	startTime := time.Now()
	s.message = s.trimMessageDots(msg)
	prevMessage := s.message

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.ticker.Stop()
				return
			case sig := <-cancel:
				var cancelMsg string
				switch sig {
				case syscall.SIGTERM, syscall.SIGINT:
					cancelMsg = "Cancelled"
					s.IsCancelled = true
				default:
					cancelMsg = "Something went wrong"
				}
				s.Stop(cancelMsg, 1)
				if s.IsCancelled {
					s.options.OnCancel()
				}
				os.Exit(0)
			case <-s.ticker.C:
				if s.isCI && s.message == prevMessage {
					continue
				}
				s.clearMessage(prevMessage)
				prevMessage = s.message
				frame := picocolors.Magenta(s.options.Frames[frameIndex])
				duration := time.Since(startTime)
				formattedFrame := s.formatFrame(s.options.Indicator, frame, s.message, duration)
				s.write(formattedFrame)
				if frameIndex+1 < len(s.options.Frames) {
					frameIndex++
				} else {
					frameIndex = 0
				}
			}
		}
	}()
}

// Updates the spinner's displayed message
func (s *SpinnerController) Message(msg string) {
	s.message = s.trimMessageDots(msg)
}

// Stops the spinner animation and displays a final message with a status indicator.
func (s *SpinnerController) Stop(msg string, code int) {
	s.stop()
	s.clearMessage(s.message)
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
		dotsRegex := regexp.MustCompile(`\.{2,}$`)
		s.message = dotsRegex.ReplaceAllString(msg, ".")
	}
	s.write(sisteransi.ShowCursor())
	s.write(fmt.Sprintf("%s %s\n", step, s.message))
}

func (s *SpinnerController) formatFrame(indicator SpinnerIndicator, frame string, message string, duration time.Duration) string {
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
