package shopstoreadmin

import (
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/samber/lo"

	urlpkg "net/url"
)

func breadcrumbs(r *http.Request, pageBreadcrumbs []breadcrumb) hb.TagInterface {
	adminHomeURL := "/admin" //AdminHomeURL(r)
	//path := utils.Req(r, "path", "")

	adminHomeBreadcrumb := lo.
		If(adminHomeURL != "", breadcrumb{
			Name: "Home",
			URL:  adminHomeURL,
		}).
		Else(breadcrumb{})

	breadcrumbItems := []breadcrumb{
		adminHomeBreadcrumb,
		{
			Name: "Shop",
			URL:  url(r, pathHome, nil),
		},
	}

	breadcrumbItems = append(breadcrumbItems, pageBreadcrumbs...)

	breadcrumbs := breadcrumbsUI(breadcrumbItems)

	return hb.Div().
		Child(breadcrumbs)
}

// func redirect(w http.ResponseWriter, r *http.Request, url string) string {
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// 	http.Redirect(w, r, url, http.StatusSeeOther)
// 	return ""
// }

type breadcrumb struct {
	Name string
	URL  string
}

func breadcrumbsUI(breadcrumbs []breadcrumb) hb.TagInterface {

	ol := hb.OL().
		Class("breadcrumb").
		Style("margin-bottom: 0px;")

	for _, breadcrumb := range breadcrumbs {

		link := hb.Hyperlink().
			HTML(breadcrumb.Name).
			Href(breadcrumb.URL)

		li := hb.LI().
			Class("breadcrumb-item").
			Child(link)

		ol.AddChild(li)
	}

	nav := hb.Nav().
		Class("d-inline-block").
		Attr("aria-label", "breadcrumb").
		Child(ol)

	return nav
}

// func redirect(w http.ResponseWriter, r *http.Request, url string) string {
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// 	http.Redirect(w, r, url, http.StatusSeeOther)
// 	return ""
// }

func url(r *http.Request, path string, params map[string]string) string {
	endpoint := r.URL.Path

	if params == nil {
		params = map[string]string{}
	}

	params["controller"] = path

	url := endpoint + query(params)

	return url
}

func query(queryData map[string]string) string {
	queryString := ""

	if len(queryData) > 0 {
		v := urlpkg.Values{}
		for key, value := range queryData {
			v.Set(key, value)
		}
		queryString += "?" + httpBuildQuery(v)
	}

	return queryString
}

func httpBuildQuery(queryData urlpkg.Values) string {
	return queryData.Encode()
}
