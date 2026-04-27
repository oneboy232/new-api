package operation_setting

import (
	"fmt"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/setting/config"
)

// 额度展示类型
const (
	QuotaDisplayTypeUSD    = "USD"
	QuotaDisplayTypeCNY    = "CNY"
	QuotaDisplayTypeTokens = "TOKENS"
	QuotaDisplayTypeCustom = "CUSTOM"
)

type GeneralSetting struct {
	DocsLink            string `json:"docs_link"`
	PingIntervalEnabled bool   `json:"ping_interval_enabled"`
	PingIntervalSeconds int    `json:"ping_interval_seconds"`
	// 当前站点额度展示类型：USD / CNY / TOKENS
	QuotaDisplayType string `json:"quota_display_type"`
	// 自定义货币符号，用于 CUSTOM 展示类型
	CustomCurrencySymbol string `json:"custom_currency_symbol"`
	// 自定义货币与美元汇率（1 USD = X Custom）
	CustomCurrencyExchangeRate float64 `json:"custom_currency_exchange_rate"`
	// 飞书 Webhook 地址，支持多行（每行一个 URL），充值成功时通知到此地址
	FeishuWebhookUrl string `json:"feishu_webhook_url"`
	// 充值时间窗口
	TopupTimeWindowEnabled bool   `json:"topup_time_window_enabled"`
	TopupTimeWindowStart   string `json:"topup_time_window_start"`
	TopupTimeWindowEnd     string `json:"topup_time_window_end"`
	// 兑换码充值
	RedemptionEnabled bool `json:"redemption_enabled"`
}

// 默认配置
var generalSetting = GeneralSetting{
	DocsLink:                   "https://docs.newapi.pro",
	PingIntervalEnabled:        false,
	PingIntervalSeconds:        60,
	QuotaDisplayType:           QuotaDisplayTypeUSD,
	CustomCurrencySymbol:       "¤",
	CustomCurrencyExchangeRate: 1.0,
	TopupTimeWindowEnabled:     false,
	TopupTimeWindowStart:       "08:00",
	TopupTimeWindowEnd:         "20:00",
		RedemptionEnabled:          true,
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("general_setting", &generalSetting)
}

func GetGeneralSetting() *GeneralSetting {
	return &generalSetting
}

// IsCurrencyDisplay 是否以货币形式展示（美元或人民币）
func IsCurrencyDisplay() bool {
	return generalSetting.QuotaDisplayType != QuotaDisplayTypeTokens
}

// IsCNYDisplay 是否以人民币展示
func IsCNYDisplay() bool {
	return generalSetting.QuotaDisplayType == QuotaDisplayTypeCNY
}

// GetQuotaDisplayType 返回额度展示类型
func GetQuotaDisplayType() string {
	return generalSetting.QuotaDisplayType
}

// GetCurrencySymbol 返回当前展示类型对应符号
func GetCurrencySymbol() string {
	switch generalSetting.QuotaDisplayType {
	case QuotaDisplayTypeUSD:
		return "$"
	case QuotaDisplayTypeCNY:
		return "¥"
	case QuotaDisplayTypeCustom:
		if generalSetting.CustomCurrencySymbol != "" {
			return generalSetting.CustomCurrencySymbol
		}
		return "¤"
	default:
		return ""
	}
}

// GetFeishuWebhookUrls 返回飞书 Webhook URL 列表，支持多行配置
func GetFeishuWebhookUrls() []string {
	urls := strings.Split(generalSetting.FeishuWebhookUrl, "\n")
	var result []string
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u != "" {
			result = append(result, u)
		}
	}
	return result
}

// IsInTopupTimeWindow 判断当前时间是否在充值时间窗口内
// 返回 (true, "") 表示在窗口内或未启用
// 返回 (false, "08:00 - 20:00") 表示不在窗口内，附带窗口时间
func IsInTopupTimeWindow() (bool, string) {
	gs := GetGeneralSetting()
	if !gs.TopupTimeWindowEnabled {
		return true, ""
	}

	start, err := time.Parse("15:04", gs.TopupTimeWindowStart)
	if err != nil {
		return true, ""
	}
	end, err := time.Parse("15:04", gs.TopupTimeWindowEnd)
	if err != nil {
		return true, ""
	}

	now := time.Now()
	nowMinutes := now.Hour()*60 + now.Minute()
	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()

	inWindow := false
	if startMinutes <= endMinutes {
		inWindow = nowMinutes >= startMinutes && nowMinutes < endMinutes
	} else {
		inWindow = nowMinutes >= startMinutes || nowMinutes < endMinutes
	}

	if inWindow {
		return true, ""
	}
	return false, fmt.Sprintf("%s - %s", gs.TopupTimeWindowStart, gs.TopupTimeWindowEnd)
}

// GetUsdToCurrencyRate 返回 1 USD = X <currency> 的 X（TOKENS 不适用）
func GetUsdToCurrencyRate(usdToCny float64) float64 {
	switch generalSetting.QuotaDisplayType {
	case QuotaDisplayTypeUSD:
		return 1
	case QuotaDisplayTypeCNY:
		return usdToCny
	case QuotaDisplayTypeCustom:
		if generalSetting.CustomCurrencyExchangeRate > 0 {
			return generalSetting.CustomCurrencyExchangeRate
		}
		return 1
	default:
		return 1
	}
}
