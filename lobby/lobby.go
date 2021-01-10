package lobby

import (
	"catans/user"
)

type Lobby struct {
	users map[string]*user.User
}
