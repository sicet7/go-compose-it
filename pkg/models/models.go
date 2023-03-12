package models

// Register database models here.
var (
	databaseModels = map[string]interface{}{
		"user": User{},
	}
)

func Get() map[string]interface{} {
	return databaseModels
}
