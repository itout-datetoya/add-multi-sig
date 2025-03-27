package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MultiSig はマルチシグに関する情報を保持します。
type MultiSig struct {
	gorm.Model
	Owner        string         `gorm:"not null" json:"owner"` // 作成者（ログインユーザー）
	Participants datatypes.JSON `gorm:"type:jsonb" json:"participants"` // EthereumアドレスのJSON配列
	Status       string         `gorm:"not null" json:"status"`           // 例: "awaiting", "partial", "completed"
	Data         datatypes.JSON `gorm:"type:jsonb" json:"data"`           // 署名に必要な中間データ
}
