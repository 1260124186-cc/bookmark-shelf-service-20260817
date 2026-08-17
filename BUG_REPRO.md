# 修复前故障复现（Docker）

## 项目与标准命令

Bookmark Shelf Service 是一个用于创建收藏夹、保存链接、移动或归档书签并生成汇总报告的本地 HTTP 服务。在仓库根目录可执行：

```sh
go build ./...
go run ./cmd/bookmark-shelf
go test ./...
```

## 环境构建与编译

已实际执行下列 linux/amd64 与 linux/arm64 镜像构建命令，并在每个镜像内执行 `go build ./...`，均成功：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t bookmark-shelf-bug001-base-amd64 .
docker run --rm --platform linux/amd64 bookmark-shelf-bug001-base-amd64 go build ./...
docker build --platform linux/arm64 -f benzhi.Dockerfile -t bookmark-shelf-bug001-base-arm64 .
docker run --rm --platform linux/arm64 bookmark-shelf-bug001-base-arm64 go build ./...
```

## 故障触发步骤

在仓库根目录执行以下命令：

```sh
docker run --rm --platform linux/arm64 -v "$PWD":/app -w /app golang:1.26.2-bookworm go test ./...
```

## 实际错误输出

```text
?   	github.com/1260124186-cc/bookmark-shelf-service-20260817/cmd/bookmark-shelf	[no test files]
?   	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain	[no test files]
--- FAIL: TestLibraryRejectsDuplicateWithinCollection (0.00s)
    library_test.go:64: expected duplicate error, got collection not found
FAIL
FAIL	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/service	0.004s
?   	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store	[no test files]
--- FAIL: TestHTTPDistinguishesDuplicateAndMissingCollection (0.00s)
    http_test.go:50: missing collection status = 409
FAIL
FAIL	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/transport	0.022s
FAIL
命令退出状态：1
```

## 期望行为

同一收藏夹再次保存相同链接时，客户端应收到冲突响应；向不存在的收藏夹保存链接时，客户端应收到资源不存在响应。其他正常书签保存请求应保持创建成功。
