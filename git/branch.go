package git

import (
	"fmt"
	"strconv"
	"time"
)

// Branch defines a git branch
type Branch struct {
	Name          string
	LastCommit    time.Time
	LastCommitter string
}

func parseUnixTimestamp(str string) time.Time {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(i, 0)
}

func (b *Branch) formatDate() string {
	if b.LastCommit.IsZero() {
		return ""
	}
	year, month, day := b.LastCommit.Date()
	return fmt.Sprintf("%v-%02v-%02v", year, int(month), day)
}
