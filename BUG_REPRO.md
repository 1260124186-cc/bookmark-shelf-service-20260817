# 修复前故障复现（Docker）

## 项目与标准命令

Bookmark Shelf Service 是一个用于按集合保存、移动和归档书签的本地 HTTP 服务。请在仓库根目录执行以下标准命令：

```sh
go build ./...
go run ./cmd/bookmark-shelf
go test ./...
```

## 环境构建与编译

已实际执行以下命令：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t bookmark-shelf-bug003:base-linux-amd64 .
docker run --rm --platform linux/amd64 bookmark-shelf-bug003:base-linux-amd64 go build ./...
docker build --platform linux/arm64 -f benzhi.Dockerfile -t bookmark-shelf-bug003:base-linux-arm64 .
docker run --rm --platform linux/arm64 bookmark-shelf-bug003:base-linux-arm64 go build ./...
```

linux/amd64 和 linux/arm64 的镜像构建及容器内编译均成功。目标故障在下面的测试命令中触发。

## 故障触发步骤

在仓库根目录执行：

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t bookmark-shelf-bug003:base-linux-arm64 .
docker run --rm --platform linux/arm64 bookmark-shelf-bug003:base-linux-arm64 go test ./...
```

## 实际错误输出

```text
?   	github.com/1260124186-cc/bookmark-shelf-service-20260817/cmd/bookmark-shelf	[no test files]
?   	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain	[no test files]
--- FAIL: TestLibraryKeepsNormalizedTagsAfterCallerReusesInput (0.00s)
    library_test.go:88: unexpected returned tags: []string{"changed", "reference", "go"}
FAIL
FAIL	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/service	0.002s
--- FAIL: TestMemoryRepositoryDoesNotRetainCallerTagSlice (0.00s)
    memory_test.go:38: stored tags changed with caller input: []string{"changed"}
FAIL
FAIL	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store	0.002s
ok  	github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/transport	0.002s
FAIL

exit status: 1
```

## 期望行为

保存书签后，调用方后续修改原标签列表不应改变已保存书签的标签；标签应保持去除空白、统一大小写并去重后的内容。
