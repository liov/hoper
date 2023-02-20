package main

import "log"

func Scan() {
	var id int
	err := db.Table(ModelTable).Select("MAX(id)").Scan(&id).Error
	if err != nil {
		log.Println(err)
	}
	log.Println(id)
}

func RawScan() {
	var exists bool
	err := db.Raw(`SELECT EXISTS(SELECT * FROM tsp_info WHERE id = 59)`).Scan(&exists).Error
	if err != nil {
		log.Println(err)
	}
	log.Println(exists)
}
