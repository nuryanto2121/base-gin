package useauth

import (
	"context"
	"strings"
	"time"

	iauth "app/interface/auth"
	ifileupload "app/interface/fileupload"
	iusers "app/interface/user"
	iusersession "app/interface/user_session"
	"app/models"

	"app/pkg/logging"
	"app/pkg/redisdb"
	"app/pkg/setting"
	util "app/pkg/utils"
)

type useAuht struct {
	repoAuth        iusers.Repository
	repoFile        ifileupload.Repository
	repoUserSession iusersession.Repository
	contextTimeOut  time.Duration
}

func NewUserAuth(a iusers.Repository, b ifileupload.Repository, c iusersession.Repository, timeout time.Duration) iauth.Usecase {
	return &useAuht{
		repoAuth:        a,
		repoFile:        b,
		repoUserSession: c,
		contextTimeOut:  timeout,
	}
}

func (u *useAuht) LoginCms(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)

	defer cancel()

	var (
		logger   = logging.Logger{}
		dataUser = &models.Users{}
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

	//get user group

	token, err := util.GenerateToken(dataUser.Id.String(), dataUser.Username, "")
	if err != nil {
		return nil, err
	}

	//save to session
	dataSession := &models.UserSession{
		UserId:      dataUser.Id,
		Token:       token,
		ExpiredDate: time.Now().Add(time.Duration(setting.AppSetting.ExpiredJwt) * time.Hour),
	}
	err = u.repoUserSession.Create(ctx, dataSession)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"users": dataUser,
		"token": token,
	}
	return response, nil
}

func (u *useAuht) LoginMobile(ctx context.Context, dataLogin *models.SosmedForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	// var (
	// 	dataUser = &models.Users{}
	// 	fb       = firebase.InitFirebase(ctx)
	// 	token    (string)
	// 	response map[string]interface{}
	// )
	// defer fb.Close()

	// dataToken, err := fb.VerifyIDToken(ctx, dataLogin.AccessToken)
	// if err != nil {
	// 	return nil, err
	// }
	// expiredDate := util.Int64ToTime(dataToken.Expires)
	// if expiredDate.Before(util.GetTimeNow()) {
	// 	return nil, models.ErrExpiredFirebaseToken
	// }

	// dataUser, _ = fb.GetUserByAccount(ctx, util.NameStruct(models.Users{}), dataLogin.Email)
	// if dataUser != nil {
	// 	if dataUser.PhoneNo == "" && !dataUser.IsActive {
	// 		response = map[string]interface{}{
	// 			"user_id":  dataUser.ID,
	// 			"token":    token,
	// 			"email":    dataUser.Email,
	// 			"name":     dataUser.Name,
	// 			"avatar":   dataUser.Avatar,
	// 			"phone_no": dataUser.PhoneNo,
	// 		}
	// 		return response, nil
	// 	}

	// 	token, err = util.GenerateToken(dataUser.ID, dataUser.Name, "")
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	redisdb.AddSession(token, dataUser.ID, time.Duration(setting.AppSetting.ExpiredJwt)*time.Hour)

	// 	response = map[string]interface{}{
	// 		"user_id":  dataUser.ID,
	// 		"token":    token,
	// 		"email":    dataUser.Email,
	// 		"name":     dataUser.Name,
	// 		"avatar":   dataUser.Avatar,
	// 		"phone_no": dataUser.PhoneNo,
	// 	}
	// 	return response, nil
	// } else {
	// 	var User models.Users
	// 	User.CreatedBy = dataLogin.Name
	// 	User.UpdatedBy = dataLogin.Name
	// 	User.IsActive = false
	// 	User.Email = dataLogin.Email
	// 	User.Name = dataLogin.Name
	// 	User.SosmedID = dataToken.UID

	// 	// err = u.repoAuth.Create(ctx, &User)
	// 	rest, _, err := fb.Create(ctx, util.NameStruct(User), User)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	response := map[string]interface{}{
		// "user_id":  rest.ID,
		// "token":    token,
		// "email":    dataLogin.Email,
		// "name":     dataLogin.Name,
		"phone_no": "",
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

	if dataReset.Passwd != dataReset.ConfirmPasswd {
		return models.ErrWrongPasswordConfirm
	}

	// account, err := util.DecryptMessage(dataReset.Account)
	// if err != nil {
	// 	return models.ErrInternalServerError
	// }

	// dataUser, err := fb.GetUserByAccount(ctx, util.NameStruct(models.Users{}), account)
	// if err != nil {
	// 	return err
	// }

	// dataUser.Password, _ = util.Hash(dataReset.Passwd)
	// dtUpdate := map[string]interface{}{
	// 	"password": dataUser.Password,
	// }

	// _, err = fb.Update(ctx, util.NameStruct(models.Users{}), dataUser.ID, dtUpdate)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (u *useAuht) Register(ctx context.Context, dataRegister models.RegisterForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	// fb := firebase.InitFirebase(ctx)
	// var User models.Users
	// defer fb.Close()
	defer cancel()

	// dataUser, err := fb.GetUserByAccount(ctx, util.NameStruct(User), dataRegister.Email)
	// if err != nil && err != models.ErrAccountNotFound {
	// 	return models.ErrInternalServerError
	// }
	// //check duplicate email
	// if dataUser != nil && dataUser.Email != "" {
	// 	return models.ErrAccountAlreadyExist
	// }

	// err = mapstructure.Decode(dataRegister, &User.AddUser)
	// if err != nil {
	// 	return err
	// }

	// User.Password, _ = util.Hash(dataRegister.Password)
	// User.CreatedBy = dataRegister.Name
	// User.UpdatedBy = dataRegister.Name
	// User.IsActive = false
	// User.CreatedAt = util.UnixNow()
	// User.UpdatedAt = util.UnixNow()

	// dtt, _, err := fb.Create(ctx, util.NameStruct(User), User)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("\nKey == %v\n", dtt.ID)
	// defer fb.Close()
	return nil
}

func (u *useAuht) Verify(ctx context.Context, dataVerify models.VerifyForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	// fb := firebase.InitFirebase(ctx)

	// defer fb.Close()
	defer cancel()
	// dataToken, err := fb.VerifyIDToken(ctx, dataVerify.AccessToken)
	// if err != nil {
	// 	return nil, err
	// }
	// expiredDate := util.Int64ToTime(dataToken.Expires)
	// if expiredDate.Before(util.GetTimeNow()) {
	// 	return nil, models.ErrExpiredFirebaseToken
	// }

	// //validasi phone no

	// dataUser, err := fb.GetUserByAccount(ctx, util.NameStruct(models.Users{}), dataVerify.Email)
	// if err != nil {
	// 	return nil, models.ErrAccountNotFound
	// }

	// dtUpdate := map[string]interface{}{
	// 	"is_active": true,
	// 	"phone_no":  dataVerify.PhoneNo,
	// 	"join_date": util.UnixNow(),
	// }

	// _, err = fb.Update(ctx, util.NameStruct(models.Users{}), dataUser.ID, dtUpdate)
	// if err != nil {
	// 	return nil, models.ErrInternalServerError
	// }

	// token, err := util.GenerateToken(dataUser.ID, dataUser.Name, "")
	// if err != nil {
	// 	return nil, err
	// }

	// redisdb.AddSession(token, dataUser.ID, time.Duration(24)*time.Hour)
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
	err = redisdb.TurncateList(token)
	if err != nil {
		return err
	}

	// conditions = append(conditions, &firebase.FirestoreConditions{
	// 	Field:    "device_token",
	// 	Operator: "==",
	// 	Value:    token,
	// })

	// conditions = append(conditions, &firebase.FirestoreConditions{
	// 	Field:    "user_id",
	// 	Operator: "==",
	// 	Value:    Claims.UserID,
	// })

	// _, err = fb.ReadByMultiKey(ctx, util.NameStruct(models.UserFcm{}), conditions).Next()
	// if err != nil {
	// 	return err
	// }

	// _, err = fb.Delete(ctx, util.NameStruct(models.UserFcm{}), docs.Ref.ID)
	// if err != nil {
	// 	return err
	// }

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
