package cronx

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
