package useauth

import (
	"context"
	"strings"
	"time"

	iauth "app/interface/auth"
	ifileupload "app/interface/fileupload"
	iroleoutlet "app/interface/role_outlet"
	itrx "app/interface/trx"
	iusers "app/interface/user"
	iuserapps "app/interface/user_apps"
	iuserrole "app/interface/user_role"
	iusersession "app/interface/user_session"
	"app/models"

	"app/pkg/logging"
	"app/pkg/setting"
	util "app/pkg/util"
)

type useAuht struct {
	repoAuth        iusers.Repository
	repoFile        ifileupload.Repository
	repoUserSession iusersession.Repository
	repoUserRole    iuserrole.Repository
	repoRoleOutlet  iroleoutlet.Repository
	repoUserApps    iuserapps.Repository
	repoTrx         itrx.Repository
	contextTimeOut  time.Duration
}

func NewUserAuth(
	repoAuth iusers.Repository, repoFile ifileupload.Repository,
	repoUserSession iusersession.Repository, repoUserRole iuserrole.Repository,
	repoRoleOutlet iroleoutlet.Repository, repoUserApps iuserapps.Repository,
	repoTrx itrx.Repository, timeout time.Duration,
) iauth.Usecase {
	return &useAuht{
		repoAuth:        repoAuth,
		repoFile:        repoFile,
		repoUserSession: repoUserSession,
		repoUserRole:    repoUserRole,
		repoRoleOutlet:  repoRoleOutlet,
		repoUserApps:    repoUserApps,
		repoTrx:         repoTrx,
		contextTimeOut:  timeout,
	}
}

func (u *useAuht) LoginCms(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		logger   = logging.Logger{}
		dataUser = &models.Users{}
		role     string
		outlets  = []*models.OutletLookUp{}
	)

	dataUser, err = u.repoAuth.GetByAccount(ctx, dataLogin.Account)
	if err != nil {
		logger.Error("error usecase.LoginCms().GetByAccount ", err)
		return nil, models.ErrUnauthorized
	}

	if !util.ComparePassword(dataUser.Password, util.GetPassword(dataLogin.Password)) {
		logger.Error("error usecase.LoginCms().ComparePassword ")
		return nil, models.ErrInvalidPassword
	}

	if !dataUser.IsActive {
		return nil, models.ErrAccountNotActive
	}

	//get outlet
	outlets, err = u.genOutletList(ctx, dataUser.Id.String())
	if dataLogin.Account != "root" {
		roles, err := u.repoUserRole.GetListByUser(ctx, "user_id", dataUser.Id.String())
		if err != nil {
			return nil, err
		}
		role = roles[0].Role
	} else {
		role = "root"
	}

	token, err := util.GenerateToken(dataUser.Id.String(), dataUser.Username, role)
	if err != nil {
		return nil, err
	}

	//save to session
	exDate := util.GetTimeNow().Add(time.Duration(setting.AppSetting.ExpiredJwt) * time.Hour)

	dataSession := &models.UserSession{
		UserId:      dataUser.Id,
		Token:       token,
		ExpiredDate: exDate,
	}
	err = u.repoUserSession.Create(ctx, dataSession)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{}
	if len(outlets) == 0 {
		response = map[string]interface{}{
			"users":   dataUser,
			"token":   token,
			"role":    role,
			"outlets": nil,
		}
	} else {
		response = map[string]interface{}{
			"users":   dataUser,
			"token":   token,
			"role":    role,
			"outlets": outlets[0],
		}
	}

	return response, nil
}

func (u *useAuht) LoginMobile(ctx context.Context, req *models.LoginForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		// logger = logging.Logger{}
		token string = ""
	)

	userApps, err := u.repoUserApps.GetByAccount(ctx, req.Account)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, models.ErrAccountNotFound
		}
		return nil, err
	}

	if !util.ComparePassword(userApps.Password, util.GetPassword(req.Password)) {
		return nil, models.ErrInvalidPassword
	}

	token, err = util.GenerateToken(userApps.Id.String(), userApps.Name, "user")
	if err != nil {
		return nil, err
	}

	//save to session
	exDate := util.GetTimeNow().Add(time.Duration(setting.AppSetting.ExpiredJwt) * time.Hour)
	dataSession := &models.UserSession{
		UserId:      userApps.Id,
		Token:       token,
		ExpiredDate: exDate,
	}
	err = u.repoUserSession.Create(ctx, dataSession)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"user_id":  userApps.Id,
		"token":    token,
		"name":     userApps.Name,
		"phone_no": userApps.PhoneNo,
	}
	return response, nil
	// }

}

func (u *useAuht) ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)

	defer cancel()

	// dataUser, err := fb.GetUserByAccount(ctx, util.NameStruct(models.Users{}), dataForgot.Account)
	// if err != nil {
	// 	return err
	// }

	// GenCode := util.GenerateNumber(6)

	// // send generate code
	// mailService := &useemailauth.Register{
	// 	Email:      dataUser.Email,
	// 	Name:       dataUser.Name,
	// 	GenerateNo: GenCode,
	// }

	// go mailService.SendRegister()
	// // if err != nil {
	// // 	return err
	// // }

	// //store to redis
	// err = redisdb.AddSession(dataUser.Email, GenCode, 5*time.Minute)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (u *useAuht) ResetPassword(ctx context.Context, dataReset *models.ResetPasswd) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)

	defer cancel()
	var (
		logger   = logging.Logger{}
		dataUser = &models.Users{}
	)

	if dataReset.Passwd != dataReset.ConfirmPasswd {
		return models.ErrWrongPasswordConfirm
	}

	// account, err := util.DecryptMessage(dataReset.Account)
	// if err != nil {
	// 	return models.ErrInternalServerError
	// }

	dataUser, err = u.repoAuth.GetByAccount(ctx, dataReset.Account)
	if err != nil {
		logger.Error("error usecase.LoginCms().GetByAccount ", err)
		return err
	}

	dataUser.Password, _ = util.Hash(dataReset.Passwd)
	dtUpdate := map[string]interface{}{
		"password": dataUser.Password,
	}

	err = u.repoAuth.Update(ctx, dataUser.Id, dtUpdate)
	if err != nil {
		logger.Error("error usecase.ResetPassword().Update ", err)
		return err
	}
	return nil
}

func (u *useAuht) Register(ctx context.Context, req models.RegisterForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		now    = util.GetTimeNow()
		logger = logging.Logger{}
	)

	if req.ConfirmasiPassword != req.Password {
		return models.ErrWrongPasswordConfirm
	}
	// req.Password, _ = util.HashAndSalt()
	req.Password, _ = util.Hash(req.Password)

	userApps, err := u.repoUserApps.GetDataBy(ctx, "phone_no", req.PhoneNo)
	if err != nil && err != models.ErrNotFound {
		return err
	}
	if userApps != nil {
		return models.ErrAccountAlreadyExist
	}

	errTx := u.repoTrx.Run(ctx, func(trxCtx context.Context) error {
		parent := &models.UserApps{
			AddUserApps: models.AddUserApps{
				Name:     req.Name,
				PhoneNo:  req.PhoneNo,
				IsParent: true,
				Password: req.Password,
				JoinDate: now,
			},
		}
		err := u.repoUserApps.Create(trxCtx, parent)
		if err != nil {
			logger.Error("error create parent ", err)
			return err
		}

		for _, val := range req.Childs {
			child := &models.UserApps{
				AddUserApps: models.AddUserApps{
					Name:     val.Name,
					ParentId: parent.Id,
					IsParent: false,
					JoinDate: now,
					DOB:      val.DOB,
				},
			}

			err := u.repoUserApps.Create(trxCtx, child)
			if err != nil {
				logger.Error("error create child ", err)
				return err
			}
		}
		return nil
	})

	return errTx
}

func (u *useAuht) Verify(ctx context.Context, dataVerify models.VerifyForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	response := map[string]interface{}{
		"user_id": "dataUser.ID",
		// "token":    token,
		// "email":    dataUser.Email,
		// "name":     dataUser.Name,
		// "avatar":   dataUser.Avatar,
		// "phone_no": dataUser.PhoneNo,
	}

	return response, nil
}
func (u *useAuht) VerifyForgot(ctx context.Context, dataVerify models.VerifyForgotForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)

	defer cancel()
	var dataUser = &models.Users{}

	if dataVerify.AccessToken != "" {

		// dataToken, err := fb.VerifyIDToken(ctx, dataVerify.AccessToken)
		// if err != nil {
		// 	return nil, err
		// }
		// expiredDate := util.Int64ToTime(dataToken.Expires)
		// if expiredDate.Before(util.GetTimeNow()) {
		// 	return nil, models.ErrExpiredFirebaseToken
		// }

		// dataUser, err = fb.GetUserByAccount(ctx, util.NameStruct(models.Users{}), dataVerify.PhoneNo)
		// if err != nil {
		// 	return nil, err
		// }

	} else {
		//validasi otp
		// if dataVerify.Otp == "" {
		// 	return nil, models.ErrOtpNotFound
		// }
		// existOTP := redisdb.GetSession(dataVerify.Email)

		// if existOTP != dataVerify.Otp {
		// 	return nil, models.ErrInvalidOTP
		// }

		// dataUser, err = fb.GetUserByAccount(ctx, util.NameStruct(models.Users{}), dataVerify.Email)
		// if err != nil {
		// 	return nil, err
		// }
	}

	account, _ := util.EncryptMessage(dataUser.Email)

	response := map[string]interface{}{
		"account": account,
	}

	return response, nil
}

func (u *useAuht) Logout(ctx context.Context, Claims util.Claims, Token string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)

	defer cancel()

	var (
	// conditions []*firebase.FirestoreConditions
	// userFcm    = &models.UserFcm{}
	)

	token := strings.Split(Token, "Bearer ")[1]

	u.repoUserSession.Delete(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (u *useAuht) CheckPhoneNo(ctx context.Context, PhoneNo string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)

	defer cancel()

	// dataUser, err := fb.GetUserByAccount(ctx, util.NameStruct(models.Users{}), PhoneNo)
	// if err != nil && err != models.ErrAccountNotFound {
	// 	return models.ErrInternalServerError
	// }

	// if dataUser != nil {
	// 	return models.ErrPhoneNoConflict
	// }

	return nil
}
