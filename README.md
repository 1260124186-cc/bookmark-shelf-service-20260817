# Bookmark Shelf Service

Bookmark Shelf Service is a local HTTP service for people who want to keep useful links organized without relying on an external account. It supports collection creation, bookmark saving with normalized tags, moving links between collections, archiving stale links, and a summary of active and archived material.

## Normal workflow

1. Create one or more collections for a research area, project, or reading list.
2. Save web links into a collection with a title and optional tags. A collection cannot contain the same URL twice.
3. Move active links as work changes, archive links that are no longer current, and read the report to see active counts per collection and the total archived count.

## API

Run the service:

```sh
go run ./cmd/bookmark-shelf
```

The server listens on `:8080` unless `BOOKMARK_SHELF_ADDR` is set.

```sh
curl -X POST http://localhost:8080/collections \
  -H 'Content-Type: application/json' \
  -d '{"name":"Reading"}'

curl -X POST http://localhost:8080/bookmarks \
  -H 'Content-Type: application/json' \
  -d '{"collection_id":"collection-001","url":"https://go.dev/doc/","title":"Go documentation","tags":["go","reference"]}'

curl http://localhost:8080/report
```

`GET /healthz` returns `{"status":"ok"}` for health checks.
