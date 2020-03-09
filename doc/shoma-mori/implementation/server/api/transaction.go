package api

import (
	"acfts/db/model"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func createSignature(transaction model.Transaction) (*big.Int, *big.Int) {
	// Remove gorm.Model, Used, Signatures from outputs
	simpleOutputs := make([]simpleOutput, len(transaction.Outputs))
	for i := range simpleOutputs {
		simpleOutputs[i] = convertOutput(transaction.Outputs[i])
	}

	// Convert []simpleOutput to bytes to get its hash
	buf := []byte(fmt.Sprintf("%v", simpleOutputs))

	// Get hash using SHA256
	h := crypto.Hash.New(crypto.SHA256)
	h.Write(buf)
	hashed := h.Sum(nil)

	fmt.Println("hash when creating")
	fmt.Println(hashed)

	// Get signature using ellipse curve cryptography
	r, s, err := ecdsa.Sign(rand.Reader, key, hashed)
	if err != nil {
		panic(err)
	}

	return r, s
}

func unlockUTXO(utxo model.Output, signature1, signature2 string) bool {
	// Convert transaction struct to bytes to get its hash
	buf := []byte(fmt.Sprintf("%v%v%v", utxo.Address1, utxo.Address2, utxo.PreviousHash))

	// Get a hash value using SHA256
	h := crypto.Hash.New(crypto.SHA256)
	h.Write(buf)
	hashed := h.Sum(nil)

	// Convert string to big.Int
	var ok bool
	n1 := new(big.Int)
	n1, ok = n1.SetString(signature1, 10)
	if !ok {
		fmt.Println("Signature1 is not valid.")
		return false
	}

	n2 := new(big.Int)
	n2, ok = n2.SetString(signature2, 10)
	if !ok {
		fmt.Println("Signature2 is not valid.")
		return false
	}

	// model.Output.Address1 and Address2 represent a public key
	// Convert Address1 and Address2 to ecdsa.PublicKey
	address1, ok1 := new(big.Int).SetString(utxo.Address1, 10)
	address2, ok2 := new(big.Int).SetString(utxo.Address2, 10)
	if !ok1 || !ok2 {
		log.Println("Error: an address of a server is not valid.")
		return false
	}
	publicKey := ecdsa.PublicKey{elliptic.P521(), address1, address2}

	if ecdsa.Verify(&publicKey, hashed, n1, n2) {
		fmt.Println("Verifyed!")
		return true
	}
	return false
}

func verifyUTXO(utxo model.Output, siblings []model.Output, n int) bool {
	// Genesis is approved without verification
	if utxo.PreviousHash == "genesis" {
		fmt.Println("genesis is approved without verification.")
		return true
	}

	// Create the same array as when creating a signature
	outputs := make([]simpleOutput, 1+len(siblings))
	for i := range outputs {
		if i < int(utxo.Index) {
			outputs[i] = convertOutput(siblings[i])
		} else if i > int(utxo.Index) {
			outputs[i] = convertOutput(siblings[i-1])
		} else {
			outputs[i] = convertOutput(utxo)
		}
	}

	// Convert transaction struct to bytes to get its hash
	buf := []byte(fmt.Sprintf("%v", outputs))

	// Get a hash value using SHA256
	h := crypto.Hash.New(crypto.SHA256)
	h.Write(buf)
	hashed := h.Sum(nil)

	fmt.Println("hash when verification")
	fmt.Println(hashed)

	valid := 0
	for _, signature := range utxo.Signatures {
		// FIXME: Should get public keys of other servers independently of clients
		address1, ok1 := new(big.Int).SetString(signature.Address1, 10)
		address2, ok2 := new(big.Int).SetString(signature.Address2, 10)
		if !ok1 || !ok2 {
			log.Println("Error: an address of a server is not valid")
			return false
		}
		serverPubKey := ecdsa.PublicKey{elliptic.P521(), address1, address2}
		signature1, ok1 := new(big.Int).SetString(signature.Signature1, 10)
		signature2, ok2 := new(big.Int).SetString(signature.Signature2, 10)
		if !ok1 || !ok2 {
			log.Println("Error: a signature of server is not valid")
			return false
		}
		if ecdsa.Verify(&serverPubKey, hashed, signature1, signature2) {
			valid++
			fmt.Println("one valid signature")
			if float64(valid) >= 2.0*float64(n)/3.0 {
				fmt.Println("Server Verifyed!")
				return true
			}
		} else {
			fmt.Println("not valid")
		}
	}
	return false
}

// Remove unnecessary properties from model.Input
type simpleInput struct {
	UTXO       simpleOutput `json:"utxo"`
	Signature1 string       `json:"signature1"`
	Signature2 string       `json:"signature2"`
}

// Remove unnecessary properties from model.Input
type simpleOutput struct {
	Amount       int    `json:"amount"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	PreviousHash string `json:"previous_hash"`
	Index        uint   `json:"index"`
}

// Remove unnecessary properties from model.Input
type simpleTransaction struct {
	Inputs  []simpleInput  `json:"inputs"`
	Outputs []simpleOutput `json:"outputs"`
}

// Convert model.Output to simpleOutput
func convertOutput(output model.Output) simpleOutput {
	json := simpleOutput{}
	json.Amount = output.Amount
	json.Address1 = output.Address1
	json.Address2 = output.Address2
	json.PreviousHash = output.PreviousHash
	json.Index = output.Index
	return json
}

// Convert model.Transaction to simpleTransaction
func convertTransaction(transaction model.Transaction) simpleTransaction {
	json := simpleTransaction{}

	json.Inputs = make([]simpleInput, len(transaction.Inputs))
	for i, input := range transaction.Inputs {
		utxo := input.UTXO
		json.Inputs[i].UTXO = convertOutput(utxo)
		json.Inputs[i].Signature1 = input.Signature1
		json.Inputs[i].Signature2 = input.Signature2
	}

	json.Outputs = make([]simpleOutput, len(transaction.Outputs))
	for i, output := range transaction.Outputs {
		json.Outputs[i] = convertOutput(output)
	}

	return json
}

// VerifyTransaction is
func VerifyTransaction(db *gorm.DB, n int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction model.Transaction
		c.BindJSON(&transaction)

		inputs := transaction.Inputs
		if len(inputs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There are no inputs.",
			})
			return
		}

		outputs := transaction.Outputs
		if len(outputs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There are no outputs.",
			})
			return
		}

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Transacton error. Rollback.",
				})
				return
			}
		}()

		inputAmount := 0
		for _, input := range inputs {
			utxo := input.UTXO
			count := 0
			tx.Where("address1 = ? AND address2 = ? AND previous_hash = ? AND output_index = ?",
				utxo.Address1, utxo.Address2, utxo.PreviousHash, utxo.Index).
				First(&utxo).Count(&count)
			if count == 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Input is not valid.",
				})
				tx.Rollback()
				return
			}
			if utxo.Used {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "UTXO is used.",
				})
				tx.Rollback()
				return
			}
			/* When choosing to delete records of outputs
			which have already been used. */
			// count := 0
			// db.First(&utxo).Count(&count)
			// if count == 0 {
			// 	c.JSON(http.StatusBadRequest, gin.H{
			// 		"message": "There are no valid UTXOs.",
			// 	})
			// 	return
			// }

			if !unlockUTXO(utxo, input.Signature1, input.Signature2) {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Could not unlock UTXO.",
				})
				tx.Rollback()
				return
			}

			// Update gorm.Model of Siblings
			tx.Where("id <> ? AND previous_hash = ?", utxo.ID, utxo.PreviousHash).
				Find(&input.Siblings)
			if !verifyUTXO(utxo, input.Siblings, n) {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "One of signatures of servers is not valid.",
				})
				tx.Rollback()
				return
			}

			inputAmount += utxo.Amount
		}

		outputAmount := 0
		for _, output := range outputs {
			outputAmount += output.Amount
		}

		if inputAmount != outputAmount {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Amount of inputs is different from amount of outputs.",
			})
			tx.Rollback()
			return
		}

		for i, input := range inputs {
			utxo := input.UTXO
			// db.Model(&utxo).Where("address1 = ? AND address2 = ? AND previous_hash = ?",
			// 	utxo.Address1, utxo.Address2, utxo.PreviousHash).Update("used", true)
			// db.Unscoped().Delete(&utxo)

			tx.Where("address1 = ? AND address2 = ? AND previous_hash = ? AND output_index = ?",
				utxo.Address1, utxo.Address2, utxo.PreviousHash, utxo.Index).
				First(&utxo).Update("used", true)
			transaction.Inputs[i].UTXO = utxo
		}
		tx.Commit()

		// Add records of outputs
		for i, output := range outputs {
			db.Create(&output)
			transaction.Outputs[i] = output
		}

		r, s := createSignature(transaction)
		simpleTx := convertTransaction(transaction)

		c.JSON(http.StatusCreated, gin.H{
			"message":     "Verified this transaction.",
			"transaction": simpleTx,
			"address1":    (&key.PublicKey).X.String(),
			"address2":    (&key.PublicKey).Y.String(),
			"signature1":  r,
			"signature2":  s,
		})
	}
}
