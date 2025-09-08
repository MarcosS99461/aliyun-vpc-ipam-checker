package alicloud

import (
	"fmt"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// RAMClient RAM客户端
type RAMClient struct {
	client *openapi.Client
	cache  map[int64]string
	mutex  sync.RWMutex
}

// NewRAMClient 创建RAM客户端
func NewRAMClient(accessKeyID, accessKeySecret string) (*RAMClient, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("ram.aliyuncs.com"),
	}

	client, err := openapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建RAM客户端失败: %v", err)
	}

	return &RAMClient{
		client: client,
		cache:  make(map[int64]string),
	}, nil
}

// GetAccountAlias 获取账号别名
func (c *RAMClient) GetAccountAlias(accountID int64) (string, error) {
	// 先检查缓存
	c.mutex.RLock()
	if alias, ok := c.cache[accountID]; ok {
		c.mutex.RUnlock()
		return alias, nil
	}
	c.mutex.RUnlock()

	// 由于RAM API的限制，这里暂时只返回账号ID
	alias := fmt.Sprintf("%d", accountID)

	// 保存到缓存
	c.mutex.Lock()
	c.cache[accountID] = alias
	c.mutex.Unlock()

	return alias, nil
}
