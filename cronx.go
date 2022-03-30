package cronx

import (
	"context"
	"strings"
	"time"

	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/robfig/cron/v3"
)

// Config defines the config for the command controller.
type Config struct {
	// Address determines the address will we serve the json and frontend status.
	// Empty string meaning we won't serve the current job status.
	Address string
	// Location describes the timezone current cron is running.
	Location *time.Location
}

var (
	defaultConfig = Config{
		Address:  ":8998",
		Location: time.Local,
	}

	commandController *CommandController
)

// Default creates a cron with default config.
// HTTP server is built in as side car by default.
func Default(interceptors ...Interceptor) {
	Custom(defaultConfig, interceptors...)
}

// New creates a cron without HTTP server built in.
func New(interceptors ...Interceptor) {
	Custom(Config{}, interceptors...)
}

// Custom creates a cron with custom config.
// For advance user, allow custom modification.
func Custom(config Config, interceptors ...Interceptor) {
	// If there is invalid config use the default config instead.
	if config.Location == nil {
		config.Location = defaultConfig.Location
	}

	// Create new command controller and start the underlying jobs.
	commandController = NewCommandController(config, interceptors...)

	// Check if client want to start a server to serve json and frontend.
	if config.Address != "" {
		go NewSideCarServer(commandController)
	}
}

// Schedule sets a job to run at specific time.
// Example:
//  @every 5m
//  0 */10 * * * * => every 10m
func Schedule(spec string, job JobItf) error {
	return schedule(spec, "", job, 1, 1)
}

// ScheduleWithName sets a job to run at specific time with a Job name
// Example:
//  @every 5m
//  0 */10 * * * * => every 10m
func ScheduleWithName(name, spec string, job JobItf) error {
	return schedule(spec, name, job, 1, 1)
}

func schedule(spec, jobName string, job JobItf, waveNumber, totalWave int64) error {
	if commandController == nil || commandController.Commander == nil {
		return errorx.New("cronx has not been initialized")
	}

	// Check if spec is correct.
	schedule, err := commandController.Parser.Parse(spec)
	if err != nil {
		downJob := NewJob(job, jobName, waveNumber, totalWave)
		downJob.Status = StatusCodeDown
		downJob.Error = err.Error()
		commandController.UnregisteredJobs = append(
			commandController.UnregisteredJobs,
			downJob,
		)
		return err
	}

	j := NewJob(job, jobName, waveNumber, totalWave)
	j.EntryID = commandController.Commander.Schedule(schedule, j)
	return nil
}

// Schedules sets a job to run multiple times at specific time.
// Symbol */,-? should never be used as separator character.
// These symbols are reserved for cron specification.
//
// Example:
//	Spec		: "0 0 1 * * *#0 0 2 * * *#0 0 3 * * *
//	Separator	: "#"
//	This input schedules the job to run 3 times.
func Schedules(spec, separator string, job JobItf) error {
	if spec == "" {
		return errorx.New("invalid specification")
	}
	if separator == "" {
		return errorx.New("invalid separator")
	}
	schedules := strings.Split(spec, separator)
	for k, v := range schedules {
		if err := schedule(v, "", job, int64(k+1), int64(len(schedules))); err != nil {
			return err
		}
	}
	return nil
}

// Every executes the given job at a fixed interval.
// The interval provided is the time between the job ending and the job being run again.
// The time that the job takes to run is not included in the interval.
// Minimal time is 1 sec.
func Every(duration time.Duration, job JobItf) {
	if commandController == nil || commandController.Commander == nil {
		return
	}

	j := NewJob(job, "", 1, 1)
	j.EntryID = commandController.Commander.Schedule(cron.Every(duration), j)
}

// Stop stops active jobs from running at the next scheduled time.
func Stop() {
	if commandController == nil || commandController.Commander == nil {
		return
	}

	commandController.Commander.Stop()
}

// GetEntries returns all the current registered jobs.
func GetEntries() []cron.Entry {
	if commandController == nil || commandController.Commander == nil {
		return nil
	}

	return commandController.Commander.Entries()
}

// GetEntry returns a snapshot of the given entry, or nil if it couldn't be found.
func GetEntry(id cron.EntryID) *cron.Entry {
	if commandController == nil || commandController.Commander == nil {
		return nil
	}

	entry := commandController.Commander.Entry(id)
	return &entry
}

// Remove removes a specific job from running.
// Get EntryID from the list job entries cronx.GetEntries().
// If job is in the middle of running, once the process is finished it will be removed.
func Remove(id cron.EntryID) {
	if commandController == nil || commandController.Commander == nil {
		return
	}

	commandController.Commander.Remove(id)
}

// GetStatusData returns all jobs status.
func GetStatusData() []StatusData {
	if commandController == nil {
		return nil
	}

	return commandController.StatusData()
}

// GetStatusJSON returns all jobs status as map[string]interface.
func GetStatusJSON() map[string]interface{} {
	if commandController == nil {
		return nil
	}

	return commandController.StatusJSON()
}

// GetInfo returns current cron check basic information.
func GetInfo() map[string]interface{} {
	if commandController == nil {
		return nil
	}

	return commandController.Info()
}

// Func is a type to allow callers to wrap a raw func.
// Example:
//	cronx.Schedule("@every 5m", cronx.Func(myFunc))
type Func func(ctx context.Context) error

func (r Func) Run(ctx context.Context) error {
	return r(ctx)
}
