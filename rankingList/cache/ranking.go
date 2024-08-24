package cache

import "fmt"

const (
	DailyRankKey = "redis-test-rank:daily"
)

func ShareKey(id string) string {
	return fmt.Sprintf("redis-test-share:%s", id)
}
