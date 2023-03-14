package entities

import "github.com/sicet7/go-compose-it/src/database/entities/user"

// Register entities here
var entities = []interface{}{
	&user.User{},
}

func List() []interface{} {
	return entities
}
