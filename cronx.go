package cronx

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/rizalgowandy/cronx/page"
	"github.com/rizalgowandy/cronx/storage"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/pagination"
	"github.com/rizalgowandy/gdk/pkg/sortx"
	"github.com/robfig/cron/v3"
)

// Default configuration for the manager.
var (
	// DefaultParser supports the v1 where the first parameter is second.
	DefaultParser = cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	DefaultInterceptors = Chain()
	DefaultLocation     = time.Local
	DefaultStorage      = storage.NewNoopClient()
	DefaultAlerter      = NewAlerter()
)

// NewManager create a command controller with a specific config.
func NewManager(opts ...Option) *Manager {
	manager := &Manager{
		commander:            nil,
		downJobs:             nil,
		createdTime:          time.Now().In(DefaultLocation),
		interceptor:          DefaultInterceptors,
		parser:               DefaultParser,
		location:             DefaultLocation,
		autoStart:            true,
		highPriorityDownJobs: true,
		storage:              DefaultStorage,
		alerter:              DefaultAlerter,
	}
	for _, opt := range opts {
		opt(manager)
	}

	commander := cron.New(
		cron.WithParser(manager.parser),
		cron.WithLocation(manager.location),
	)
	if manager.autoStart {
		commander.Start()
	}
	manager.commander = commander
	manager.createdTime = time.Now().In(manager.location)
	return manager
}

// Manager controls all the underlying job.
type Manager struct {
	// commander holds all the underlying cron jobs.
	commander *cron.Cron
	// downJobs describes the list of jobs that have been failed to be registered.
	downJobs []*Job
	// createdTime describes when the command controller created.
	createdTime time.Time

	// Configured using Options.
	//
	// interceptor holds middleware that will be executed before current job operation.
	interceptor Interceptor
	// parser is a custom parser to support v1 that contains second as first parameter.
	parser cron.Parser
	// location describes the timezone current cron is running.
	// By default, the timezone will be the same timezone as the server.
	location *time.Location
	// autoStart determines if the cron will be started automatically or not.
	autoStart bool
	// highPriorityDownJobs determines if the down jobs will be put at the top or bottom of the list.
	highPriorityDownJobs bool
	// storage determines where do we record and read the history data.
	storage storage.Client
	// alerter sends an alert on certain unwanted event.
	alerter AlerterItf
}

// Schedule sets a job to run at specific time.
// Example:
//
//	@every 5m
//	0 */10 * * * * => every 10m
func (m *Manager) Schedule(spec string, job JobItf) error {
	return m.schedule(spec, job, 1, 1)
}

// ScheduleFunc adds a func to the Cron to be run on the given schedule.
func (m *Manager) ScheduleFunc(spec, name string, cmd func(ctx context.Context) error) error {
	return m.Schedule(spec, NewFuncJob(name, cmd))
}

// Schedules sets a job to run multiple times at specific time.
// Symbol */,-? should never be used as separator character.
// These symbols are reserved for cron specification.
//
// Example:
//
//	Spec		: "0 0 1 * * *#0 0 2 * * *#0 0 3 * * *
//	Separator	: "#"
//	This input schedules the job to run 3 times.
func (m *Manager) Schedules(spec, separator string, job JobItf) error {
	if spec == "" {
		return errorx.New("invalid specification")
	}
	if separator == "" {
		return errorx.New("invalid separator")
	}
	schedules := strings.Split(spec, separator)
	for k, v := range schedules {
		if err := m.schedule(v, job, int64(k+1), int64(len(schedules))); err != nil {
			return err
		}
	}
	return nil
}

// SchedulesFunc adds a func to the Cron to be run on the given schedules.
func (m *Manager) SchedulesFunc(
	spec, separator, name string,
	cmd func(ctx context.Context) error,
) error {
	return m.Schedules(spec, separator, NewFuncJob(name, cmd))
}

func (m *Manager) schedule(spec string, job JobItf, waveNumber, totalWave int64) error {
	// Check if spec is correct.
	schedule, err := m.parser.Parse(spec)
	if err != nil {
		downJob := NewJob(m, job, waveNumber, totalWave)
		downJob.Status = StatusCodeDown
		downJob.Error = err.Error()
		m.downJobs = append(m.downJobs, downJob)
		return err
	}

	j := NewJob(m, job, waveNumber, totalWave)
	j.EntryID = m.commander.Schedule(schedule, j)
	return nil
}

// Start starts jobs from running at the next scheduled time.
func (m *Manager) Start() {
	m.commander.Start()
}

// Stop stops active jobs from running at the next scheduled time.
func (m *Manager) Stop() {
	m.commander.Stop()
}

// GetEntries returns all the current registered jobs.
func (m *Manager) GetEntries() []cron.Entry {
	return m.commander.Entries()
}

// GetEntry returns a snapshot of the given entry, or nil if it couldn't be found.
func (m *Manager) GetEntry(id cron.EntryID) *cron.Entry {
	entry := m.commander.Entry(id)
	return &entry
}

// Remove removes a specific job from running.
// Get EntryID from the list job entries manager.GetEntries().
// If job is in the middle of running, once the process is finished it will be removed.
func (m *Manager) Remove(id cron.EntryID) {
	m.commander.Remove(id)
}

// GetInfo returns command controller basic information.
func (m *Manager) GetInfo() map[string]interface{} {
	currentTime := time.Now().In(m.location)

	return map[string]interface{}{
		"data": map[string]interface{}{
			"location":     m.location.String(),
			"created_time": m.createdTime.String(),
			"current_time": currentTime.String(),
			"up_time":      currentTime.Sub(m.createdTime).String(),
		},
	}
}

// GetStatusData returns all jobs status for status page.
func (m *Manager) GetStatusData(sortQuery string) StatusPageData {
	// Default sorting is by id in ascending order.
	if sortQuery == "" {
		sortQuery = page.ColumnID
	}

	// Get status data.
	entries := m.commander.Entries()
	totalEntries := len(entries)

	data := make([]StatusData, totalEntries)
	totalData := totalEntries
	for k, v := range entries {
		data[k].ID = v.ID
		data[k].Job = v.Job.(*Job)
		data[k].Next = v.Next
		data[k].Prev = v.Prev
	}

	// Sort data.
	sorts := sortx.NewSorts(sortQuery)
	for _, v := range sorts {
		sorter := NewStatusDataSorter(v.Key, v.Order, data)
		sort.Sort(sorter)
	}

	downs := m.downJobs
	totalDowns := len(downs)

	totalJobs := totalData + totalDowns
	listStatus := make([]StatusData, totalJobs)

	if m.highPriorityDownJobs {
		// Register down jobs.
		for k, v := range downs {
			listStatus[k].Job = v
		}

		// Register other jobs.
		for k, v := range data {
			idx := totalDowns + k
			listStatus[idx].ID = v.ID
			listStatus[idx].Job = v.Job
			listStatus[idx].Next = v.Next
			listStatus[idx].Prev = v.Prev
		}
	} else {
		// Register other jobs.
		for k, v := range data {
			listStatus[k].ID = v.ID
			listStatus[k].Job = v.Job
			listStatus[k].Next = v.Next
			listStatus[k].Prev = v.Prev
		}

		// Register down jobs.
		for k, v := range downs {
			idx := totalData + k
			listStatus[idx].Job = v
		}
	}

	return StatusPageData{
		Data: listStatus,
		Sort: pagination.Sort{
			Query:   sortQuery,
			Columns: sorts.Map(),
		},
	}
}

// GetHistoryData returns run histories for history page.
func (m *Manager) GetHistoryData(
	ctx context.Context,
	req *Request,
) (HistoryPageData, error) {
	if err := req.Validate(); err != nil {
		return HistoryPageData{}, errorx.E(err)
	}

	// Sort data.
	sorts := sortx.NewSorts(req.Sort)

	// Get data from storage.
	data, err := m.storage.ReadHistories(ctx, &storage.HistoryFilter{
		Sorts:         sorts,
		Limit:         req.Limit,
		StartingAfter: req.StartingAfter,
		EndingBefore:  req.EndingBefore,
	})
	if err != nil {
		if errorx.Is(err, errorx.CodeNotFound) {
			return HistoryPageData{}, nil
		}
		return HistoryPageData{}, errorx.E(err)
	}

	// Create pagination data.
	paginationResp := Response{
		Sort:          req.Sort,
		StartingAfter: req.StartingAfter,
		EndingBefore:  req.EndingBefore,
		Total:         len(data),
		Yielded:       len(data),
		Limit:         req.Limit,
		PreviousURI:   nil,
		NextURI:       nil,
		CursorRange:   nil,
	}
	if len(data) > 0 {
		paginationResp.CursorRange = []int64{
			data[0].ID,
			data[len(data)-1].ID,
		}

		if next, _ := m.storage.ReadHistories(ctx, &storage.HistoryFilter{
			Sorts:         sorts,
			Limit:         1,
			StartingAfter: paginationResp.NextPageCursor(),
		}); len(next) > 0 {
			paginationResp.NextURI = paginationResp.NextPageRequest().URI(&req.url)
		}
		if prev, _ := m.storage.ReadHistories(ctx, &storage.HistoryFilter{
			Sorts:        sorts,
			Limit:        1,
			EndingBefore: paginationResp.PrevPageCursor(),
		}); len(prev) > 0 {
			paginationResp.PreviousURI = paginationResp.PrevPageRequest().URI(&req.url)
		}
	}

	for k := range data {
		data[k].CreatedAt = data[k].CreatedAt.In(m.location)
		data[k].StartedAt = data[k].StartedAt.In(m.location)
		data[k].FinishedAt = data[k].FinishedAt.In(m.location)
	}

	return HistoryPageData{
		Data:       data,
		Pagination: paginationResp,
		Sort: pagination.Sort{
			Query:   req.Sort,
			Columns: sorts.Map(),
		},
	}, nil
}
