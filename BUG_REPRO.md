# 修复前故障复现（Docker）

## 项目与标准命令

Bookmark Shelf Service 是一个用于创建收藏夹、保存和归档书签，并生成书签统计报表的本地 HTTP 服务。在仓库根目录执行以下标准命令：

```sh
go build ./...
go run ./cmd/bookmark-shelf
go test ./...
```

## 环境构建与编译

已实际构建 linux/arm64 和 linux/amd64 镜像，并在对应容器内执行 `go build ./...`，两种平台的镜像构建和容器内编译均成功。

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t bookmark-shelf-service-bug-002-base:arm64 .
docker run --rm --platform linux/arm64 bookmark-shelf-service-bug-002-base:arm64 go build ./...
docker build --platform linux/amd64 -f benzhi.Dockerfile -t bookmark-shelf-service-bug-002-base:amd64 .
docker run --rm --platform linux/amd64 bookmark-shelf-service-bug-002-base:amd64 go build ./...
```

## 故障触发步骤

在仓库根目录执行：

```sh
go test ./internal/service -run TestLibraryStopsReportWhenRequestIsCanceled -count=1
```

## 实际错误输出

```text
--- FAIL: TestLibraryStopsReportWhenRequestIsCanceled (0.00s)
    library_test.go:79: expected canceled report, got <nil>
FAIL
FAIL	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/service	2.173s
FAIL
```

## 期望行为

当调用方传入的请求已取消时，报表操作应停止并返回取消错误，而不是继续返回完整统计结果。
