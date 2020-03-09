package api

import (
	"acfts/db/model"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// key is a private key of this server
var key *ecdsa.PrivateKey

// CreateGenesis makes a genesis transaction
func CreateGenesis(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		genesis := model.Output{}
		db.Where("previous_hash = ?", "genesis").First(&genesis)
		if db.NewRecord(genesis) {
			var err error
			key, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
			if err != nil {
				panic(err)
			}

			var output model.Output
			c.BindJSON(&output)

			genesis = model.Output{
				Amount:       output.Amount,
				Address1:     output.Address1,
				Address2:     output.Address2,
				PreviousHash: "genesis",
				Index:        0,
				Used:         false,
			}
			db.Create(&genesis)

			json := convertOutput(genesis)
			c.JSON(http.StatusCreated, gin.H{
				"message": "Genesis is created.",
				"genesis": json,
			})
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Genesis has already created.",
				"genesis": genesis,
			})
		}
	}
}
