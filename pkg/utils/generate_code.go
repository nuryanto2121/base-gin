package util

import (
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"app/pkg/setting"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GenCode(t *models.TmpCode) string {
	var (
		ctx = context.Background()
		// conn          = postgres.Conn
		code   string = ""
		logger        = logging.Logger{}
		// tmpCode = models.TmpCode{}
	)
	dbConn := db.NewDBdelegate(setting.DatabaseSetting.Debug)
	dbConn.Init()

	conn := dbConn.Get(ctx)
	// declare needed variable
	tm := time.Now()

	years := fmt.Sprintf("%d", tm.Year())
	month := tm.Month()

	// make roman of month
	roman := map[time.Month]string{time.January: "I", time.February: "II", time.March: "III", time.April: "IV",
		time.May: "V", time.June: "VI", time.July: "VII", time.August: "VIII", time.September: "IX", time.October: "X",
		time.November: "XI", time.December: "XII"}

	var isNewPrefix bool = false
	if err := conn.Model(models.TmpCode{}).Where("prefix = ?", t.Prefix).First(t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			isNewPrefix = true
		}
	}

	if isNewPrefix {
		// generate new code, cause last code not found
		t.SeqNo = 1
		code = fmt.Sprintf("%s/%s/%s/%04d", t.Prefix, years, roman[month], t.SeqNo)
		t.Prefix = t.Prefix
		t.Code = code

		if err := conn.Create(t).Error; err != nil {
			logger.Error("error create tmp_code  ", err)
		}
	} else {
		splitCode := strings.Split(t.Code, "/")
		// if code is different in year or month, suffix will be 001
		if splitCode[2] != roman[month] || splitCode[1] != years {
			t.SeqNo = 1
			code = fmt.Sprintf("%s/%s/%s/%04d", t.Prefix, years, roman[month], t.SeqNo)
		} else {
			// increment suffix +1
			t.SeqNo += 1
			code = fmt.Sprintf("%s/%s/%s/%04d", t.Prefix, years, roman[month], t.SeqNo)
		}
		t.Code = code
		if err := conn.Model(models.TmpCode{}).Where("id=?", t.ID).Updates(t).Error; err != nil {
			logger.Error("error update tmp_code  ", err)
		}
	}

	return code
}
