package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func SendMessage(webhookURL string, stats struct {
	Total    int
	Active   int
	Managed  int
	Filtered int
}, details string) error {
	if webhookURL == "" {
		return nil
	}

	// 直接使用原始格式
	text := fmt.Sprintf("⚠️ VPC资源统计\nVPC总数: %d\n活跃VPC: %d\nIPAM管理: %d\n符合筛选条件: %d\n\n%s",
		stats.Total, stats.Active, stats.Managed, stats.Filtered, details)

	msg := Message{
		MsgType: "text",
	}
	msg.Content.Text = text

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send webhook message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
