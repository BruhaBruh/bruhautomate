package executors

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/BruhaBruh/bruhautomate/internal/flag"
	"github.com/BruhaBruh/bruhautomate/internal/plural"
)

const TimerExecutorType = "timer"

type TimerExecutor struct {
	Type     string        `yaml:"type"`
	Time     time.Duration `yaml:"time"`
	Executor Executor      `yaml:"executor"`
}

var _ Executor = (*TimerExecutor)(nil)

func (e *TimerExecutor) Execute(
	args []string,
	flags *flag.Flags,
	commandFlags *flag.Flags,
	executeCommand func(
		name string,
		args []string,
		flags *flag.Flags,
		commandFlags *flag.Flags,
	) error,
) error {
	duration := e.duration(flags)

	notifications := e.notifications(duration)
	var wg sync.WaitGroup
	wg.Add(len(notifications))
	for _, notification := range notifications {
		remaining := max(duration-notification, 0)
		notificationCopy := notification
		go func(d time.Duration) {
			timer := time.NewTimer(d)
			<-timer.C
			flag := &flag.Flag{
				Name:     "countdown",
				HasValue: true,
				Value:    e.formatDuration(notificationCopy),
			}
			flags.AddFlag(flag)
			commandFlags.AddFlag(flag)
			e.Executor.Execute(args, flags, commandFlags, executeCommand)
			wg.Done()
		}(remaining)
	}
	wg.Wait()
	return nil
}

func (e *TimerExecutor) Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag {
	flags := make([]flag.Flag, 0, 16)
	flags = append(flags, flag.Flag{
		Name:        "time",
		Shortcut:    "t",
		Description: "Time for timer. Example 5h30m15s",
		HasValue:    true,
	})
	flags = append(flags, e.Executor.Flags(commandFlagsGetter)...)
	return flags
}

func (e *TimerExecutor) duration(flags *flag.Flags) time.Duration {
	flag, err := flags.Find("time")
	if err != nil {
		return e.Time
	}
	if duration, err := time.ParseDuration(flag.Value); err == nil {
		return duration
	}
	return e.Time
}

func (e *TimerExecutor) notifications(d time.Duration) []time.Duration {
	var notifications []time.Duration
	seen := make(map[time.Duration]bool)

	add := func(t time.Duration) {
		if !seen[t] {
			notifications = append(notifications, t)
			seen[t] = true
		}
	}

	add(d)

	for h := d / time.Hour; h >= 1; h-- {
		add(h * time.Hour)
	}

	for h := d / time.Hour; h >= 0; h-- {
		t := h*time.Hour + 30*time.Minute
		if t < d && t > 0 {
			add(t)
		}
	}

	for m := d / time.Minute; m >= 5; m-- {
		if m%5 == 0 && m*60 < time.Duration(d.Seconds()) {
			add(time.Duration(m) * time.Minute)
		}
	}

	if d >= time.Minute {
		add(time.Minute)
	}

	if d >= 30*time.Second {
		add(30 * time.Second)
	}

	if d >= 10*time.Second {
		add(10 * time.Second)
	}

	for s := 5; s >= 1; s-- {
		if s > int(d.Seconds()) {
			continue
		}
		add(time.Duration(s) * time.Second)
	}

	for i := 0; i < len(notifications)/2; i++ {
		j := len(notifications) - 1 - i
		notifications[i], notifications[j] = notifications[j], notifications[i]
	}

	return notifications
}

func (e *TimerExecutor) formatDuration(d time.Duration) string {
	totalSec := int(d.Seconds())

	h := totalSec / 3600
	m := (totalSec % 3600) / 60
	s := totalSec % 60

	parts := make([]string, 0, 3)

	if h > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", h, plural.Make(h, "час", "часа", "часов")))
	}
	if m > 0 && len(parts) > 0 && s == 0 {
		parts = append(parts, fmt.Sprintf("и %d %s", m, plural.Make(m, "минуту", "минуты", "минут")))
	} else if m > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", m, plural.Make(m, "минуту", "минуты", "минут")))
	}
	if s > 0 && len(parts) > 0 {
		parts = append(parts, fmt.Sprintf("и %d %s", s, plural.Make(s, "секунду", "секунды", "секунд")))
	} else if s > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d %s", s, plural.Make(s, "секунду", "секунды", "секунд")))
	}

	return strings.Join(parts, " ")
}
