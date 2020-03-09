package model

import "github.com/jinzhu/gorm"

// Output represents transaction output
type Output struct {
	gorm.Model
	Address
	Amount       int         `json:"amount"`
	PreviousHash string      `json:"previous_hash"`
	Index        uint        `json:"index" gorm:"column:output_index"`
	Used         bool        `json:"used"`
	Signatures   []Signature `json:"server_signatures" gorm:"foreignkey:OutputID"`
}
