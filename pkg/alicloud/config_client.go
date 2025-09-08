package alicloud

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ConfigClient struct {
	accessKeyID     string
	accessKeySecret string
	regionID        string
	aggregatorID    string
}

type VPCResource struct {
	AccountId      int64  `json:"AccountId"`
	ResourceId     string `json:"ResourceId"`
	ResourceName   string `json:"ResourceName"`
	Region         string `json:"Region"`
	ResourceStatus string `json:"ResourceStatus"`
	ResourceType   string `json:"ResourceType"`
}

func NewConfigClient(accessKeyID, accessKeySecret, regionID, aggregatorID string) (*ConfigClient, error) {
	return &ConfigClient{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		regionID:        regionID,
		aggregatorID:    aggregatorID,
	}, nil
}

func (c *ConfigClient) ListVPCs() ([]VPCResource, error) {
	var vpcs []VPCResource
	nextToken := ""

	for {
		// 准备请求参数
		params := map[string]string{
			"Action":           "ListAggregateDiscoveredResources",
			"Format":           "JSON",
			"Version":          "2020-09-07",
			"SignatureMethod":  "HMAC-SHA1",
			"SignatureVersion": "1.0",
			"RegionId":         c.regionID,
			"AggregatorId":     c.aggregatorID,
			"ResourceType":     "ACS::VPC::VPC",
			"MaxResults":       "100",
		}

		if nextToken != "" {
			params["NextToken"] = nextToken
		}

		// 发送请求
		resp, err := c.sendRequest(params)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		// 解析响应
		var result struct {
			DiscoveredResourceProfiles struct {
				DiscoveredResourceProfileList []VPCResource `json:"DiscoveredResourceProfileList"`
				NextToken                     *string       `json:"NextToken"`
			} `json:"DiscoveredResourceProfiles"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		resp.Body.Close()

		// 添加VPC到结果列表
		vpcs = append(vpcs, result.DiscoveredResourceProfiles.DiscoveredResourceProfileList...)

		// 检查是否还有更多数据
		if result.DiscoveredResourceProfiles.NextToken == nil || *result.DiscoveredResourceProfiles.NextToken == "" {
			break
		}
		nextToken = *result.DiscoveredResourceProfiles.NextToken
	}

	return vpcs, nil
}

func (c *ConfigClient) sendRequest(params map[string]string) (*http.Response, error) {
	// 添加公共参数
	params["AccessKeyId"] = c.accessKeyID
	params["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	params["SignatureNonce"] = fmt.Sprintf("%d", time.Now().UnixNano())

	// 生成签名
	signature := c.generateSignature(params)
	params["Signature"] = signature

	// 构建请求URL
	endpoint := "https://config." + c.regionID + ".aliyuncs.com"
	reqURL := endpoint + "/?" + c.buildQuery(params)

	// 发送请求
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

func (c *ConfigClient) generateSignature(params map[string]string) string {
	// 1. 按参数名称排序
	var keys []string
	for k := range params {
		if k != "Signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 2. 构建规范化请求字符串
	var canonicalizedQueryString bytes.Buffer
	for i, k := range keys {
		if i > 0 {
			canonicalizedQueryString.WriteByte('&')
		}
		canonicalizedQueryString.WriteString(url.QueryEscape(k))
		canonicalizedQueryString.WriteByte('=')
		canonicalizedQueryString.WriteString(url.QueryEscape(params[k]))
	}

	// 3. 构建签名字符串
	stringToSign := "GET&%2F&" + url.QueryEscape(canonicalizedQueryString.String())

	// 4. 计算签名
	key := []byte(c.accessKeySecret + "&")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signature
}

func (c *ConfigClient) buildQuery(params map[string]string) string {
	var pairs []string
	for k, v := range params {
		pairs = append(pairs, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	return strings.Join(pairs, "&")
}
