package model

import "github.com/jinzhu/gorm"

// Cluster is
type Cluster struct {
	gorm.Model
	URL     string
	Clients []Client `gorm:"foreignkey:ClusterID"`
}
