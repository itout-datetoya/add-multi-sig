package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User はEthereumアドレスで識別されるユーザー情報です。
type User struct {
	gorm.Model
	Address string `gorm:"uniqueIndex;not null" json:"address"` // ユーザーアドレス
	Pubkey string  `gorm:"not null" json:"pubkey"` // paillier公開鍵
	MultiSigs datatypes.JSON `gorm:"type:jsonb" json:"multisigs"` // 参加マルチシグアドレスリスト
}
