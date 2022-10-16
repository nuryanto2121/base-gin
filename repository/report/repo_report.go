package reporeport

import (
	ireport "app/interface/report"
	"app/models"
	"app/pkg/db"
	"app/pkg/logging"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type repoReport struct {
	db db.DBGormDelegate
}

func NewRepoReport(Conn db.DBGormDelegate) ireport.Repository {
	return &repoReport{Conn}
}

// GetReport implements ireport.Repository
func (r *repoReport) GetReport(ctx context.Context, sWhere, startDate, endDate, userId string) ([]*models.ReportResponse, error) {
	var (
		logger = logging.Logger{}
		conn   = r.db.Get(ctx)
		result = []*models.ReportResponse{}
	)

	query := fmt.Sprintf(`
	select o.outlet_name
		,o.outlet_city 
		,t.id as transaction_id
		,td.ticket_no 
		,t.transaction_code
		,t.transaction_date
		,CASE 
			WHEN t.payment_code ='3000001' then 'BCA'
			WHEN t.payment_code ='3000002' then 'CC / Kartu Kredit'
			WHEN t.payment_code ='3000003' then 'Tunai'
			WHEN t.payment_code ='3000004' then 'Kasir'
			WHEN t.payment_code ='3000005' then 'BCA VA'
			WHEN t.payment_code ='3000006' then 'Online Payment'
			else ''
		end as payment_method
		,CASE 
			WHEN t.status_transaction ='2000001' then 'Booked'
			WHEN t.status_transaction ='2000002' then 'Check In'
			WHEN t.status_transaction ='2000003' then 'Check Out'
			WHEN t.status_transaction ='2000004' then 'Draf'
			WHEN t.status_transaction ='2000005' then 'Active'
			WHEN t.status_transaction ='2000006' then 'Finish'
			when t.status_transaction = 2000007 then 'Delta'
			when t.status_transaction = 2000008 then 'Overtime'
			else ''
		end as status_transaction
		,td.is_overtime
		,td.is_overtime_paid
		,case 
			when sm.is_bracelet = true then concat(sm.sku_name,' - ',sm.duration,' Jam' )
			else sm.sku_name 
		end as sku_name
		,case when sm.is_bracelet = true then 1 else  0 end as num_of_kids_booked
		,case 
			when sm.is_bracelet = true 
			then (case when date(td.check_in) = date('0001-01-01') then 0 else 1 end)
			else  0 
		 end as num_of_kids_check_in
		,case 
			when sm.is_bracelet = true 
			then (case when date(td.check_in) = date('0001-01-01') then 1 else 0 end)
			else  0 
		 end as delta_booked_vs_check_in
		,t.total_ticket as total_booked
		,td.product_qty as qty
		,td.price as amount
		,td.amount as total_amount		
	FROM outlets o 
	inner join role_outlet ro
		on o.id = ro.outlet_id
	inner join "transaction" t
		on o.id = t.outlet_id 
	inner  join transaction_detail td 
		on t.id = td.transaction_id 
	left join sku_management sm 
		on sm.id =td.product_id 
	WHERE ro.user_id= ? 	
	and  (date(t.transaction_date) between ? and ?)
	%s
	order by o.outlet_city,o.outlet_name,t.transaction_date,t.transaction_code  
	`, sWhere)

	err := conn.Raw(query, userId, startDate, endDate).Scan(&result).Error
	if err != nil {
		logger.Error("repo order GetList ", err)
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
