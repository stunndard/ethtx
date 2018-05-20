package routes

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/kataras/iris.v6"
)

func CheckTx(ctx *iris.Context) {

	// check chainId
	chainId, ok := new(big.Int).SetString(ctx.PostValue("chainId"), 10)
	if !ok {
		renderJSONError(ctx, "102", "invalid chainId")
		log.Println("invalid chainId", ctx.PostValue("chainId"))
		return
	}

	// check privKey
	s := ctx.PostValue("privKey")
	if len(s) != 64 {
		renderJSONError(ctx, "102", "invalid privKey length")
		log.Println("invalid privKey length", s)
		return
	}
	if _, err := hex.DecodeString(s); err != nil {
		renderJSONError(ctx, "102", "invalid privKey")
		log.Println("invalid privKey", s)
		return
	}
	senderPrivKey, err := crypto.HexToECDSA(s)
	if err != nil {
		renderJSONError(ctx, "102", "invalid privKey")
		log.Println("invalid privKey", s, err)
		return
	}

	// check sendTo
	s = ctx.PostValue("sendTo")
	if !common.IsHexAddress(s) {
		renderJSONError(ctx, "102", "invalid sendTo address")
		log.Println("invalid sendTo address", s)
		return
	}
	sendTo := common.HexToAddress(s)

	// check nonce
	s = ctx.PostValue("nonce")
	nonce, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		renderJSONError(ctx, "102", "invalid nonce")
		log.Println("invalid nonce", s, err)
		return
	}

	// check amount
	s = ctx.PostValue("amount")
	amountF, ok := new(big.Float).SetString(s)
	if !ok {
		renderJSONError(ctx, "102", "invalid amount")
		log.Println("invalid amount", s)
		return
	}
	multipierF, _ := new(big.Float).SetString("1000000000000000000")
	amountF = amountF.Mul(amountF, multipierF) //
	amount := new(big.Int)
	amountF.Int(amount)

	// check if there's also amount in wei
	// if there is, it has priority
	s = ctx.PostValue("amountWei")
	amountWei := new(big.Int)
	if s != "" {
		amountWei, ok = amountWei.SetString(s, 10)
		if !ok {
			renderJSONError(ctx, "102", "invalid amountWei")
			log.Println("invalid amountWei", s)
			return
		}
		amount = amountWei
	}

	// check gasLimit
	s = ctx.PostValue("gasLimit")
	gasLimit, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		renderJSONError(ctx, "102", "invalid gasLimit")
		log.Println("invalid gasLimit", s)
		return
	}

	// check gasPrice
	s = ctx.PostValue("gasPrice")
	gasPrice, ok := new(big.Int).SetString(s, 10)
	if !ok {
		renderJSONError(ctx, "102", "invalid gasPrice")
		log.Println("invalid gasPrice", s)
		return
	}

	// save all checked params to ctx
	ctx.Set("chainId", chainId)
	ctx.Set("senderPrivKey", senderPrivKey)
	ctx.Set("sendTo", sendTo)
	ctx.Set("nonce", nonce)
	ctx.Set("amount", amount)
	ctx.Set("gasLimit", gasLimit)
	ctx.Set("gasPrice", gasPrice)

	ctx.Next()
}

func SignTx(ctx *iris.Context) {
	chainId := ctx.Get("chainId").(*big.Int)
	senderPrivKey := ctx.Get("senderPrivKey").(*ecdsa.PrivateKey)
	sendTo := ctx.Get("sendTo").(common.Address)
	nonce := ctx.Get("nonce").(uint64)
	amount := ctx.Get("amount").(*big.Int)
	gasLimit := ctx.Get("gasLimit").(uint64)
	gasPrice := ctx.Get("gasPrice").(*big.Int)

	log.Println("signing tx")
	log.Println("chainId  :", chainId.String())
	if RedactLogs {
		log.Println("privKey  : 0x...")
	} else {
		log.Println("privKey  :", fmt.Sprintf("0x%x", senderPrivKey.D))
	}
	log.Println("sendTo   :", sendTo.String())
	log.Println("nonce    :", nonce)
	log.Println("amountWei:", amount.String())
	log.Println("gasLimit :", gasLimit)
	log.Println("gasPrice :", gasPrice)

	tx := types.NewTransaction(nonce, sendTo, amount, gasLimit, gasPrice, nil)

	signer := types.NewEIP155Signer(chainId)
	signedTx, err := types.SignTx(tx, signer, senderPrivKey)
	if err != nil {
		renderJSONError(ctx, "102", "error signing transaction")
		log.Println("error signing transaction", err)
		return
	}
	log.Println("tx successfully signed")
	//log.Println(signedTx)

	var buff bytes.Buffer
	if err := signedTx.EncodeRLP(&buff); err != nil {
		renderJSONError(ctx, "102", "error encoding rlp transaction")
		log.Println("error encode rlp transaction", err)
		return
	}
	sTx := fmt.Sprintf("0x%x", buff.Bytes())
	log.Println("signedTx :", sTx)


	ctx.JSON(iris.StatusOK, &iris.Map{
		"result": "ok",
		"signedTx": sTx,
	})
}