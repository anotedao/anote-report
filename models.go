package main

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Address string `gorm:"size:255;uniqueIndex"`
	Balance uint64
}
