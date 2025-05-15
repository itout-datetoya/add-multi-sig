package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MultiSig はマルチシグに関する情報を保持します。
type MultiSig struct {
	gorm.Model
	Address      string         `gorm:"uniqueIndex;not null" json:"address"` // マルチシグ公開鍵
	Owner        string         `gorm:"not null" json:"owner"` // 作成者アドレス
	Participants datatypes.JSON `gorm:"type:jsonb" json:"participants"` // 参加者アドレスのJSON配列
	Status       string         `gorm:"not null" json:"status"`           // "awaiting", "partial", "completed"
	Data         datatypes.JSON `gorm:"type:jsonb" json:"data"`           // 署名に必要な中間データ
}
