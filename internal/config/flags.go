package config

import (
	"flag"
	"fmt"
	"strings"
)

// Flags 结构体定义了程序支持的命令行参数
type Flags struct {
	// ExcludeAccounts 存储要排除的账号ID列表
	ExcludeAccounts []int64

	// OnlyUnmanaged 表示是否只显示未被IPAM管理的VPC
	OnlyUnmanaged bool

	// WebhookURL 存储Lark机器人的webhook地址
	WebhookURL string
}

// ParseFlags 解析命令行参数并返回Flags结构体
// 支持的参数：
// -exclude: 要排除的账号ID列表，用逗号分隔
// -unmanaged: 只显示未被IPAM管理的VPC
// -webhook: Lark机器人的webhook地址
func ParseFlags() *Flags {
	// 定义命令行参数
	excludeAccounts := flag.String("exclude", "", "要排除的账号ID列表，用逗号分隔，例如: 5130150745510468,5020439416629852")
	onlyUnmanaged := flag.Bool("unmanaged", false, "只显示未被IPAM管理的VPC")
	webhookURL := flag.String("webhook", "", "Webhook URL，用于发送结果")

	// 解析命令行参数
	flag.Parse()

	// 处理排除账号列表
	var excludeAccountsList []int64
	if *excludeAccounts != "" {
		accounts := strings.Split(*excludeAccounts, ",")
		for _, account := range accounts {
			var accountID int64
			if _, err := fmt.Sscanf(account, "%d", &accountID); err == nil {
				excludeAccountsList = append(excludeAccountsList, accountID)
			}
		}
	}

	return &Flags{
		ExcludeAccounts: excludeAccountsList,
		OnlyUnmanaged:   *onlyUnmanaged,
		WebhookURL:      *webhookURL,
	}
}
