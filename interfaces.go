package shopstoreadmin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/shopstore"
)

type Layout interface {
	SetTitle(title string)
	SetScriptURLs(scripts []string)
	SetScripts(scripts []string)
	SetStyleURLs(styles []string)
	SetStyles(styles []string)
	SetBody(string)
	Render(w http.ResponseWriter, r *http.Request) string
}

type UiOptionsInterface interface {
	GetResponseWriter() http.ResponseWriter
	SetResponseWriter(http.ResponseWriter) UiOptionsInterface

	GetRequest() *http.Request
	SetRequest(*http.Request) UiOptionsInterface

	GetStore() shopstore.StoreInterface
	SetStore(shopstore.StoreInterface) UiOptionsInterface

	GetLogger() *slog.Logger
	SetLogger(*slog.Logger) UiOptionsInterface

	GetLayout() Layout
	SetLayout(Layout) UiOptionsInterface

	GetHomeURL() string
	SetHomeURL(string) UiOptionsInterface

	GetWebsiteUrl() string
	SetWebsiteUrl(string) UiOptionsInterface

	Validate() error
}
