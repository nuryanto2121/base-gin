package useauth

import (
	"app/models"
	"encoding/json"
	"fmt"
)

func genRole(userRole []*models.UserRoleDesc) (string, string) {
	var (
		roles string
		// outlets  = []*models.OutletLookUp{}
		dtOutlet []map[string]interface{}
	)

	for _, val := range userRole {
		// dtOutlet := val.Outlets
		if err := json.Unmarshal([]byte(val.Outlets), &dtOutlet); err != nil {
			fmt.Errorf("unMarshal ", err)
		}
		fmt.Println(dtOutlet)
		roles = val.Role
	}
	if len(dtOutlet) > 0 {
		dd := dtOutlet[0]["outlet_id"]
		return roles, fmt.Sprintf("%s", dd)
	}
	return roles, ""
}
