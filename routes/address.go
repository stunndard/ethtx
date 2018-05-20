package routes

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/kataras/iris.v6"
)

func GenerateAddress(ctx *iris.Context) {

	privKey, err := crypto.GenerateKey()
	if err != nil {
		log.Println("error generating private key", err)
		renderJSONError(ctx, "101", "error generating key")
		return
	}

	//pubKey :=
	address := crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	log.Println("key successfully generated")

	privKey1 := fmt.Sprintf("0x%x", privKey.D)
	if RedactLogs {
		log.Println("private key: 0x...")
	} else {
		log.Println("private key:", privKey1)
	}

	pubKey1 := fmt.Sprintf("0x%x, 0x%x", privKey.PublicKey.X, privKey.PublicKey.Y)
	log.Println("public  key:", pubKey1)
	log.Println("address    :", address)

	ctx.JSON(iris.StatusOK, &iris.Map{
		"privKey": privKey1,
		"pubKey": pubKey1,
		"address": address,
	})
}
