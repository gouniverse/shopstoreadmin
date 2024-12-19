package shopstoreadmin

import (
	"net/http"

	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore/admin/shared"
	"github.com/samber/lo"
)

// =============================================================================
// == CONSTRUCTOR
// =============================================================================

func home(options UiOptionsInterface) shared.PageInterface {
	return &homeController{
		opts: options,
	}
}

// =============================================================================
// == CONTROLLER
// =============================================================================

type homeController struct {
	opts UiOptionsInterface
}

type homeControllerData struct{}

func (c *homeController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData()

	c.opts.GetLayout().SetTitle("Dashboard | Shop")

	if errorMessage != "" {
		c.opts.GetLayout().SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.opts.GetLayout().Render(w, r))
	}

	htmxScript := `setTimeout(() => async function() {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
			await script.onload
		}
	}, 1000);`

	swalScript := `setTimeout(() => async function() {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
			await script.onload
		}
	}, 1000);`

	// cdn.Jquery_3_7_1(),
	// // `https://cdnjs.cloudflare.com/ajax/libs/Chart.js/1.0.2/Chart.min.js`,
	// `https://cdn.jsdelivr.net/npm/chart.js`,

	c.opts.GetLayout().SetBody(c.page(data).ToHTML())
	c.opts.GetLayout().SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(c.opts.GetLayout().Render(w, r))
}

func (c *homeController) ToHTML() string {
	return c.ToTag(nil, nil).ToHTML()
}

// == PRIVATE METHODS ==========================================================

func (c *homeController) prepareData() (data homeControllerData, errorMessage string) {
	return homeControllerData{}, ""
}

func (c *homeController) page(_ homeControllerData) hb.TagInterface {
	breadcrumbs := breadcrumbs(c.opts.GetRequest(), []breadcrumb{})

	title := hb.Heading1().
		HTML("Shop. Home")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(Header(c.opts)).
		Child(hb.HR()).
		Child(title).
		Child(hb.BR()).
		Child(c.tiles())
}

// == PRIVATE METHODS ==========================================================

func (c *homeController) tiles() hb.TagInterface {
	ordersURL := url(c.opts.GetRequest(), pathOrders, nil)
	productsURL := url(c.opts.GetRequest(), pathProducts, nil)
	discountsURL := url(c.opts.GetRequest(), pathDiscounts, nil)

	tiles := []map[string]string{
		{
			"title": "Order Manager",
			"icon":  "bi-cart",
			"link":  ordersURL,
		},
		{
			"title": "Product Manager",
			"icon":  "bi-box",
			"link":  productsURL,
		},
		{
			"title": "Discount Manager",
			"icon":  "bi-percent",
			"link":  discountsURL,
		},
	}

	cards := lo.Map(tiles, func(tile map[string]string, index int) hb.TagInterface {
		target := lo.ValueOr(tile, "target", "")
		card := hb.Div().
			Class("card").
			Class("bg-transparent border round-10 shadow-lg h-100 pt-4").
			OnMouseOver(`
			this.style.setProperty('background-color', 'beige', 'important');
			this.style.setProperty('scale', 1.1);
			this.style.setProperty('border', '4px solid moccasin', 'important');
			`).
			OnMouseOut(`
			this.style.setProperty('background-color', 'transparent', 'important');
			this.style.setProperty('scale', 1);
			this.style.setProperty('border', '0px solid moccasin', 'important');
			`).
			Style("margin:0px 0px 20px 0px;").
			Children([]hb.TagInterface{
				hb.Div().Class("card-body").
					Class("d-flex flex-column justify-content-evenly").
					Children([]hb.TagInterface{
						hb.Div().
							Child(hb.I().Class(`bi ` + tile["icon"])).Style(`font-size:36px;color: red;`).
							Style("text-align:center;padding:10px;"),
						hb.Heading5().
							HTML(tile["title"]).
							Style("text-align:center;padding:10px;"),
					}),
			})

		link := hb.Hyperlink().
			Href(tile["link"]).
			AttrIf(target != "", "target", target).
			Child(card)

		column := hb.Div().
			Class("col-sm-6 col-md-4 col-lg-3").
			Child(link)

		return column
	})

	return hb.Div().Class("row").Children(cards)
}
