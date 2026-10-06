package currentuser

import "sync"

const currentUserID = 1

var (
	once     sync.Once
	instance int
)

func GetCurrentUserID() int {
	once.Do(func() {
		instance = currentUserID
	})

	return instance
}
