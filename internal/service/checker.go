package service

import (
	"fmt"
	"strings"
	"vpc-ipam-checker/internal/config"
	"vpc-ipam-checker/pkg/alicloud"
	"vpc-ipam-checker/pkg/lark"
)

// Checker VPC检查服务
type Checker struct {
	cfg        *config.Config
	configVPC  *alicloud.Client
	ipamVPC    *alicloud.Client
	larkClient *lark.Client
}

// NewChecker 创建新的检查服务
func NewChecker(cfg *config.Config) (*Checker, error) {
	// 创建配置中心VPC客户端（使用STS）
	configVPC, err := alicloud.NewClientWithSTS(
		cfg.ConfigCenterAccessKeyID,
		cfg.ConfigCenterAccessKeySecret,
		cfg.ConfigCenterSTSRoleARN,
		cfg.ConfigCenterRoleSessionName,
		cfg.RegionID,
		cfg.STSDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("创建配置中心客户端失败: %v", err)
	}

	// 创建IPAM VPC客户端（使用STS）
	ipamVPC, err := alicloud.NewClientWithSTS(
		cfg.IPAMAccessKeyID,
		cfg.IPAMAccessKeySecret,
		cfg.IPAMSTSRoleARN,
		cfg.IPAMRoleSessionName,
		cfg.RegionID,
		cfg.STSDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("创建IPAM客户端失败: %v", err)
	}

	// 创建Lark客户端
	larkClient := lark.NewClient(cfg.LarkWebhookURL)

	return &Checker{
		cfg:        cfg,
		configVPC:  configVPC,
		ipamVPC:    ipamVPC,
		larkClient: larkClient,
	}, nil
}

// Run 运行检查
func (c *Checker) Run() error {
	// 获取配置中心的VPC列表
	configVPCs, err := c.configVPC.GetVPCs()
	if err != nil {
		return fmt.Errorf("获取配置中心VPC列表失败: %v", err)
	}

	// 获取IPAM的VPC列表
	ipamVPCs, err := c.ipamVPC.GetVPCs()
	if err != nil {
		return fmt.Errorf("获取IPAM VPC列表失败: %v", err)
	}

	// 创建IPAM VPC ID映射
	ipamVPCMap := make(map[string]bool)
	for _, vpc := range ipamVPCs {
		ipamVPCMap[vpc.VpcID] = true
	}

	// 找出未被IPAM管理的VPC
	var unmanagedVPCs []alicloud.VPC
	for _, vpc := range configVPCs {
		if !ipamVPCMap[vpc.VpcID] {
			unmanagedVPCs = append(unmanagedVPCs, vpc)
		}
	}

	// 如果有未管理的VPC，发送通知
	if len(unmanagedVPCs) > 0 {
		return c.sendNotification(unmanagedVPCs)
	}

	return nil
}

// sendNotification 发送通知
func (c *Checker) sendNotification(vpcs []alicloud.VPC) error {
	var builder strings.Builder
	builder.WriteString("以下VPC未被IPAM管理：\n\n")

	for _, vpc := range vpcs {
		builder.WriteString(fmt.Sprintf("VPC ID: %s\n", vpc.VpcID))
		builder.WriteString(fmt.Sprintf("User ID: %s\n", vpc.UserID))
		builder.WriteString(fmt.Sprintf("VPC Name: %s\n", vpc.VpcName))
		builder.WriteString(fmt.Sprintf("Region: %s\n", vpc.RegionID))
		builder.WriteString(fmt.Sprintf("Created At: %s\n", vpc.CreatedAt.Format("2006-01-02 15:04:05")))
		builder.WriteString("-------------------\n")
	}

	return c.larkClient.SendText(builder.String())
}
