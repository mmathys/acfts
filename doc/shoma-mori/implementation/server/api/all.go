package api

import (
	"acfts/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeleteAll is
func DeleteAll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		output := model.Output{}
		db.DropTable(&output)
		db.AutoMigrate(&output)

		c.JSON(http.StatusOK, gin.H{
			"message": "all data is deleted.",
		})
	}
}
