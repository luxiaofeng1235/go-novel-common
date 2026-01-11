package db

import (
	lock "github.com/sjy3/go-keylock"
	"go-novel/global"
)

func InitKeyLock() {
	var keyLock = lock.NewKeyLock()
	global.KeyLock = keyLock
}
