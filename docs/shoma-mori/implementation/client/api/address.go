package api

import (
	"acfts-client/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAddrs is
func GetAddrs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clients := []model.Client{}
		// ch <- true
		// FIXME: clusetr idがリセットされないので1とは限らない
		db.Where("cluster_id = ?", 1).Find(&clients)
		// <-ch

		addrs := make([]model.Address, len(clients))
		for i, client := range clients {
			addrs[i] = client.Address
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "success",
			"addresses": addrs,
		})
	}
}
