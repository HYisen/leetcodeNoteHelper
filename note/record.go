package note

import (
	"fmt"
	"leetcodeNoteHelper/date"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	ProblemID  int
	Difficulty Difficulty
	Simple     bool
	Begin      time.Time
	End        time.Time
}

func NewRecord(lines []string, d date.Date) (*Record, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("too short %d lines", len(lines))
	}

	problemID, err := parseHeaderLine(lines[0])
	if err != nil {
		return nil, fmt.Errorf("parse header line %s: %v", lines[0], err)
	}

	ret, err := parseFirstContentLine(lines[1], d)
	if err != nil {
		return nil, fmt.Errorf("parse first content line %s: %v", lines[1], err)
	}
	ret.ProblemID = problemID
	if ret.End != (time.Time{}) {
		return ret, nil
	}

	begin, end, err := findPair(lines[2:], d)
	if err != nil {
		return nil, err
	}
	if begin != (time.Time{}) {
		ret.Begin = begin
	}
	ret.End = end
	return ret, nil
}

func fixJustPastMidnight(t time.Time) time.Time {
	// Skip zero value as it shall be an empty value.
	if t == (time.Time{}) {
		return time.Time{}
	}
	if date.BelongToYesterday(t) {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func parseHeaderLine(s string) (problemID int, err error) {
	const headerLinePrefix = "## "
	if !strings.HasPrefix(s, headerLinePrefix) {
		return 0, fmt.Errorf("bad header line not begin with %s", headerLinePrefix)
	}
	id, err := strconv.Atoi(s[len(headerLinePrefix):])
	if err != nil {
		return 0, fmt.Errorf("parse problem ID in header line: %v", err)
	}
	return id, nil
}

// parseFirstContentLine try its best to part the first content line.
// If simple as one line content, returns a partial Record (only Difficulty, Simple, Begin, End).
// If may not simple, returns a smaller partial Record (only Difficulty, Simple, Begin).
// So shall you check partial to decide whether to fulfill End later.
// If facing error that could be neither simple nor not, returns not nil err.
// Those optional results are not named return values because too many.
func parseFirstContentLine(s string, d date.Date) (partial *Record, err error) {
	parts := strings.Split(s, " ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("bad line parts")
	}
	difficulty, err := parseDifficulty(parts[0])
	if err != nil {
		return nil, fmt.Errorf("parse difficulty part: %v", err)
	}
	begin, err := parseTime(parts[1], d)
	if err != nil {
		return nil, fmt.Errorf("parse begin %s: %v", parts[1], err)
	}
	begin = fixJustPastMidnight(begin)
	if len(parts) >= 3 {
		end, err := parseTime(parts[2], d)
		// If not ok, may the following content would make it valid, so just silent ignore without return here.
		if err == nil {
			end = fixJustPastMidnight(end)
			return &Record{
				Difficulty: difficulty,
				Simple:     true,
				Begin:      begin,
				End:        end,
			}, nil
		}
	}
	return &Record{
		Difficulty: difficulty,
		Simple:     false,
		Begin:      begin,
	}, nil
}

// findPair tries its best to find the pair, begin may be zero as it's optional in content.
func findPair(rest []string, d date.Date) (begin, end time.Time, err error) {
	for _, line := range rest {
		parts := strings.Split(line, " ")
		if len(parts) == 0 {
			continue
		}
		parts = append(parts, "")
		one, oneErr := parseTime(parts[0], d)
		two, twoErr := parseTime(parts[1], d)
		if twoErr == nil && oneErr == nil {
			begin = fixJustPastMidnight(one)
			end = fixJustPastMidnight(two)
		} else if twoErr != nil && oneErr == nil {
			end = fixJustPastMidnight(one)
		}
	}
	if end == (time.Time{}) {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to find end in %v", rest)
	}
	return begin, end, nil
}

func parseTime(s string, d date.Date) (time.Time, error) {
	if len(s) != 4 {
		return time.Time{}, fmt.Errorf("not 4 width")
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		return time.Time{}, fmt.Errorf("not number: %v", err)
	}
	hour := num / 100
	minute := num % 100
	return time.Date(d.Year, d.Month, d.Day, hour, minute, 0, 0, time.Local), nil
}
