package main

import (
	"bytes"
	"fmt"
	"log"
	"text/tabwriter"

	"github.com/MarcosS99461/aliyun-vpc-ipam-checker/internal/config"
	"github.com/MarcosS99461/aliyun-vpc-ipam-checker/pkg/alicloud"
	"github.com/MarcosS99461/aliyun-vpc-ipam-checker/pkg/webhook"
)

func main() {
	// 解析命令行参数
	flags := config.ParseFlags()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建配置审计客户端
	configClient, err := alicloud.NewConfigClient(cfg.AccessKeyID, cfg.AccessKeySecret, cfg.RegionID, cfg.AggregatorID)
	if err != nil {
		log.Fatalf("创建配置审计客户端失败: %v", err)
	}

	// 获取VPC列表
	vpcs, err := configClient.ListVPCs()
	if err != nil {
		log.Fatalf("获取VPC列表失败: %v", err)
	}

	// 如果提供了IPAM凭证，则获取IPAM管理的VPC信息
	var vpcToPoolName map[string]string
	if cfg.IPAMAccessKeyID != "" && cfg.IPAMAccessKeySecret != "" {
		ipamClient, err := alicloud.NewIPAMClient(cfg.IPAMAccessKeyID, cfg.IPAMAccessKeySecret, cfg.RegionID)
		if err != nil {
			log.Fatalf("创建IPAM客户端失败: %v", err)
		}

		vpcToPoolName, err = ipamClient.ListManagedVPCs()
		if err != nil {
			log.Fatalf("获取IPAM管理的VPC列表失败: %v", err)
		}
	}

	// 统计信息
	totalVPCs := 0
	activeVPCs := 0
	ipamManagedVPCs := len(vpcToPoolName)
	filteredVPCs := 0

	// 创建输出缓冲区
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	// 写入表头
	fmt.Fprintln(w, "账号ID\tVPC ID\tVPC名称\t地域\t状态\tIPAM池")
	fmt.Fprintln(w, "--------\t--------\t--------\t--------\t--------\t--------")

	// 处理每个VPC
	for _, vpc := range vpcs {
		if vpc.ResourceType == "ACS::VPC::VPC" {
			totalVPCs++
			if vpc.ResourceStatus == "Available" {
				activeVPCs++
			}

			// 检查是否需要排除此账号
			shouldExclude := false
			for _, excludeAccount := range flags.ExcludeAccounts {
				if vpc.AccountId == excludeAccount {
					shouldExclude = true
					break
				}
			}
			if shouldExclude {
				continue
			}

			// 获取IPAM池信息
			ipamPool := vpcToPoolName[vpc.ResourceId]
			isManaged := ipamPool != ""
			if ipamPool == "" {
				ipamPool = "未托管"
			}

			// 如果只显示未托管的VPC，则跳过已托管的
			if flags.OnlyUnmanaged && isManaged {
				continue
			}

			// 写入VPC信息
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
				vpc.AccountId,
				vpc.ResourceId,
				vpc.ResourceName,
				vpc.Region,
				vpc.ResourceStatus,
				ipamPool,
			)
			filteredVPCs++
		}
	}

	// 刷新输出缓冲区
	w.Flush()

	// 打印统计信息和结果
	fmt.Printf("找到 %d 个VPC资源 (活跃: %d, IPAM管理: %d, 符合筛选条件: %d)\n\n", 
		totalVPCs, activeVPCs, ipamManagedVPCs, filteredVPCs)
	fmt.Print(buf.String())

	// 如果提供了webhook URL，发送结果
	if flags.WebhookURL != "" {
		stats := struct {
			Total    int
			Active   int
			Managed  int
			Filtered int
		}{
			Total:    totalVPCs,
			Active:   activeVPCs,
			Managed:  ipamManagedVPCs,
			Filtered: filteredVPCs,
		}

		if err := webhook.SendMessage(flags.WebhookURL, stats, buf.String()); err != nil {
			log.Printf("发送webhook消息失败: %v", err)
		} else {
			fmt.Println("\n已将结果发送到webhook")
		}
	}
}