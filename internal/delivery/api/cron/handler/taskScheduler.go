package handler

import (
	"log"

	"gorm.io/gorm"
)

func get_hostname(db *gorm.DB) {
	ips := []string{}
	db.Table("servers").
		Select("ip as IP").
		Where("active = true").
		Scan(&ips)

	log.Printf("Active IPs: %+v\n", ips)
}
