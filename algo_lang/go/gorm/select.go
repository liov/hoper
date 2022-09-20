package main

import "log"

func Select() {
	var tspSettleInfos []*TspSettleInfo
	err := db.Table(TableNameTspSettleInfo).Find(&tspSettleInfos).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, tspSettleInfo := range tspSettleInfos {
		log.Println(tspSettleInfo)
	}
}

func SelectStrTime() {
	var tspSettleInfos []*TspSettleInfo2
	err := db.Table(TableNameTspSettleInfo).Find(&tspSettleInfos).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, tspSettleInfo := range tspSettleInfos {
		log.Println(tspSettleInfo.ContractEndDay)
	}
}
