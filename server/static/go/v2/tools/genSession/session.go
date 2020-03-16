package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

type Session struct {
	UserID       int    `json:"userId"`
	UserName     string `json:"userName"`
	UserRealName string `json:"userRealName"`
	ClientIp     string `json:"clientIp"`
}

type ErpSession struct {
	PiId          int      `json:"piId"`
	EmployeeId    int      `json:"employeeId"`
	EmployeeName  string   `json:"employeeName"`
	DeptId        int      `json:"deptId"`
	DeptName      string   `json:"deptName"`
	CompId        int      `json:"compId"`
	CompName      string   `json:"compName"`
	EnglishName   string   `json:"englishName"`
	Phone         string   `json:"phone"`
	RoleCodeList  []string `json:"roleCodeList"`
	Type          int      `json:"type"`
	FilterIds     []int    `json:"filterIds"`
	FilterDeptIds []int    `json:"filterDeptIds"`
}

func main() {
	sess := &Session{UserID: 699}
	data, _ := json.Marshal(sess)
	log.Println("X-Session-Data", base64.StdEncoding.EncodeToString(data))
}
