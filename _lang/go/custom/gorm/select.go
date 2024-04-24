package main

import "log"

func Select() {
	var models []*Model
	err := db.Table(ModelTable).Find(&models).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models {
		log.Println(model)
	}
}

func SelectStrTime() {
	var models []*ModelA
	err := db.Table(ModelTable).Find(&models).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models {
		log.Println(model.K)
	}
}
