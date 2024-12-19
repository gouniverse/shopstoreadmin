# shopstoreadmin <a href="https://gitpod.io/#https://github.com/gouniverse/shopstoreadmin" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/gouniverse/shopstoreadmin/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/shopstoreadmin/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/shopstoreadmin)](https://goreportcard.com/report/github.com/gouniverse/shopstoreadmin)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/shopstoreadmin)](https://pkg.go.dev/github.com/gouniverse/shopstoreadmin)

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Description

shopstoreadmin provides a UI for the [shopstore](https://github.com/Shopify/shopstore) library.

## Usage

```go
ui, err := shopstoreadmin.New(shopstoreadmin.NewUiOptions().
		SetLogger(&config.Logger).
		SetResponseWriter(w).
		SetRequest(r).
		SetLayout(NewShopAdminLayout()).
		SetStore(config.ShopStore))

if err != nil {
    config.Logger.Error("Error", "error", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))

    return ""
}

return ui.ToHTML()
```

```go
// ===========================================================================
// == LAYOUT
// ===========================================================================

var _ shopstoreadmin.Layout = (*shopAdminLayout)(nil)

func NewShopAdminLayout() shopstoreadmin.Layout {
	return &shopAdminLayout{}
}

type shopAdminLayout struct {
	title string
	body  string

	scriptURLs []string
	scripts    []string

	styleURLs []string
	styles    []string
}

func (a *shopAdminLayout) SetTitle(title string) {
	a.title = title
}

func (a *shopAdminLayout) SetBody(body string) {
	a.body = body
}

func (a *shopAdminLayout) SetScriptURLs(urls []string) {
	a.scriptURLs = urls
}

func (a *shopAdminLayout) SetScripts(scripts []string) {
	a.scripts = scripts
}

func (a *shopAdminLayout) SetStyleURLs(urls []string) {
	a.styleURLs = urls
}

func (a *shopAdminLayout) SetStyles(styles []string) {
	a.styles = styles
}

func (a *shopAdminLayout) Render(w http.ResponseWriter, r *http.Request) string {
	return layouts.NewAdminLayout(r, layouts.Options{
		Title:      a.title,
		Content:    hb.Raw(a.body),
		ScriptURLs: a.scriptURLs,
		Scripts:    a.scripts,
		StyleURLs:  a.styleURLs,
		Styles:     a.styles,
	}).ToHTML()
}
```