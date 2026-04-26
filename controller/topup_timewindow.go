package controller

import (
	"fmt"
	"net/http"

	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/gin-gonic/gin"
)

// checkTopupTimeWindow 检查当前是否在充值时间窗口内
// 不在窗口内时返回 false 并写入 JSON 错误响应
func checkTopupTimeWindow(c *gin.Context) bool {
	inWindow, windowStr := operation_setting.IsInTopupTimeWindow()
	if inWindow {
		return true
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "error",
		"data":    fmt.Sprintf("当前不在充值时间窗口内（窗口时间：%s）", windowStr),
	})
	return false
}
