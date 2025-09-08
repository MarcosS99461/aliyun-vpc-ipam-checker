package alicloud

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	vpcipam "github.com/alibabacloud-go/vpcipam-20230228/client"
)

type IPAMClient struct {
	client *vpcipam.Client
}

func NewIPAMClient(accessKeyID, accessKeySecret, regionID string) (*IPAMClient, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(regionID),
	}
	config.Endpoint = tea.String("vpcipam." + regionID + ".aliyuncs.com")
	client, err := vpcipam.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPAM client: %w", err)
	}
	return &IPAMClient{client: client}, nil
}

func (c *IPAMClient) ListManagedVPCs() (map[string]string, error) {
	// 存储VPC ID到IPAM池名称的映射
	vpcToPoolName := make(map[string]string)

	// 获取所有IPAM池
	nextToken := ""
	for {
		request := &vpcipam.ListIpamPoolsRequest{
			RegionId:   tea.String("ap-southeast-1"),
			MaxResults: tea.Int32(100),
		}
		if nextToken != "" {
			request.NextToken = tea.String(nextToken)
		}

		response, err := c.client.ListIpamPools(request)
		if err != nil {
			return nil, fmt.Errorf("failed to list IPAM pools: %w", err)
		}

		// 对每个IPAM池获取其分配信息
		for _, pool := range response.Body.IpamPools {
			allocRequest := &vpcipam.ListIpamPoolAllocationsRequest{
				RegionId:   tea.String("ap-southeast-1"),
				IpamPoolId: pool.IpamPoolId,
				MaxResults: tea.Int32(100),
			}

			// 获取该池的所有分配
			allocNextToken := ""
			for {
				if allocNextToken != "" {
					allocRequest.NextToken = tea.String(allocNextToken)
				}

				allocResponse, err := c.client.ListIpamPoolAllocations(allocRequest)
				if err != nil {
					return nil, fmt.Errorf("failed to list allocations for pool %s: %w", *pool.IpamPoolId, err)
				}

				// 记录VPC到池名称的映射
				for _, alloc := range allocResponse.Body.IpamPoolAllocations {
					if *alloc.ResourceType == "VPC" {
						vpcToPoolName[*alloc.ResourceId] = *pool.IpamPoolName
					}
				}

				// 检查是否还有更多分配
				if allocResponse.Body.NextToken == nil || *allocResponse.Body.NextToken == "" {
					break
				}
				allocNextToken = *allocResponse.Body.NextToken
			}
		}

		// 检查是否还有更多池
		if response.Body.NextToken == nil || *response.Body.NextToken == "" {
			break
		}
		nextToken = *response.Body.NextToken
	}

	return vpcToPoolName, nil
}
