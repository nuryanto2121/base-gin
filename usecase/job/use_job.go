package usejob

import (
	ijob "app/interface/job"
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
	"time"
)

type useJob struct {
	repoTransaction itransaction.Repository
	repoTransDetail itransactiondetail.Repository
	repoTrx         itrx.Repository
	useSMS          isms.Usecase
	contextTimeOut  time.Duration
}

func NewUseJob(
	a itransaction.Repository,
	a1 itransactiondetail.Repository,
	b isms.Usecase,
	c itrx.Repository,
	timeout time.Duration,
) ijob.Usecase {
	return &useJob{
		repoTransaction: a,
		repoTransDetail: a1,
		useSMS:          b,
		repoTrx:         c,
		contextTimeOut:  timeout,
	}
}

// NotificationSms implements ijob.Usecase
func (u *useJob) NotificationSms(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		logger = logging.Logger{}
		// wg     sync.WaitGroup
		// sLock      = sync.RWMutex{}
		paramQuery = models.ParamList{}
	)
	logger.Info("job NotificationSms")

	sQuery := fmt.Sprintf("td.flag_notif_send = 0 and date(td.check_in) = '%s' and td.is_children = true", util.GetTimeNow().Format("2006-01-02"))
	paramQuery = models.ParamList{
		Page:       1,
		PerPage:    1000000,
		SortField:  "",
		InitSearch: sQuery,
	}

	facilityDetail, err := u.repoTransaction.GetList(ctx, paramQuery)
	if err != nil {
		logger.Error("cron error getlist transaction detail ", err)
		return err
	}

	dataUpdate := map[string]interface{}{
		"flag_notif_send": 1, //prosess
	}

	// sQuery = strings.ReplaceAll(sQuery, "td.", "")

	// err = u.repoTransDetail.UpdateBy(ctx, sQuery, dataUpdate)
	// if err != nil {
	// 	logger.Error("cron error update transaction detail ", err)
	// 	return err
	// }

	for _, val := range facilityDetail {
		// wg.Add(1)
		fmt.Println(val)
		// val.CheckIn
		//check if lama bermain kurang 15mnt lagi then send message
		now := util.GetTimeNow()
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
		if timeLeft <= float64(setting.AppSetting.MinSendNotif) {
			//send sms
			message := fmt.Sprintf("Pelanggan yang terhormat, waktu bermain ananda %s tersisa %f menit ", val.Name, timeLeft)
			err = u.useSMS.Send(ctx, val.PhoneNo, message, "")
			if err != nil {
				// ret
				dataUpdate["flag_notif_send"] = 2 //error or failed send sms
			} else {
				dataUpdate["flag_notif_send"] = 3 //terkirim
			}
		} else {
			dataUpdate["flag_notif_send"] = 0 //belum terkirim

		}
		sQuery = fmt.Sprintf("ticket_no = '%s'", val.TicketNo)

		err = u.repoTransDetail.UpdateBy(ctx, sQuery, dataUpdate)
		if err != nil {
			logger.Error("cron error update transaction detail ", err)
			// return err
		}

	}

	return nil
}
