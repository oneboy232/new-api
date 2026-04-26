package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

func init() {
	// 注册配额增加回调，发送系统级飞书通知
	model.OnQuotaIncreased = func(userId int, quota int) {
		user, err := model.GetUserById(userId, false)
		if err != nil || user == nil {
			common.SysLog(fmt.Sprintf("failed to get user %d for quota increased notification: %s", userId, err.Error()))
			return
		}

		amount := float64(quota) / common.QuotaPerUnit
		content := fmt.Sprintf("用户：%s\n金额：$%.2f\n增加额度：%d",
			user.Username,
			amount,
			quota)

		SendFeishuNotify("配额增加通知", content)
	}
}

// feishuCardMessage 飞书交互式卡片消息结构
type feishuCardMessage struct {
	MsgType string         `json:"msg_type"`
	Card    feishuCardBody `json:"card"`
}

type feishuCardBody struct {
	Header   feishuCardHeader   `json:"header"`
	Elements []feishuCardElement `json:"elements"`
}

type feishuCardHeader struct {
	Title feishuText `json:"title"`
}

type feishuCardElement struct {
	Tag  string     `json:"tag"`
	Text feishuText `json:"text"`
}

type feishuText struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// SendFeishuNotify 发送飞书通知到所有已配置的 Webhook URL
func SendFeishuNotify(title string, content string) {
	urls := getFeishuWebhookUrls()
	if len(urls) == 0 {
		return
	}

	// 构建飞书消息
	lines := strings.Split(content, "\n")
	var elementText string
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "：", 2)
		if len(parts) == 2 {
			elementText += fmt.Sprintf("**%s**：%s\n", parts[0], parts[1])
		} else {
			elementText += line + "\n"
		}
	}

	msg := feishuCardMessage{
		MsgType: "interactive",
		Card: feishuCardBody{
			Header: feishuCardHeader{
				Title: feishuText{
					Tag:     "plain_text",
					Content: title,
				},
			},
			Elements: []feishuCardElement{
				{
					Tag: "div",
					Text: feishuText{
						Tag:     "lark_md",
						Content: strings.TrimSpace(elementText),
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(msg)
	if err != nil {
		common.SysLog(fmt.Sprintf("failed to marshal feishu message: %s", err.Error()))
		return
	}

	for _, webhookUrl := range urls {
		sendFeishuSingle(webhookUrl, payloadBytes)
	}
}

func sendFeishuSingle(webhookUrl string, payloadBytes []byte) {
	var req *http.Request
	var resp *http.Response
	var err error

	if system_setting.EnableWorker() {
		workerReq := &WorkerRequest{
			URL:    webhookUrl,
			Key:    system_setting.WorkerValidKey,
			Method: http.MethodPost,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			Body: payloadBytes,
		}

		resp, err = DoWorkerRequest(workerReq)
		if err != nil {
			common.SysLog(fmt.Sprintf("failed to send feishu notify through worker: %s", err.Error()))
			return
		}
		defer resp.Body.Close()
	} else {
		fetchSetting := system_setting.GetFetchSetting()
		if err := common.ValidateURLWithFetchSetting(webhookUrl, fetchSetting.EnableSSRFProtection, fetchSetting.AllowPrivateIp, fetchSetting.DomainFilterMode, fetchSetting.IpFilterMode, fetchSetting.DomainList, fetchSetting.IpList, fetchSetting.AllowedPorts, fetchSetting.ApplyIPFilterForDomain); err != nil {
			common.SysLog(fmt.Sprintf("feishu webhook URL validation failed: %s", err.Error()))
			return
		}

		req, err = http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(payloadBytes))
		if err != nil {
			common.SysLog(fmt.Sprintf("failed to create feishu request: %s", err.Error()))
			return
		}

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("User-Agent", "NewAPI-Feishu-Notify/1.0")

		client := GetHttpClient()
		resp, err = client.Do(req)
		if err != nil {
			common.SysLog(fmt.Sprintf("failed to send feishu notify: %s", err.Error()))
			return
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		common.SysLog(fmt.Sprintf("feishu webhook request failed with status code: %d", resp.StatusCode))
	}
}

// getFeishuWebhookUrls 从设置中获取飞书 Webhook URL 列表
func getFeishuWebhookUrls() []string {
	return operation_setting.GetFeishuWebhookUrls()
}
