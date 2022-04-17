package cronx

// byNextRun is a wrapper for sorting the entry array by next run time
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
