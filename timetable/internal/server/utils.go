package server

import (
	"time"

	models "github.com/gnom48/hospital-api-lib"
)

func IsUserInRole(roles []models.Role, roleId string) bool {
	for _, role := range roles {
		if role.Id == roleId {
			return true
		}
	}
	return false
}

func IsValidTime(from, to time.Time) bool {
	if from.After(to) {
		return false
	}

	if to.Sub(from) > 12*time.Hour {
		return false
	}

	if from.Second() != 0 || from.Minute()%30 != 0 || to.Second() != 0 || to.Minute()%30 != 0 {
		return false
	}

	return true
}
