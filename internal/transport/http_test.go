package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/service"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/transport"
)

func TestHTTPWorkflow(t *testing.T) {
	t.Parallel()
	handler := transport.NewHandler(service.NewLibrary(store.NewMemoryRepository()))

	collection := request(t, handler, http.MethodPost, "/collections", `{"name":"Reading"}`)
	if collection.Code != http.StatusCreated {
		t.Fatalf("collection status = %d", collection.Code)
	}
	bookmark := request(t, handler, http.MethodPost, "/bookmarks", `{"collection_id":"collection-001","url":"https://example.com","title":"Example","tags":["news"]}`)
	if bookmark.Code != http.StatusCreated {
		t.Fatalf("bookmark status = %d, body=%s", bookmark.Code, bookmark.Body.String())
	}
	archived := request(t, handler, http.MethodPost, "/bookmarks/bookmark-002/archive", "")
	if archived.Code != http.StatusOK {
		t.Fatalf("archive status = %d", archived.Code)
	}
	report := request(t, handler, http.MethodGet, "/report", "")
	if report.Code != http.StatusOK || !bytes.Contains(report.Body.Bytes(), []byte(`"archived_count":1`)) {
		t.Fatalf("report = %d %s", report.Code, report.Body.String())
	}
}

func request(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}
