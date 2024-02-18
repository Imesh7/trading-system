package balance

import "time"

type Balance struct {
	statusTime time.Time
	amount int64
}