package model

import "github.com/jinzhu/gorm"

// Signature represents signatures of servers
type Signature struct {
	gorm.Model
	Address
	Signature1 string `json:"signature1"`
	Signature2 string `json:"signature2"`
	OutputID   uint
}
