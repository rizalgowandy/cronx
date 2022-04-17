package cronx

import (
	"sort"
	"strings"
)

type SortOrder int64

const (
	SortOrderUndefined SortOrder = iota
	SortOrderAscending
	SortOrderDescending
)

type SortKey string

const (
	SortKeyID      SortKey = "id"
	SortKeyStatus  SortKey = "status"
	SortKeyPrevRun SortKey = "prev_run"
	SortKeyNextRun SortKey = "next_run"
	SortKeyLatency SortKey = "latency"
)

func NewStatusDataSorter(key SortKey, order SortOrder, data []StatusData) sort.Interface {
	var sorter sort.Interface
	switch key {
	case SortKeyID:
		sorter = byID(data)
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
	case SortOrderUndefined, SortOrderAscending:
		return sorter
	case SortOrderDescending:
		return sort.Reverse(sorter)
	default:
		return sorter
	}
}

type Sort struct {
	Key   SortKey
	Order SortOrder
}

// NewSorts create sorting based on
// Format:
//	sort=key1:asc,key2:desc,key3:asc
func NewSorts(qs string) []Sort {
	sorts := strings.Split(qs, ",")

	var res []Sort
	for _, v := range sorts {
		kv := strings.Split(v, ":")

		s := Sort{
			Key:   SortKey(kv[0]),
			Order: SortOrderAscending,
		}
		if len(kv) == 2 {
			switch kv[1] {
			case "asc":
				s.Order = SortOrderAscending
			case "desc":
				s.Order = SortOrderDescending
			default:
				s.Order = SortOrderUndefined
			}
		}
		res = append(res, s)
	}

	return res
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
