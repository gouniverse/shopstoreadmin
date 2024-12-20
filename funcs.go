package shopstoreadmin

import (
	"net/http"
	"strings"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore"
	"github.com/samber/lo"
	"github.com/spf13/cast"

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

func header(opts UiOptionsInterface) hb.TagInterface {
	urlHome := opts.GetHomeURL()
	urlShop := url(opts.GetRequest(), pathHome, nil)
	urlOrders := url(opts.GetRequest(), pathOrders, nil)
	urlProducts := url(opts.GetRequest(), pathProducts, nil)
	urlDiscounts := url(opts.GetRequest(), pathDiscounts, nil)

	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(urlHome).
		Class("nav-link")

	linkShop := hb.NewHyperlink().
		HTML("Shop").
		Href(urlShop).
		Class("nav-link")

	linkOrders := hb.Hyperlink().
		HTML("Orders").
		Href(urlOrders).
		Class("nav-link")

	linkDiscounts := hb.Hyperlink().
		HTML("Discounts").
		Href(urlDiscounts).
		Class("nav-link")

	linkProducts := hb.Hyperlink().
		HTML("Products ").
		Href(urlProducts).
		Class("nav-link")

	productsCount, err := opts.GetStore().ProductCount(opts.GetRequest().Context(), shopstore.NewProductQuery())

	if err != nil {
		opts.GetLogger().Error(err.Error())
		productsCount = -1
	}

	ordersCount, err := opts.GetStore().OrderCount(opts.GetRequest().Context(), shopstore.NewOrderQuery())

	if err != nil {
		opts.GetLogger().Error(err.Error())
		ordersCount = -1
	}

	discountsCount, err := opts.GetStore().DiscountCount(opts.GetRequest().Context(), shopstore.NewDiscountQuery())

	if err != nil {
		opts.GetLogger().Error(err.Error())
		discountsCount = -1
	}

	ulNav := hb.NewUL().
		Class("nav  nav-pills justify-content-center").
		Child(hb.NewLI().
			Class("nav-item").Child(linkHome)).
		Child(hb.NewLI().
			Class("nav-item").Child(linkShop)).
		Child(hb.LI().
			Class("nav-item").
			Child(linkOrders.
				Child(hb.Span().
					Class("badge bg-secondary ms-2").
					HTML(cast.ToString(ordersCount))))).
		Child(hb.LI().
			Child(linkProducts.
				Child(hb.Span().
					Class("badge bg-secondary ms-2").
					HTML(cast.ToString(productsCount))))).
		Child(hb.LI().
			Child(linkDiscounts.
				Child(hb.Span().
					Class("badge bg-secondary ms-2").
					HTML(cast.ToString(discountsCount)))))

	divCard := hb.NewDiv().Class("card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Class("card-body").Style("padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav))
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

// isDate checks if a string is a valid date
// the format is YYYY-MM-DD
//
// Business logic:
// - checks if the string contains 2 dashes
// - checks if the string does not contain colons
// - checks the first dash is at position 4
// - checks the second dash is at position 7
func isDate(value string) bool {
	countDashes := strings.Count(value, "-")

	if countDashes != 2 {
		return false
	}

	countColons := strings.Count(value, ":")

	if countColons > 0 {
		return false
	}

	if strings.Index(value, "-") != 4 {
		return false
	}

	if strings.LastIndex(value, "-") != 7 {
		return false
	}

	return true
}

// isDateTime checks if a string is a valid datetime
// the format is YYYY-MM-DD HH:MM:SS
//
// Business logic:
// - checks if the string contains 2 dashes
// - checks if the string contains 2 colons
// - checks the first dash is at position 4
// - checks the second dash is at position 7
// - checks the first colon is at position 10
// - checks the second colon is at position 13
func isDateTime(value string) bool {
	countDashes := strings.Count(value, "-")

	if countDashes != 2 {
		return false
	}

	countColons := strings.Count(value, ":")

	if countColons != 2 {
		return false
	}

	if strings.Index(value, "-") != 4 {
		return false
	}

	if strings.LastIndex(value, "-") != 7 {
		return false
	}

	if strings.Index(value, ":") != 13 {
		return false
	}

	if strings.LastIndex(value, ":") != 16 {
		return false
	}

	return true
}
