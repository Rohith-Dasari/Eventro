package authorisation

import (
	"eventro/models"
	"strings"
)

// signup
func AssignRole(email string) models.Role {
	if strings.HasSuffix(email, "@eadmin.com") {
		return models.Admin
	} else if strings.HasSuffix(email, "@ehost.com") {
		return models.Host
	}
	return models.Customer
}
