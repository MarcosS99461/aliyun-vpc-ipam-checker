#!/bin/bash

# 版本号
VERSION="v1.0.0"
BINARY_NAME="vpc-checker"

# 创建发布目录
mkdir -p releases

# 编译不同平台的二进制文件
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o releases/${BINARY_NAME}-${VERSION}-darwin-amd64 ./cmd/vpc-checker

# macOS (Apple Silicon/M1)
GOOS=darwin GOARCH=arm64 go build -o releases/${BINARY_NAME}-${VERSION}-darwin-arm64 ./cmd/vpc-checker

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o releases/${BINARY_NAME}-${VERSION}-linux-amd64 ./cmd/vpc-checker

# Linux (arm64)
GOOS=linux GOARCH=arm64 go build -o releases/${BINARY_NAME}-${VERSION}-linux-arm64 ./cmd/vpc-checker

# Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o releases/${BINARY_NAME}-${VERSION}-windows-amd64.exe ./cmd/vpc-checker

# 为每个二进制文件创建压缩包
cd releases
for file in *; do
    if [ -f "$file" ]; then
        tar -czf "${file}.tar.gz" "$file"
        rm "$file"
    fi
done

echo "Build complete! Release files are in the releases directory."
