package account

import "time"

type Account struct {
	UserId    int64
	Amount    int64
	CreatedAt time.Time
}
