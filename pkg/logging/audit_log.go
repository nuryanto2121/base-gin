package logging

import (

	// "gitlab.com/369-engineer/369backend/account/pkg/monggodb"

	"time"
)

type AuditLog struct {
	AuditLogID int64     `json:"audit_log_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Level      string    `json:"level"`
	UUID       string    `json:"uuid"`
	FuncName   string    `json:"func_name"`
	FileName   string    `json:"file_name"`
	Line       int       `json:"line"`
	Time       string    `json:"time"`
	Message    string    `json:"message"`
}

func (a *AuditLog) saveAudit() {

	// ctx := context.Background()
	// idx := strconv.FormatInt(util.UnixNow(), 10)
	// a.CreatedAt = util.GetTimeNow()
	// a.UpdatedAt = util.GetTimeNow()
	// val, err := elasticsearch.Client.Index().
	// 	Index("audit_log_chat").
	// 	Id(idx).
	// 	BodyJson(a).
	// 	Do(ctx)

	// if err != nil {
	// 	fmt.Printf("\nerror === >%#v\n", err)
	// }
	// fmt.Printf("\value === >%#v\n", val)
	// a.ID = util.GetTimeNow().Unix()
	// result, err := monggodb.MCon.Collection("auditlogs").InsertOne(context.TODO(), a)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// fmt.Println("Inserted a single document: ", result.InsertedID)

}
