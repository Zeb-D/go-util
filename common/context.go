package common

import (
	"context"
	"time"
)

//锁相关

func NewContextWithTimeout(n time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), n*time.Second)
}

//NewContext 创建一个Context
func NewContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
