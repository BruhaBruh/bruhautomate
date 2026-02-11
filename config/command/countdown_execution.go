package command

import (
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/BruhaBruh/bruhautomate/pkg/duration"
	"github.com/BruhaBruh/bruhautomate/pkg/plural"
)

type CountdownExecution struct {
	Type      string            `yaml:"type"`
	Time      duration.Duration `yaml:"time"`
	Execution Executable        `yaml:"execution"`
}

func (e *CountdownExecution) Execute(commands []Command, args []string, flags map[string]string) error {
	d := e.time(flags)

	notifications := e.notifications(d)
	var wg sync.WaitGroup
	wg.Add(len(notifications))
	for _, notification := range notifications {
		remaining := max(d-notification, 0)
		notificationCopy := notification
		go func(d time.Duration) {
			timer := time.NewTimer(d)
			<-timer.C
			countdownFlags := make(map[string]string, len(flags)+1)
			maps.Copy(countdownFlags, flags)
			countdownFlags["countdown"] = e.formatDuration(notificationCopy)
			e.Execution.Execute(commands, args, countdownFlags)
			wg.Done()
		}(remaining)
	}
	wg.Wait()
	return nil
}

func (e *CountdownExecution) time(flags map[string]string) time.Duration {
	if value, ok := flags["t"]; ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if value, ok := flags["time"]; ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return e.Time.Duration
}

func (e *CountdownExecution) formatDuration(d time.Duration) string {
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

func (e *CountdownExecution) notifications(d time.Duration) []time.Duration {
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
