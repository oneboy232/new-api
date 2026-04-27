package service

import (
	"fmt"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
)

// NotifyTopupByTradeNo 根据订单号查询充值记录，并推送充值成功通知给用户。
// 用户需在通知设置中配置通知方式（webhook/bark/gotify/email）后才能收到。
func NotifyTopupByTradeNo(tradeNo string) {
	if tradeNo == "" {
		return
	}
	topUp := model.GetTopUpByTradeNo(tradeNo)
	if topUp == nil {
		return
	}
	if topUp.Status != common.TopUpStatusSuccess {
		return
	}

	user, err := model.GetUserById(topUp.UserId, false)
	if err != nil || user == nil {
		return
	}

	userSetting := user.GetSetting()
	if userSetting.NotifyType == "" {
		return // 用户未配置通知方式，跳过
	}

	// 格式化充值信息
	providerName := topUp.PaymentProvider
	if providerName == "" {
		providerName = topUp.PaymentMethod
	}

	title := "充值成功通知"
	content := fmt.Sprintf("充值成功！\n用户：%s\n金额：$%.2f\n到账额度：%d\n支付方式：%s",
		user.Username,
		topUp.Money,
		topUp.Amount,
		providerName)

	// 发送飞书群通知
	feishuTitle := "充值成功通知"
	feishuContent := fmt.Sprintf("用户：%s\n金额：$%.2f\n到账额度：%d\n支付方式：%s\n订单号：%s",
		user.Username,
		topUp.Money,
		topUp.Amount,
		providerName,
		topUp.TradeNo)
	SendFeishuNotify(feishuTitle, feishuContent)

	// 发送用户通知（需用户个人配置通知方式）
	notification := dto.NewNotify(dto.NotifyTypeTopup, title, content, nil)
	if err := NotifyUser(topUp.UserId, user.Email, userSetting, notification); err != nil {
		common.SysLog(fmt.Sprintf("failed to send topup notification to user %d: %s", topUp.UserId, err.Error()))
	}
}
