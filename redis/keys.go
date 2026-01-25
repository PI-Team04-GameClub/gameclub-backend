package redis

import "fmt"

const (
	KeyGameAll     = "game:all"
	KeyGameByID    = "game:id:%s"
	KeyUserByID    = "user:id:%s"
	KeyUserByEmail = "user:email:%s"
)

func GameByIDKey(id string) string {
	return fmt.Sprintf(KeyGameByID, id)
}

func GameByIDKeyUint(id uint) string {
	return fmt.Sprintf(KeyGameByID, fmt.Sprint(id))
}

func UserByIDKey(id string) string {
	return fmt.Sprintf(KeyUserByID, id)
}

func UserByEmailKey(email string) string {
	return fmt.Sprintf(KeyUserByEmail, email)
}
