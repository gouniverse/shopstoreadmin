package shopstoreadmin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gouniverse/shopstore"
)

func NewUiOptions() UiOptionsInterface {
	return &uiOptionsImplementation{}
}

type uiOptionsImplementation struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Store          shopstore.StoreInterface
	Logger         *slog.Logger
	Layout         Layout
	HomeURL        string
	WebsiteUrl     string
}

func (options *uiOptionsImplementation) Validate() error {
	if options.ResponseWriter == nil {
		return errors.New("options.ResponseWriter is required")
	}

	if options.Request == nil {
		return errors.New("options.Request is required")
	}

	if options.Store == nil {
		return errors.New("options.Store is required")
	}

	if options.Logger == nil {
		return errors.New("options.Logger is required")
	}

	if options.Layout == nil {
		return errors.New("options.Layout is required")
	}

	return nil
}

func (options *uiOptionsImplementation) GetResponseWriter() http.ResponseWriter {
	return options.ResponseWriter
}

func (options *uiOptionsImplementation) SetResponseWriter(responseWriter http.ResponseWriter) UiOptionsInterface {
	options.ResponseWriter = responseWriter
	return options
}

func (options *uiOptionsImplementation) GetRequest() *http.Request {
	return options.Request
}

func (options *uiOptionsImplementation) SetRequest(request *http.Request) UiOptionsInterface {
	options.Request = request
	return options
}

func (options *uiOptionsImplementation) GetStore() shopstore.StoreInterface {
	return options.Store
}

func (options *uiOptionsImplementation) SetStore(store shopstore.StoreInterface) UiOptionsInterface {
	options.Store = store
	return options
}

func (options *uiOptionsImplementation) GetLogger() *slog.Logger {
	return options.Logger
}

func (options *uiOptionsImplementation) SetLogger(logger *slog.Logger) UiOptionsInterface {
	options.Logger = logger
	return options
}

func (options *uiOptionsImplementation) GetLayout() Layout {
	return options.Layout
}

func (options *uiOptionsImplementation) SetLayout(layout Layout) UiOptionsInterface {
	options.Layout = layout
	return options
}

func (options *uiOptionsImplementation) GetHomeURL() string {
	return options.HomeURL
}

func (options *uiOptionsImplementation) SetHomeURL(homeURL string) UiOptionsInterface {
	options.HomeURL = homeURL
	return options
}

func (options *uiOptionsImplementation) GetWebsiteUrl() string {
	return options.WebsiteUrl
}

func (options *uiOptionsImplementation) SetWebsiteUrl(websiteUrl string) UiOptionsInterface {
	options.WebsiteUrl = websiteUrl
	return options
}
