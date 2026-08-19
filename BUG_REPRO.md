# 修复前故障复现（Docker）

## 项目与标准命令

Bookmark Shelf Service 是一个本地书签整理 HTTP 服务，支持创建收藏夹、保存链接、移动和归档书签，并生成汇总报告。在仓库根目录可执行：

```sh
go build ./...
go run ./cmd/bookmark-shelf
go test ./...
```

## 环境构建与编译

已实际执行下列 linux/amd64 与 linux/arm64 镜像构建命令，并在每个镜像内执行 `go build ./...`，均成功：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t bookmark-shelf-bug005-base-amd64 .
docker run --rm --platform linux/amd64 bookmark-shelf-bug005-base-amd64 go build ./...
docker build --platform linux/arm64 -f benzhi.Dockerfile -t bookmark-shelf-bug005-base-arm64 .
docker run --rm --platform linux/arm64 bookmark-shelf-bug005-base-arm64 go build ./...
```

两个平台的镜像构建和容器内编译均成功；目标故障由下节命令触发。

## 故障触发步骤

在仓库根目录执行以下命令：

```sh
docker run --rm --platform linux/arm64 bookmark-shelf-bug005-base-arm64 sh -c 'go test -count=1 -run TestArchiveMissDoesNotBlockLaterWrites ./internal/store; go test -count=1 -run TestArchiveMissDoesNotBlockSavingBookmark ./internal/store'
```

## 实际错误输出

```text
+ go test -count=1 -run TestArchiveMissDoesNotBlockLaterWrites ./internal/store
--- FAIL: TestArchiveMissDoesNotBlockLaterWrites (0.10s)
    archive_lock_test.go:29: later write remained blocked after archive miss
FAIL
FAIL	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store	0.105s
FAIL
+ go test -count=1 -run TestArchiveMissDoesNotBlockSavingBookmark ./internal/store
--- FAIL: TestArchiveMissDoesNotBlockSavingBookmark (0.11s)
    archive_lock_test.go:57: bookmark save remained blocked after archive miss
FAIL
FAIL	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store	0.108s
FAIL
```

## 期望行为

归档不存在的书签时应返回找不到书签的错误；随后创建收藏夹和保存链接的请求均应及时完成，不应被遗留的锁阻塞。
