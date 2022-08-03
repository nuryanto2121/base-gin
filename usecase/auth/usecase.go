package useauth

import (
	"app/models"
	"fmt"
)

func genRole(userRole []*models.UserRoleDesc) ([]string, []*models.OutletList) {
	var (
		roles   = make([]string, len(userRole))
		outlets = []*models.OutletList{}
	)

	for i, val := range userRole {
		dtOutlet := val.Outlets
		fmt.Println(dtOutlet)
		roles[i] = val.Role
	}
	return roles, outlets
}
