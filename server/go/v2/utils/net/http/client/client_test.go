package client

import (
	"net/http"
	"testing"
	"time"
)

type OperatorReq struct {
	Id           uint64    `json:"id" explain:"运营商id" example:"1"`
	Ids          []uint64  `json:"ids" explain:"运营商id" example:"[1,2]"`
	Name         string    `json:"name" explain:"运营商名称" example:"xxx"`
	StartTime    time.Time `json:"startTime" explain:"起始时间"`
	EndTime      time.Time `json:"endTime" explain:"结束时间"`
	PageNo       int       `json:"pageNo" explain:"页数" example:"1"`
	PageNum      int       `json:"pageNum" explain:"每页数量" example:"20"`
	PlatformType uint8     `json:"platformType" explain:"平台类型,必传,传255为所有状态"`
}

type OperatorPublicInfo struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name" explain:"'公司名称"`
	ShortName string `json:"shortName" explain:"运营商名称/公司简称"`
}

type OperatorPublicList struct {
	List  []OperatorPublicInfo `json:"list"`
	Total uint64               `json:"total"`
}

func TestClient(t *testing.T){
	req := &OperatorReq{
		PageNo:       1,
		PageNum:      10,
		PlatformType: 1,
	}
	res:=&OperatorPublicList{}
	err := NewRequest(`http://operator-center.openmng/api/operator/allList/v1`,http.MethodPost,req).
		SetHeader("Erp_User_Para","e30=").HTTPRequest(CommonResponse(res))
	if err!=nil{
		t.Fatal(err)
	}
	t.Log(res)
}
