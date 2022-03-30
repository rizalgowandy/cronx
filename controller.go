package cronx

import (
	"time"

	"github.com/robfig/cron/v3"
)

// CommandController controls all the underlying job.
type CommandController struct {
	// Commander holds all the underlying cron jobs.
	Commander *cron.Cron
	// Interceptor holds middleware that will be executed before current job operation.
	Interceptor Interceptor
	// Parser is a custom parser to support v1 that contains second as first parameter.
	Parser cron.Parser
	// UnregisteredJobs describes the list of jobs that have been failed to be registered.
	UnregisteredJobs []*Job
	// Address determines the address will we serve the json and frontend status.
	// Empty string meaning we won't serve the current job status.
	Address string
	// Location describes the timezone current cron is running.
	// By default the timezone will be the same timezone as the server.
	Location *time.Location
	// CreatedTime describes when the command controller created.
	CreatedTime time.Time
}

// NewCommandController create a command controller with a specific config.
func NewCommandController(config Config, interceptors ...Interceptor) *CommandController {
	if config.Location == nil {
		config.Location = defaultConfig.Location
	}

	// Support the v1 where the first parameter is second.
	parser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	// Create the commander.
	commander := cron.New(
		cron.WithParser(parser),
		cron.WithLocation(config.Location),
	)
	commander.Start()

	// Create command controller.
	return &CommandController{
		Commander:        commander,
		Interceptor:      Chain(interceptors...),
		Parser:           parser,
		UnregisteredJobs: nil,
		Address:          config.Address,
		Location:         config.Location,
		CreatedTime:      time.Now().In(config.Location),
	}
}

// Info returns command controller basic information.
func (c *CommandController) Info() map[string]interface{} {
	if c.Location == nil {
		c.Location = defaultConfig.Location
	}

	currentTime := time.Now().In(c.Location)

	return map[string]interface{}{
		"data": map[string]interface{}{
			"location":     c.Location.String(),
			"created_time": c.CreatedTime.String(),
			"current_time": currentTime.String(),
			"up_time":      currentTime.Sub(c.CreatedTime).String(),
		},
	}
}

// StatusData returns all jobs status.
func (c *CommandController) StatusData() []StatusData {
	if c.Commander == nil {
		return nil
	}

	entries := c.Commander.Entries()
	totalEntries := len(entries)

	downs := c.UnregisteredJobs
	totalDowns := len(downs)

	totalJobs := totalEntries + totalDowns
	listStatus := make([]StatusData, totalJobs)

	// Register down jobs.
	for k, v := range downs {
		listStatus[k].Job = v
	}

	// Register other jobs.
	for k, v := range entries {
		idx := totalDowns + k
		listStatus[idx].ID = v.ID
		listStatus[idx].Job = v.Job.(*Job)
		listStatus[idx].Next = v.Next
		listStatus[idx].Prev = v.Prev
	}

	return listStatus
}

// StatusJSON returns all jobs status as map[string]interface.
func (c *CommandController) StatusJSON() map[string]interface{} {
	return map[string]interface{}{
		"data": c.StatusData(),
	}
}
