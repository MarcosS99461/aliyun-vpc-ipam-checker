# aliyun-vpc-ipam-checker

[![Go Report Card](https://goreportcard.com/badge/github.com/MarcosS99461/aliyun-vpc-ipam-checker)](https://goreportcard.com/report/github.com/MarcosS99461/aliyun-vpc-ipam-checker)
[![License](https://img.shields.io/github/license/MarcosS99461/aliyun-vpc-ipam-checker)](https://github.com/MarcosS99461/aliyun-vpc-ipam-checker/blob/main/LICENSE)

这个工具用于检查阿里云账号中未被 IPAM（IP 地址管理）系统管理的 VPC 资源。它通过比对配置审计中心的 VPC 资源列表和 IPAM 系统中的 VPC 资源，找出那些尚未纳入 IPAM 管理的 VPC。

## 下载和安装

### 预编译二进制文件

我们为以下平台提供了预编译的二进制文件：

- macOS
  - [Apple Silicon (M1/M2)](https://github.com/MarcosS99461/aliyun-vpc-ipam-checker/releases/download/v1.0.0/vpc-checker-v1.0.0-darwin-arm64.tar.gz)
  - [Intel](https://github.com/MarcosS99461/aliyun-vpc-ipam-checker/releases/download/v1.0.0/vpc-checker-v1.0.0-darwin-amd64.tar.gz)
- Linux
  - [x86_64/amd64](https://github.com/MarcosS99461/aliyun-vpc-ipam-checker/releases/download/v1.0.0/vpc-checker-v1.0.0-linux-amd64.tar.gz)
  - [arm64](https://github.com/MarcosS99461/aliyun-vpc-ipam-checker/releases/download/v1.0.0/vpc-checker-v1.0.0-linux-arm64.tar.gz)
- Windows
  - [64-bit](https://github.com/MarcosS99461/aliyun-vpc-ipam-checker/releases/download/v1.0.0/vpc-checker-v1.0.0-windows-amd64.exe.tar.gz)

下载后解压并运行：
```bash
# macOS/Linux
tar xzf vpc-checker-*.tar.gz
chmod +x vpc-checker-*
./vpc-checker-* -h

# Windows
# 解压 .tar.gz 文件后直接运行 .exe 文件
```

### 从源码构建

需要 Go 1.20 或更高版本：

```bash
git clone https://github.com/MarcosS99461/aliyun-vpc-ipam-checker.git
cd aliyun-vpc-ipam-checker
go build -v ./cmd/vpc-checker
```

## 功能特性

- 通过配置审计中心获取所有 VPC 资源
- 通过 IPAM 系统获取已托管的 VPC 资源
- 支持按账号 ID 排除特定账号
- 支持仅显示未托管的 VPC 资源
- 支持将结果发送到 Lark 机器人

## 配置

在项目根目录创建 `.env` 文件，包含以下配置：

```env
# 配置审计账号凭证（必需）
ACCESS_KEY_ID=your_access_key_id
ACCESS_KEY_SECRET=your_access_key_secret
REGION_ID=ap-southeast-1
AGGREGATOR_ID=your_aggregator_id

# IPAM账号凭证（可选）
ALIBABA_CLOUD_IPAM_ACCESS_KEY_ID=your_ipam_access_key_id
ALIBABA_CLOUD_IPAM_ACCESS_KEY_SECRET=your_ipam_access_key_secret
```

## 使用方法

基本用法：
```bash
./vpc-checker
```

排除特定账号：
```bash
./vpc-checker -exclude 5130150745510468,5020439416629852
```

只显示未托管的 VPC：
```bash
./vpc-checker -unmanaged
```

发送结果到 Lark 机器人：
```bash
./vpc-checker -webhook https://open.larksuite.com/open-apis/bot/v2/hook/your-webhook-url
```

组合使用：
```bash
./vpc-checker -exclude 5130150745510468,5020439416629852 -unmanaged -webhook https://open.larksuite.com/open-apis/bot/v2/hook/your-webhook-url
```

## 系统要求

- 操作系统：
  - macOS 10.15+ (Intel 或 Apple Silicon)
  - Linux (x86_64 或 arm64)
  - Windows 10+ (64-bit)
- 内存：至少 512MB 可用内存
- 磁盘空间：约 20MB

## 权限要求

### 配置审计账号
需要具有以下权限：
- `config:ListAggregateDiscoveredResources`
- `config:GetAggregator`

### IPAM账号
需要具有以下权限：
- `vpc:ListIpamPools`
- `vpc:ListIpamPoolAllocations`

## 常见问题

1. 无法获取配置审计数据
   - 检查 AK/SK 是否正确
   - 验证账号是否有足够权限
   - 确认聚合器 ID 是否正确

2. IPAM 信息不完整
   - 检查 IPAM AK/SK 是否正确
   - 验证 IPAM 账号权限
   - 确认 IPAM 池是否正确配置

3. Webhook 发送失败
   - 验证 webhook URL 是否有效
   - 检查网络连接
   - 确认 Lark 机器人是否正常工作

## 许可证

[MIT License](LICENSE)