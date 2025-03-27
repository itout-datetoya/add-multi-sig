package models

import "gorm.io/gorm"

// User はEthereumアドレスで識別されるユーザー情報です。
type User struct {
	gorm.Model
	Address string `gorm:"uniqueIndex;not null" json:"address"`
	// その他の情報（例：ニックネーム、公開鍵など）を追加可能
}
