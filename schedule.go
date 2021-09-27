package gotik

import (
	"fmt"
	"strings"
	"time"
)

type Schedule struct {
	ID        string        `json:"id"`
	Comment   string        `json:"comment"`
	Disabled  bool          `json:"disabled"`
	Name      string        `json:"name"`
	NextRun   string        `json:"next-run"`
	Owner     string        `json:"owner"`
	Policy    []string      `json:"policy"`
	StartDate string        `json:"start-date"`
	StartTime string        `json:"start-time"`
	Interval  time.Duration `json:"interval"`
	OnEvent   string        `json:"on-event"`
}

func (s *Schedule) String() string {
	return fmt.Sprintf("%s %s %s %s %s", s.Name, s.StartDate, s.StartTime, s.Comment, s.NextRun)
}

func parseSchedule(props map[string]string) Schedule {
	entry := Schedule{
		ID:        props[".id"],
		Comment:   props["comment"],
		Disabled:  parseBool(props["disabled"]),
		Name:      props["name"],
		NextRun:   props["next-run"],
		Owner:     props["owner"],
		StartDate: props["start-date"],
		StartTime: props["start-time"],
		Interval:  parseDuration(props["interval"]),
		OnEvent:   props["on-event"],
	}
	entry.Policy = strings.Split(props["policy"], ",")
	return entry
}

// Returns a list of all scheduler items
func (c *Client) GetScheduler() ([]Schedule, error) {
	entries := make([]Schedule, 0, 8)
	detail, err := c.RunCmd("/system/scheduler/print")
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseSchedule(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// Add a new Scheduler item
func (c *Client) AddSchedule(s Schedule) (string, error) {
	if err := s.validator(); err != nil {
		return "", err
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/system/scheduler/add")
	parts = append(parts, s.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

// RemoveSchedule removes a Scheduler item by ID. No check is made to see if the item exists.
// Passing an empty ID is a null but successful operation.
func (c *Client) RemoveSchedule(id string) error {
	if len(id) > 0 {
		_, err := c.Run("/system/scheduler/remove", "=.id="+id)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateSchedule update an existing Scheduler item
func (c *Client) UpdateSchedule(s Schedule) (string, error) {
	if err := s.validator(); err != nil {
		return "", err
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/system/script/set", "=.id="+s.ID)
	parts = append(parts, s.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (s *Schedule) validator() error {
	if len(s.Name) == 0 {
		return fmt.Errorf("invalid name supplied")
	}
	if len(s.StartTime) == 0 {
		s.StartTime = "startup"
	} else {
		// check that it's hh:mm:ss
	}
	if len(s.StartDate) > 0 {
		//check that it's mmm/dd/yyyy
	}
	return nil
}

func (s *Schedule) parter() []string {
	parts := make([]string, 0, 10)
	parts = append(parts, fmt.Sprintf("=name=%s", s.Name))
	parts = append(parts, fmt.Sprintf("=disabled=%t", s.Disabled))
	if len(s.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", s.Comment))
	}
	if len(s.Policy) > 0 {
		parts = append(parts, fmt.Sprintf("=policy=%s", strings.Join(s.Policy, ",")))
	}
	if len(s.StartDate) > 0 {
		parts = append(parts, fmt.Sprintf("=start-date=%s", s.StartDate))
	}
	if len(s.StartTime) > 0 {
		parts = append(parts, fmt.Sprintf("=start-time=%s", s.StartTime))
	}
	parts = append(parts, fmt.Sprintf("=interval=%s", s.Interval.Round(time.Second).String()))
	parts = append(parts, fmt.Sprintf("=on-event=%s", s.OnEvent))
	return parts
}
