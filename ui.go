package shopstoreadmin

import (
	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
)

func New(options UiOptionsInterface) (hb.TagInterface, error) {
	err := options.Validate()

	if err != nil {
		return nil, err
	}

	// ui := ui{
	// 	response:   options.GetResponseWriter(),
	// 	request:    options.GetRequest(),
	// 	store:      options.GetStore(),
	// 	logger:     *options.GetLogger(),
	// 	layout:     options.GetLayout(),
	// 	homeURL:    options.GetHomeURL(),
	// 	websiteUrl: options.GetWebsiteUrl(),
	// }

	return handler(options), nil
}

func handler(options UiOptionsInterface) hb.TagInterface {
	controller := utils.Req(options.GetRequest(), "controller", "")

	if controller == "" {
		controller = pathHome
	}

	if controller == pathHome {
		return home(options)
	}

	if controller == pathDiscounts {
		// return visitorActivity(*ui)
	}

	if controller == pathProductCreate {
		return productCreate(options)
	}

	if controller == pathProducts {
		return productManager(options)
	}

	options.GetLayout().SetBody(hb.H1().HTML(controller).ToHTML())
	return hb.Raw(options.GetLayout().Render(options.GetResponseWriter(), options.GetRequest()))
	// redirect(a.response, a.request, url(a.request, pathQueueManager, map[string]string{}))
	// return nil
}

// type ui struct {
// 	response   http.ResponseWriter
// 	request    *http.Request
// 	store      shopstore.StoreInterface
// 	logger     slog.Logger
// 	layout     Layout
// 	homeURL    string
// 	websiteUrl string
// }

// func (ui *ui) handler() hb.TagInterface {
// 	controller := utils.Req(ui.request, "controller", "")

// 	if controller == "" {
// 		controller = pathHome
// 	}

// 	if controller == pathHome {
// 		return home(*ui)
// 	}

// 	if controller == pathDiscounts {
// 		// return visitorActivity(*ui)
// 	}

// 	if controller == pathProducts {
// 		// return visitorPaths(*ui)
// 	}

// 	ui.layout.SetBody(hb.H1().HTML(controller).ToHTML())
// 	return hb.Raw(ui.layout.Render(ui.response, ui.request))
// 	// redirect(a.response, a.request, url(a.request, pathQueueManager, map[string]string{}))
// 	// return nil
// }

// type Layout interface {
// 	SetTitle(title string)
// 	SetScriptURLs(scripts []string)
// 	SetScripts(scripts []string)
// 	SetStyleURLs(styles []string)
// 	SetStyles(styles []string)
// 	SetBody(string)
// 	Render(w http.ResponseWriter, r *http.Request) string
// }

// type UIOptions struct {
// 	ResponseWriter http.ResponseWriter
// 	Request        *http.Request
// 	Logger         *slog.Logger
// 	Store          shopstore.StoreInterface
// 	Layout         Layout
// 	HomeURL        string
// 	WebsiteUrl     string
// }

// type PageInterface interface {
// 	hb.TagInterface
// 	ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface
// }
