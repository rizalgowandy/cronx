package cronx

import (
	"sort"
	"strings"
	"unicode"

	"github.com/rizalgowandy/gdk/pkg/sortx"
)

const (
	SortKeyID      sortx.Key = "id"
	SortKeyName    sortx.Key = "name"
	SortKeyStatus  sortx.Key = "status"
	SortKeyPrevRun sortx.Key = "prev_run"
	SortKeyNextRun sortx.Key = "next_run"
	SortKeyLatency sortx.Key = "latency"
)

func NewStatusDataSorter(key sortx.Key, order sortx.Order, data []StatusData) sort.Interface {
	var sorter sort.Interface
	switch key {
	case SortKeyID:
		sorter = byID(data)
	case SortKeyName:
		sorter = byName(data)
	case SortKeyStatus:
		sorter = byStatus(data)
	case SortKeyPrevRun:
		sorter = byPrevRun(data)
	case SortKeyNextRun:
		sorter = byNextRun(data)
	case SortKeyLatency:
		sorter = byLatency(data)
	default:
		sorter = byID(data)
	}
	switch order {
	case sortx.OrderAscending:
		return sorter
	case sortx.OrderDescending:
		return sort.Reverse(sorter)
	default:
		return sorter
	}
}

// byNextRun is a wrapper for sorting the entry array by next run time.
// (with zero time at the end).
type byNextRun []StatusData

func (s byNextRun) Len() int      { return len(s) }
func (s byNextRun) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byNextRun) Less(i, j int) bool {
	// Two zero times should return false.
	// Otherwise, zero is "greater" than any other time.
	// (To sort it at the end of the list.)
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
}

// byPrevRun is a wrapper for sorting the entry array by prev run time.
// (with zero time at the end).
type byPrevRun []StatusData

func (s byPrevRun) Len() int      { return len(s) }
func (s byPrevRun) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byPrevRun) Less(i, j int) bool {
	// Two zero times should return false.
	// Otherwise, zero is "greater" than any other time.
	// (To sort it at the end of the list.)
	if s[i].Prev.IsZero() {
		return false
	}
	if s[j].Prev.IsZero() {
		return true
	}
	return s[i].Prev.Before(s[j].Prev)
}

// byLatency is a wrapper for sorting the entry array by latency.
type byLatency []StatusData

func (s byLatency) Len() int      { return len(s) }
func (s byLatency) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byLatency) Less(i, j int) bool {
	return s[i].Job.latency < s[j].Job.latency
}

// byStatus is a wrapper for sorting the entry array by status.
type byStatus []StatusData

func (s byStatus) Len() int      { return len(s) }
func (s byStatus) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byStatus) Less(i, j int) bool {
	return s[i].Job.status < s[j].Job.status
}

// byID is a wrapper for sorting the entry array by entry id.
type byID []StatusData

func (s byID) Len() int      { return len(s) }
func (s byID) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byID) Less(i, j int) bool {
	return s[i].Job.EntryID < s[j].Job.EntryID
}

// byName is a wrapper for sorting the entry array by name.
type byName []StatusData

func (s byName) Len() int      { return len(s) }
func (s byName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool {
	t := strings.Map(unicode.ToUpper, s[i].Job.Name)
	u := strings.Map(unicode.ToUpper, s[j].Job.Name)
	return t < u
}
