package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"multisigservice/db"
	"multisigservice/models"

	"gorm.io/datatypes"
	"github.com/google/uuid"
	"time"
)

// CreateMultiSigHandler は、ログインユーザー（Owner）と2名の参加者からマルチシグを作成しDBに登録します。
func CreateMultiSigHandler(c *gin.Context) {
	var req struct {
		Owner        string   `json:"owner"`        // ログイン済みのユーザーアドレス
		Participants []string `json:"participants"` // 参加者のEthereumアドレス（2名）
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Participants) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Require owner and exactly 2 participants"})
		return
	}

	// マルチシグIDを生成し、初期状態を設定
	newMultiSig := models.MultiSig{
		Owner:        req.Owner,
		Participants: datatypes.JSON([]byte(mustMarshal(req.Participants))),
		Status:       "awaiting",
		Data:         datatypes.JSON([]byte(`{}`)),
	}
	newMultiSig.ID = uuid.New().String()

	if err := db.DB.Create(&newMultiSig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error on create"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "MultiSig created", "multisig": newMultiSig})
}

// GetMultiSigListHandler は、指定ユーザーが参加しているマルチシグの一覧を返します。
// 簡単のため全件取得していますが、実際はフィルタ処理を実装してください。
func GetMultiSigListHandler(c *gin.Context) {
	var multisigs []models.MultiSig
	if err := db.DB.Find(&multisigs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error on list"})
		return
	}
	c.JSON(http.StatusOK, multisigs)
}

// GetMultiSigDataHandler は、指定マルチシグの署名用データ（例としてプレースホルダー）を生成し返します。
func GetMultiSigDataHandler(c *gin.Context) {
	id := c.Param("id")
	var ms models.MultiSig
	if err := db.DB.First(&ms, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "MultiSig not found"})
		return
	}
	// 状態に応じたデータ（ここではシンプルにタイムスタンプ付きの文字列を例示）
	dataToSign := "data-placeholder-" + time.Now().String()
	dataJSON, _ := json.Marshal(map[string]string{"dataToSign": dataToSign})
	ms.Data = datatypes.JSON(dataJSON)

	// DB上も更新
	db.DB.Save(&ms)
	c.JSON(http.StatusOK, gin.H{"dataToSign": dataToSign})
}

// UpdateMultiSigDataHandler は、クライアントから送信された署名情報をDB上のマルチシグに反映し状態更新します。
func UpdateMultiSigDataHandler(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		Signature string `json:"signature"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signature payload"})
		return
	}

	var ms models.MultiSig
	if err := db.DB.First(&ms, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "MultiSig not found"})
		return
	}

	// 受信した署名情報をJSON形式で保存（必要に応じたプロトコルロジックの実装を行う）
	dataMap := map[string]interface{}{}
	// 既存のDataを読み込む
	json.Unmarshal(ms.Data, &dataMap)
	dataMap["signature"] = payload.Signature
	updatedData, _ := json.Marshal(dataMap)
	ms.Data = datatypes.JSON(updatedData)

	// 状態更新（例：一方の署名を受領した場合）
	ms.Status = "partial"

	if err := db.DB.Save(&ms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error on update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signature data updated", "multisig": ms})
}

// mustMarshal は、JSON変換に失敗した場合にpanicする簡易関数です。
func mustMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
