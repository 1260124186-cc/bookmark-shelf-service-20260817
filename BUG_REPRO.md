# 修复前故障复现（Docker）

## 项目与标准命令

Bookmark Shelf Service 是一个本地 HTTP 服务，用于创建收藏夹、保存和归档书签，并查询书签统计信息。在仓库根目录可执行以下标准命令：

```sh
go build ./...
go run ./cmd/bookmark-shelf
go test ./...
```

## 环境构建与编译

已实际执行以下 `linux/arm64` 命令，镜像构建和容器内编译均成功：

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t bookmark-shelf-bug004-base-arm64 .
docker run --rm --platform linux/arm64 bookmark-shelf-bug004-base-arm64 go build ./...
```

已实际执行以下 `linux/amd64` 命令，镜像构建和容器内编译均成功：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t bookmark-shelf-bug004-base-amd64 .
docker run --rm --platform linux/amd64 bookmark-shelf-bug004-base-amd64 go build ./...
```

两个平台的镜像构建、容器内 `go build ./...`、健康检查和报表请求均成功；目标故障由下节命令触发。

## 故障触发步骤

在仓库根目录执行：

```sh
go test ./internal/transport -run '^TestHTTPReturnsNotFoundForUnknownBookmarkArchive$' -count=1
```

## 实际错误输出

```text
--- FAIL: TestHTTPReturnsNotFoundForUnknownBookmarkArchive (0.00s)
    http_test.go:41: archive missing status = 200
FAIL
FAIL	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/transport	1.581s
FAIL
```

命令退出码：1。

## 期望行为

客户端归档不存在的书签时，应收到 `404 Not Found` 错误响应，而不是 `200 OK` 和不完整的成功响应；调用方应能据此判断归档未完成。
