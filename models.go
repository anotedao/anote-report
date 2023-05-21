package main

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Address string `gorm:"size:255;uniqueIndex"`
	Balance uint64
	New     uint64
}

type KeyValue struct {
	gorm.Model
	Key      string `gorm:"size:255;uniqueIndex"`
	ValueInt uint64 `gorm:"type:int"`
	ValueStr string `gorm:"type:string"`
}
