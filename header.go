package shopstoreadmin

import (
	// "project/config"
	// "project/internal/links"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore"
	"github.com/spf13/cast"
)

func Header(opts UiOptionsInterface) hb.TagInterface {
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
