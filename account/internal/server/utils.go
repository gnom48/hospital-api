package server

import models "github.com/gnom48/hospital-api-lib"

func IsUserInRole(roles []models.Role, roleId string) bool {
	for _, role := range roles {
		if role.Id == roleId {
			return true
		}
	}
	return false
}
