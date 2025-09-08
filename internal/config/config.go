package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config 结构体定义了程序运行所需的配置信息
type Config struct {
	// 配置审计相关配置（必需）
	AccessKeyID     string // 访问密钥ID
	AccessKeySecret string // 访问密钥密码
	RegionID        string // 地域ID
	AggregatorID    string // 聚合器ID

	// IPAM相关配置（可选）
	IPAMAccessKeyID     string // IPAM系统的访问密钥ID
	IPAMAccessKeySecret string // IPAM系统的访问密钥密码
}

// LoadConfig 从环境变量加载配置信息
// 支持从.env文件加载配置
// 配置审计的凭证是必需的，IPAM的凭证是可选的
func LoadConfig() (*Config, error) {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("加载.env文件失败: %v", err)
	}

	// 检查必需的环境变量
	accessKeyID := os.Getenv("ACCESS_KEY_ID")
	if accessKeyID == "" {
		return nil, fmt.Errorf("环境变量ACCESS_KEY_ID未设置")
	}

	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	if accessKeySecret == "" {
		return nil, fmt.Errorf("环境变量ACCESS_KEY_SECRET未设置")
	}

	regionID := os.Getenv("REGION_ID")
	if regionID == "" {
		return nil, fmt.Errorf("环境变量REGION_ID未设置")
	}

	aggregatorID := os.Getenv("AGGREGATOR_ID")
	if aggregatorID == "" {
		return nil, fmt.Errorf("环境变量AGGREGATOR_ID未设置")
	}

	// IPAM相关配置（可选）
	ipamAccessKeyID := os.Getenv("ALIBABA_CLOUD_IPAM_ACCESS_KEY_ID")
	ipamAccessKeySecret := os.Getenv("ALIBABA_CLOUD_IPAM_ACCESS_KEY_SECRET")

	return &Config{
		AccessKeyID:         accessKeyID,
		AccessKeySecret:     accessKeySecret,
		RegionID:            regionID,
		AggregatorID:        aggregatorID,
		IPAMAccessKeyID:     ipamAccessKeyID,
		IPAMAccessKeySecret: ipamAccessKeySecret,
	}, nil
}
