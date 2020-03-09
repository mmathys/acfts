package api

import (
	"acfts-client/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ReceiveUTXO is
func ReceiveUTXO(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// log.Println("received utxos")
		j := struct {
			UTXOs []model.Output `json:"outputs"`
		}{}
		c.BindJSON(&j)

		// bytes, _ := json.MarshalIndent(j, "", "    ")
		// log.Println("utxos")
		// log.Println(string(bytes))

		// FIXME: 1つ1つのUTXOの署名を検証する
		for _, utxo := range j.UTXOs {
			db.Create(&utxo)
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Received a utxo.",
			"utxos":   j.UTXOs,
		})
	}
}

// ClearOutputs is
func ClearOutputs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		output := model.Output{}
		db.DropTable(&output)
		db.AutoMigrate(&output)

		c.JSON(http.StatusOK, gin.H{
			"message": "all data is deleted.",
		})
	}
}
