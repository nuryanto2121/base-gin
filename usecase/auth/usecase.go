package useauth

import (
	"app/models"
	"app/pkg/logging"
	"context"
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
		// if err := json.Unmarshal([]byte(val.Outlets), &dtOutlet); err != nil {
		// 	fmt.Errorf("unMarshal ", err)
		// }
		fmt.Println(dtOutlet)
		roles = val.Role
	}
	if len(dtOutlet) > 0 {
		dd := dtOutlet[0]["outlet_id"]
		return roles, fmt.Sprintf("%s", dd)
	}
	return roles, ""
}
func (u *useAuht) genOutletList(ctx context.Context, userId string) ([]*models.OutletLookUp, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		logger  = logging.Logger{}
		outlets = []*models.OutletLookUp{}
	)

	query := models.ParamList{
		Page:       1,
		PerPage:    1000,
		SortField:  "outlet_name",
		InitSearch: fmt.Sprintf("a.user_id='%s' ", userId),
	}
	outlets, err := u.repoRoleOutlet.GetList(ctx, query)
	if err != nil {
		logger.Error("error get role outlet list ", err)
	}

	return outlets, nil
}
