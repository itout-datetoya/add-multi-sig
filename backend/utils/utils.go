package utils

import "fmt"

// LogInfo はログ出力用のユーティリティ関数です。
func LogInfo(msg string) {
	fmt.Println("[INFO]", msg)
}
