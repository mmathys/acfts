package model

import (
	"github.com/jinzhu/gorm"
)

// Output represents transaction output
type Output struct {
	gorm.Model
	Amount       int         `json:"amount"`
	Address1     string      `json:"address1"`
	Address2     string      `json:"address2"`
	PreviousHash string      `json:"previous_hash"`
	Index        uint        `json:"index" gorm:"column:output_index"`
	Used         bool        `json:"used"`
	Signatures   []Signature `json:"server_signatures"`
}
