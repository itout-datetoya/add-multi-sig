package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"

	"multisigservice/db"
	"multisigservice/models"
)

// challengeStore はアドレス毎のチャレンジ（nonce）を一時保存するストアです。
// 実運用では有効期限やキャッシュクリア機構の実装が望まれます。
var challengeStore = struct {
	sync.RWMutex
	m map[string]string // key: Ethereum address (小文字), value: challenge
}{m: make(map[string]string)}

// ChallengeHandler は指定アドレスに対しランダムなチャレンジを発行します。
func ChallengeHandler(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "address is required"})
		return
	}

	// ランダムなnonceを生成（シンプルな例）
	nonceBytes := make([]byte, 16)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate challenge"})
		return
	}
	challenge := hex.EncodeToString(nonceBytes)

	// ストアに保存
	challengeStore.Lock()
	challengeStore.m[strings.ToLower(address)] = challenge
	challengeStore.Unlock()

	c.JSON(http.StatusOK, gin.H{"challenge": challenge})
}

// LoginRequest はログイン時のリクエストデータです。
type LoginRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

// LoginHandler はチャレンジ署名方式により署名検証を行い、ユーザーをDBに登録します。
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Address == "" || req.Signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	// ストアからチャレンジを取得
	challengeStore.RLock()
	challenge, exists := challengeStore.m[strings.ToLower(req.Address)]
	challengeStore.RUnlock()

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No challenge found for address"})
		return
	}

	// 署名検証
	valid, err := verifySignature(challenge, req.Signature, req.Address)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Signature verification failed"})
		return
	}

	// 認証成功の場合、ユーザーをDBに登録（既存の場合は更新）
	user := models.User{Address: req.Address}
	// GORMのSaveはプライマリキーに基づいて更新・作成を行う
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database error"})
		return
	}

	// 使用済みチャレンジを削除
	challengeStore.Lock()
	delete(challengeStore.m, strings.ToLower(req.Address))
	challengeStore.Unlock()

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logged in", "address": req.Address})
}

// verifySignature は、チャレンジメッセージと署名から署名者のアドレスが一致するか検証します。
// Ethereumのpersonal_signでは、メッセージの先頭に定型文字列が付加されます。
func verifySignature(message, signatureHex, expectedAddress string) (bool, error) {
	// 署名はhex文字列なのでデコード
	sig, err := hex.DecodeString(strings.TrimPrefix(signatureHex, "0x"))
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}
	// リカバリIDの補正（27,28に合わせる）
	if sig[64] != 27 && sig[64] != 28 {
		sig[64] += 27
	}

	// Ethereum仕様に基づくメッセージの前処理
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	msg := []byte(prefix + message)

	// Keccak256ハッシュの計算
	hash := crypto.Keccak256(msg)

	// 署名から公開鍵を復元
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %v", err)
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()

	if strings.ToLower(recoveredAddr) != strings.ToLower(expectedAddress) {
		return false, nil
	}

	return true, nil
}
