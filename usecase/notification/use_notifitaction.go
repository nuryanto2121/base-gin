package usenotification

import (
	inotification "app/interface/notification"
	isms "app/interface/sms"
	itransaction "app/interface/transaction"
	itransactiondetail "app/interface/transaction_detail"
	itrx "app/interface/trx"
	"app/models"
	"app/pkg/logging"
	"app/pkg/setting"
	"app/pkg/util"
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

type useNotification struct {
	repoTransaction itransaction.Repository
	repoTransDetail itransactiondetail.Repository
	repoTrx         itrx.Repository
	useSMS          isms.Usecase
	contextTimeOut  time.Duration
}

func NewUseNotification(
	a itransaction.Repository,
	a1 itransactiondetail.Repository,
	b isms.Usecase,
	c itrx.Repository,
	timeout time.Duration,
) inotification.Usecase {
	return &useNotification{
		repoTransaction: a,
		repoTransDetail: a1,
		useSMS:          b,
		repoTrx:         c,
		contextTimeOut:  timeout,
	}
}

// NotificationSms implements inotification.Usecase
func (u *useNotification) NotificationSms(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		logger     = logging.Logger{}
		wg         sync.WaitGroup
		sLock      = sync.RWMutex{}
		paramQuery = models.ParamList{}
	)
	logger.Info("job NotificationSms ", util.GetTimeNow())

	sQuery := fmt.Sprintf("td.flag_notif_send = %d and date(td.check_in) = '%s' and td.is_children = true", models.SMS_NOT_READY, util.GetTimeNow().Format("2006-01-02"))
	paramQuery = models.ParamList{
		Page:       1,
		PerPage:    1000000,
		SortField:  "",
		InitSearch: sQuery,
	}

	facilityDetail, err := u.repoTransaction.GetList(ctx, paramQuery)
	if err != nil {
		logger.Error("cron error getlist transaction detail ", err)
		// return err
	}

	dataUpdate := map[string]interface{}{
		"flag_notif_send": models.SMS_PROCESS, //prosess
	}

	sQuery = strings.ReplaceAll(sQuery, "td.", "")

	err = u.repoTransDetail.UpdateBy(ctx, sQuery, dataUpdate)
	if err != nil {
		logger.Error("cron error update transaction detail ", err)
		// return err
	}

	for _, val := range facilityDetail {
		wg.Add(1)
		fmt.Println(val)
		go u.newMethod(ctx, &wg, val, dataUpdate, &sLock)

	}
	wg.Wait()
	logger.Info("job NotificationSms finish ", util.GetTimeNow())
	// return nil
}

func (u *useNotification) newMethod(ctx context.Context, wg *sync.WaitGroup, val *models.TransactionList, dataUpdate map[string]interface{}, sLock *sync.RWMutex) {
	defer wg.Done()

	var (
		now    = util.GetTimeNow()
		logger = logging.Logger{}
		sQuery string
	)
	checkOut := val.CheckOut
	if checkOut.IsZero() {
		checkOut = val.CheckIn.Add(time.Duration(val.Duration) * time.Hour)
	}
	fmt.Println("now")
	fmt.Println(now)
	fmt.Println("val.CheckIn")
	fmt.Println(val.CheckIn)
	fmt.Println("checkOut")
	fmt.Println(checkOut)
	diff := checkOut.Sub(now).Minutes()
	fmt.Println(diff)
	timeLeft := math.Round(diff)
	//check if lama bermain kurang 15mnt lagi then send message
	if timeLeft <= float64(setting.AppSetting.MinSendNotif) {

		message := fmt.Sprintf("Pelanggan yang terhormat, waktu bermain ananda %s tersisa %d menit ", val.Name, int(timeLeft))
		err := u.useSMS.Send(ctx, val.ParentId, val.PhoneNo, message, "")
		if err != nil {

			logger.Error("failed send sms ", err)
			sLock.Lock()
			dataUpdate["flag_notif_send"] = models.SMS_FAILED //error or failed send sms
			sLock.Unlock()
		} else {
			sLock.Lock()
			dataUpdate["flag_notif_send"] = models.SMS_SUCCESS //send sms
			sLock.Unlock()

		}
	} else {
		sLock.Lock()
		dataUpdate["flag_notif_send"] = models.SMS_NOT_READY //belum terkirim
		sLock.Unlock()

	}
	sQuery = fmt.Sprintf("ticket_no = '%s'", val.TicketNo)

	err := u.repoTransDetail.UpdateBy(ctx, sQuery, dataUpdate)
	if err != nil {
		logger.Error("cron error update transaction detail ", err)

	}
	logger.Info("send sms success")
}
