# Bookmark Shelf Service Docker 说明

Bookmark Shelf Service 用于在本地整理网页书签：先建立收藏夹，再保存、迁移或归档链接，并生成当前收藏夹的活跃链接统计。服务不依赖数据库或外部网络服务，运行时状态保存在内存中。

在仓库根目录执行本地构建、启动和测试：

```sh
go build ./...
go run ./cmd/bookmark-shelf
go test ./...
```

默认 HTTP 端口为 `8080`。健康检查为 `GET /healthz`，返回状态码 `200` 和 `{"status":"ok"}`；业务验收可请求 `GET /report`，或创建收藏夹后调用 `POST /bookmarks` 保存链接。

Docker 使用 `benzhi.Dockerfile`，基础镜像为 `golang:1.26.2-bookworm`。`go.mod` 固定 Go 语言版本为 `1.26.2`，Dockerfile 设置 `GOTOOLCHAIN=local`，以便容器只使用镜像中已有的工具链。镜像会复制源码、执行 `go mod download`，并在构建阶段执行 `go build ./...`。

在 Apple Silicon 或其他 arm64 环境验收：

```sh
./build_benzhi_docker.sh bookmark-shelf-service:arm64 linux/arm64
```

在 amd64 环境或支持多架构构建的 Docker 环境验收：

```sh
./build_benzhi_docker.sh bookmark-shelf-service:amd64 linux/amd64
```

脚本依次构建指定平台镜像、启动容器、在容器内执行 `go build ./...`，然后请求健康检查和报表接口。也可以逐步手工执行：

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t bookmark-shelf-service:arm64 .
docker run -d --name bookmark-shelf-arm64 -p 8080:8080 bookmark-shelf-service:arm64
docker exec bookmark-shelf-arm64 go build ./...
curl --fail http://localhost:8080/healthz
curl --fail http://localhost:8080/report
docker rm -f bookmark-shelf-arm64
```

退出码为 `0` 表示镜像构建、容器内编译和接口请求均通过；任一命令非零退出，或健康检查未返回 `200`，均表示验收未通过。
