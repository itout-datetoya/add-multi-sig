package handlers

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/gin-gonic/gin"
	"multisigservice/db"
	"multisigservice/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// CreateMultiSigHandler は、ログインユーザー（Owner）と2名の参加者からマルチシグを作成しDBに登録します。
func CreateMultiSigHandler(c *gin.Context) {
	var req struct {
		Owner        string   `json:"owner"`        // ログイン済みのユーザーアドレス
		Participants []string `json:"participants"` // 参加者のEthereumアドレス（2名）
		Address      string   `json:"address"`        // マルチシグ公開鍵のアドレス
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Participants) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Require owner and exactly 2 participants"})
		return
	}

	// マルチシグIDを生成し、初期状態を設定
	newMultiSig := models.MultiSig{
		Address:      req.Address,
		Owner:        req.Owner,
		Participants: datatypes.JSON([]byte(mustMarshal(req.Participants))),
		Status:       "awaiting",
		Data:         datatypes.JSON([]byte(`{}`)),
	}

	if err := db.DB.Create(&newMultiSig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error on create"})
		return
	}

	for _, address := range req.Participants {
    	var user models.User
    	if err := db.DB.First(&user, "address = ?", address).Error; err != nil {
        	if errors.Is(err, gorm.ErrRecordNotFound) {
            	continue
        	} else {
            	c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error fetching user"})
				return
        	}
    	}

    	var msAddresses []string
    	if err := json.Unmarshal(user.MultiSigs, &msAddresses); err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user's multisigs list"})
        	return
    	}
		msAddresses = append(msAddresses, req.Address)
		
		user = models.User{Address: address, MultiSigs: datatypes.JSON([]byte(mustMarshal(msAddresses)))}
		if err := db.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "MultiSig created", "multisig": newMultiSig})
}

// GetMultiSigListHandler は、指定ユーザーが参加しているマルチシグの一覧を返します。
func GetMultiSigListHandler(c *gin.Context) {
    userAddress := c.Query("address")
    if userAddress == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "address parameter is required"})
        return
    }

    var user models.User
    if err := db.DB.First(&user, "address = ?", userAddress).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error fetching user"})
        }
        return
    }

    var msAddresses []string
    if err := json.Unmarshal(user.MultiSigs, &msAddresses); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user's multisigs list"})
        return
    }

    if len(msAddresses) == 0 {
        c.JSON(http.StatusOK, []models.MultiSig{})
        return
    }

    var multisigs []models.MultiSig
    if err := db.DB.
        Where("address IN ?", msAddresses).
        Find(&multisigs).
        Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error on list"})
        return
    }

    // 6. 取得結果を返却
    c.JSON(http.StatusOK, multisigs)
}

// GetMultiSigDataHandler は、指定マルチシグの署名用データ（例としてプレースホルダー）を生成し返します。
func GetMultiSigDataHandler(c *gin.Context) {
	address := c.Param("address")
	var ms models.MultiSig
	if err := db.DB.First(&ms, "address = ?", address).Error; err != nil {
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
	address := c.Param("address")
	var payload struct {
		Signature string `json:"signature"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signature payload"})
		return
	}

	var ms models.MultiSig
	if err := db.DB.First(&ms, "address = ?", address).Error; err != nil {
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
