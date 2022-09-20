package main

import "time"

type Model struct {
	Id        int
	A         int
	B         int
	C         int
	D         int
	E         int
	F         string
	H         string
	I         string
	J         string
	K         string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

const TableNameTspSettleInfo = `tsp_settle_info`

type TspSettleInfo struct {
	Id                       int       `json:"id" explain:""`
	TspId                    int       `json:"tspId" explain:"结算对象Id"`
	ContractEndDay           time.Time `json:"contractEndDay" explain:"合同到期日"`
	SettleMethod             int       `json:"settleMethod" explain:"结算方式, 0:预付款抵扣/日结 1:现款/日结 2:现款/月结"`
	HourOfDay                int       `json:"hourOfDay" explain:"订座票归属日的起始小时数, 0:0点,6:6点"`
	DayBillAutoConfirm       int       `json:"dayBillAutoConfirm" explain:"日账单自动收款确认, 0:否 1:是"`
	MonthBillAutoConfirm     int       `json:"monthBillAutoConfirm" explain:"月账单自动申请对账, 0:否 1:是"`
	EffectiveStartDay        time.Time `json:"effectiveStartDay" explain:"结算方式生效起始日期"`
	PrepaymentWarningBalance int       `json:"prepaymentWarningBalance" explain:"预付款余额告警值, 结算方式为预付日结时有效"`
	PiId                     string    `json:"piId" explain:"erp的流程Id, 修改时更新当前流程Id"`
	UpdatedAt                time.Time `json:"updatedAt" explain:"修改时间"`
	CreatedAt                time.Time `json:"createdAt" explain:"创建时间"`
	Status                   int       `json:"status" explain:"状态"`
}

type TspSettleInfo2 struct {
	Id                       int       `json:"id" explain:""`
	TspId                    int       `json:"tspId" explain:"结算对象Id"`
	ContractEndDay           string    `json:"contractEndDay" explain:"合同到期日"`
	SettleMethod             int       `json:"settleMethod" explain:"结算方式, 0:预付款抵扣/日结 1:现款/日结 2:现款/月结"`
	HourOfDay                int       `json:"hourOfDay" explain:"订座票归属日的起始小时数, 0:0点,6:6点"`
	DayBillAutoConfirm       int       `json:"dayBillAutoConfirm" explain:"日账单自动收款确认, 0:否 1:是"`
	MonthBillAutoConfirm     int       `json:"monthBillAutoConfirm" explain:"月账单自动申请对账, 0:否 1:是"`
	EffectiveStartDay        string    `json:"effectiveStartDay" explain:"结算方式生效起始日期"`
	PrepaymentWarningBalance int       `json:"prepaymentWarningBalance" explain:"预付款余额告警值, 结算方式为预付日结时有效"`
	PiId                     string    `json:"piId" explain:"erp的流程Id, 修改时更新当前流程Id"`
	UpdatedAt                time.Time `json:"updatedAt" explain:"修改时间"`
	CreatedAt                time.Time `json:"createdAt" explain:"创建时间"`
	Status                   int       `json:"status" explain:"状态"`
}
