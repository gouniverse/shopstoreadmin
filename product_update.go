package shopstoreadmin

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gouniverse/base/req"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const viewContent = "content"
const viewMetadata = "metadata"
const viewMedia = "media"
const viewSettings = "settings"

// ===========================================================================
// == CONSTRUCTOR
// ===========================================================================

func productUpdate(opts UiOptionsInterface) pageInterface {
	return &productUpdateController{
		opts: opts,
	}
}

// ===========================================================================
// == CONTROLLER
// ===========================================================================

type productUpdateController struct {
	opts UiOptionsInterface
}

// ===========================================================================
// == INTERFACE IMPLEMENTATION
// ===========================================================================

func (c *productUpdateController) ToTag() hb.TagInterface {
	data, errorMessage := c.prepareDataAndValidate()

	if errorMessage != "" {
		return hb.Div().Class("alert alert-danger").Child(hb.Text(errorMessage))
	}

	if c.opts.GetRequest().Method == http.MethodPost {
		return c.form(data)
	}

	c.opts.GetLayout().SetTitle("Edit Product | Shop")
	c.opts.GetLayout().SetBody(c.page(data).ToHTML())
	c.opts.GetLayout().SetStyleURLs([]string{
		cdn.TrumbowygCss_2_27_3(),
	})
	c.opts.GetLayout().SetScriptURLs([]string{
		cdn.Htmx_2_0_0(),
		cdn.Sweetalert2_10(),
		cdn.Jquery_3_7_1(),
		cdn.TrumbowygJs_2_27_3(),
	})
	c.opts.GetLayout().SetStyles([]string{
		`
.htmx-indicator {
    display: none;
}
.htmx-request.htmx-indicator {
    display: inline-block;
}
		`,
	})

	return hb.Raw(c.opts.GetLayout().Render(c.opts.GetResponseWriter(), c.opts.GetRequest()))
}

func (controller *productUpdateController) ToHTML() string {
	return controller.ToTag().ToHTML()
}

// ===========================================================================
// == METHODS
// ===========================================================================

func (c *productUpdateController) page(data productUpdateControllerData) hb.TagInterface {
	productManegerURL := url(c.opts.GetRequest(), pathProducts, map[string]string{})

	productUpdateURL := url(c.opts.GetRequest(), pathProductUpdate, map[string]string{
		"productID": data.productID,
	})

	breadcrumbs := breadcrumbs(c.opts.GetRequest(), []breadcrumb{
		{
			Name: "Product Manager",
			URL:  productManegerURL,
		},
		{
			Name: "Edit Product",
			URL:  productUpdateURL,
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(productManegerURL)

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		Child(hb.Div().ID("ButtonSaveIndicator").Class("spinner-border spinner-border-sm ms-2 htmx-indicator")).
		HxIndicator("#ButtonSaveIndicator").
		HxInclude("#FormProductUpdate").
		HxPost(productUpdateURL).
		HxTarget("#FormProductUpdate")

	title := hb.Heading1().
		Class("mb-3").
		Text("Edit Product: ").
		Text(data.product.Title()).
		Child(buttonCancel)

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTMLIf(data.view == viewContent, "Product Contents").
					HTMLIf(data.view == viewMetadata, "Product Metadata").
					HTMLIf(data.view == viewSettings, "Product Settings").
					Style("margin-bottom:0;display:inline-block;")).
				Child(buttonSave),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(c.form(data)))

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(header(c.opts)).
		Child(hb.HR()).
		Child(title).
		Child(c.tabs(data)).
		Child(card).
		// Child(cardProductDetails).
		Child(hb.BR()).
		Child(hb.BR()).
		Child(hb.BR()).
		Child(hb.BR())
	// Child(cardProductMetadata)
}

func (controller *productUpdateController) tabs(data productUpdateControllerData) hb.TagInterface {
	viewContentURL := url(controller.opts.GetRequest(), pathProductUpdate, map[string]string{
		"product_id": data.productID,
		"view":       viewContent,
	})

	viewMedia := url(controller.opts.GetRequest(), pathProductUpdate, map[string]string{
		"product_id": data.productID,
		"view":       viewMedia,
	})

	viewMetadataURL := url(controller.opts.GetRequest(), pathProductUpdate, map[string]string{
		"product_id": data.productID,
		"view":       viewMetadata,
	})

	viewSettingsURL := url(controller.opts.GetRequest(), pathProductUpdate, map[string]string{
		"product_id": data.productID,
		"view":       viewSettings,
	})

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == viewContent, "active").
				Href(viewContentURL).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == viewMedia, "active").
				Href(viewMedia).
				HTML("Media"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == viewMetadata, "active").
				Href(viewMetadataURL).
				HTML("Metas"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == viewSettings, "active").
				Href(viewSettingsURL).
				HTML("Settings")))

	return tabs
}

func (c *productUpdateController) form(data productUpdateControllerData) hb.TagInterface {
	formProductUpdate := form.NewForm(form.FormOptions{
		ID: "FormProductUpdate",
	})

	if data.view == viewContent {
		formProductUpdate.SetFields(c.formContentFields(data))
	}

	if data.view == viewMetadata {
		formProductUpdate.SetFields(c.formMetadataFields(data))
	}

	if data.view == viewMedia {
		formProductUpdate.SetFields(c.formMediaFields(data))
	}

	if data.view == viewSettings {
		formProductUpdate.SetFields(c.formSettingsFields(data))
	}

	if data.formErrorMessage != "" {
		formProductUpdate.AddField(form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
		}))
	}

	if data.formSuccessMessage != "" {
		formProductUpdate.AddField(form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}).ToHTML(),
		}))
	}

	productID := form.NewField(form.FieldOptions{
		Label:     "Product ID",
		Name:      "product_id",
		Type:      form.FORM_FIELD_TYPE_STRING,
		Value:     data.productID,
		Readonly:  true,
		Invisible: true,
	})

	view := form.NewField(form.FieldOptions{
		Label:     "View",
		Name:      "view",
		Type:      form.FORM_FIELD_TYPE_HIDDEN,
		Value:     data.view,
		Readonly:  true,
		Invisible: true,
	})

	formProductUpdate.AddField(productID)
	formProductUpdate.AddField(view)

	return formProductUpdate.Build()

}

func (c *productUpdateController) formContentFields(data productUpdateControllerData) []form.FieldInterface {
	title := form.NewField(form.FieldOptions{
		Label: "Title",
		Name:  "product_title",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formTitle,
		Help:  `The title of the product.`,
	})

	description := form.NewField(form.FieldOptions{
		Label: "Description",
		Name:  "product_description",
		Type:  form.FORM_FIELD_TYPE_HTMLAREA,
		Value: data.formDescription,
		Help:  `The description of the product.`,
	})

	shortDescription := form.NewField(form.FieldOptions{
		Label: "Short Description",
		Name:  "product_short_description",
		Type:  form.FORM_FIELD_TYPE_HTMLAREA,
		Value: data.formShortDescription,
		Help:  `The short description of the product.`,
	})

	return []form.FieldInterface{
		title,
		shortDescription,
		description,
	}
}

func (c *productUpdateController) formMediaFields(data productUpdateControllerData) []form.FieldInterface {
	repeaterAddURL := url(c.opts.GetRequest(), pathProductUpdate, map[string]string{
		"product_id": data.productID,
		"view":       viewMedia,
		"action":     "add",
	})

	repeaterMoveUpURL := url(c.opts.GetRequest(), pathProductUpdate, map[string]string{
		"product_id": data.productID,
		"view":       viewMedia,
		"action":     "move_up",
	})

	repeaterMoveDownURL := url(c.opts.GetRequest(), pathProductUpdate, map[string]string{
		"product_id": data.productID,
		"view":       viewMedia,
		"action":     "move_down",
	})

	repeaterRemoveURL := url(c.opts.GetRequest(), pathProductUpdate, map[string]string{
		"product_id": data.productID,
		"view":       viewMedia,
		"action":     "remove",
	})

	fieldID := form.NewField(form.FieldOptions{
		ID:    "product_media_id",
		Label: "ID",
		Name:  "id",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Help:  `The ID of the media.`,
	})

	fieldTitle := form.NewField(form.FieldOptions{
		ID:    "product_media_title",
		Label: "Title",
		Name:  "title",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Help:  `The Title of the media.`,
	})

	fieldURL := form.NewField(form.FieldOptions{
		ID:    "product_media_url",
		Label: "URL",
		Name:  "url",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Help:  `The URL of the media.`,
	})

	repeater := form.NewRepeater(form.RepeaterOptions{
		Label:               "Media",
		Help:                `The media of the product.`,
		Name:                "product_media",
		Values:              data.formMedia,
		RepeaterAddUrl:      repeaterAddURL,
		RepeaterMoveUpUrl:   repeaterMoveUpURL,
		RepeaterMoveDownUrl: repeaterMoveDownURL,
		RepeaterRemoveUrl:   repeaterRemoveURL,
		Fields: []form.FieldInterface{
			fieldID,
			fieldTitle,
			fieldURL,
		},
	})

	return []form.FieldInterface{
		repeater,
	}
}

func (c *productUpdateController) formMetadataFields(data productUpdateControllerData) []form.FieldInterface {
	metas := data.formMetas

	fields := []form.FieldInterface{}

	index := 0
	keys := lo.Keys(metas)
	slices.Sort(keys)
	for _, key := range keys {
		value := metas[key]
		background := lo.Ternary(index%2 == 0, "bg-light", "bg-white")
		fieldsMeta := []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Help:  `Opening row`,
				Value: `<div id="Row` + cast.ToString(index) + `" class="row ` + background + ` py-2">`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Help:  `Opening column 1`,
				Value: `<div class="col-3">`,
			}),
			form.NewField(form.FieldOptions{
				Label: `Key`,
				Name:  `product_meta[` + cast.ToString(index) + `][key]`,
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: key,
				// Help:  "The metadata value.",
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Help:  `Closing column 1`,
				Value: `</div>`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Help:  `Opening column 2`,
				Value: `<div class="col-8">`,
			}),
			form.NewField(form.FieldOptions{
				Label: `Value`,
				Name:  `product_meta[` + cast.ToString(index) + `][value]`,
				Type:  form.FORM_FIELD_TYPE_TEXTAREA,
				Value: value,
				// Help:  "The metadata value.",
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Help:  `Closing column 2`,
				Value: `</div>`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Help:  `Opening column 3`,
				Value: `<div class="col-1">`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: `<button onclick="document.getElementById('Row` + cast.ToString(index) + `').innerHTML='';" type="button" class="btn btn-sm btn-danger">x</button>`,
				Help:  "Delete...",
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Help:  `Closing column 3`,
				Value: `</div>`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Help:  `Closing the row.`,
				Value: `</div>`,
			}),
		}

		fields = append(fields, fieldsMeta...)

		index++
	}

	fieldsNewMeta := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: `<hr />`,
		}),
		form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: `<div class="row bg-info py-2">`,
		}),
		form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: `<h3>New Meta</h3>`,
		}),
		form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: `<div class="col-6">`,
		}),

		form.NewField(form.FieldOptions{
			Label: `Key`,
			Name:  `product_meta[` + cast.ToString(index) + `][key]`,
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: "",
			// Help:  "The metadata value.",
		}),

		form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: `</div>`,
		}),

		form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: `<div class="col-6">`,
		}),

		form.NewField(form.FieldOptions{
			Label: `Value`,
			Name:  `product_meta[` + cast.ToString(index) + `][value]`,
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: "",
			// Help:  "The metadata value.",
		}),

		form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: `</div>`,
		}),

		form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: `</div>`,
		}),
	}

	fields = append(fields, fieldsNewMeta...)

	return fields
}

func (c *productUpdateController) formSettingsFields(data productUpdateControllerData) []form.FieldInterface {
	status := form.NewField(form.FieldOptions{
		Label: "Status",
		Name:  "product_status",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formStatus,
		Help:  `The status of the product.`,
		Options: []form.FieldOption{
			{
				Value: "- not selected -",
				Key:   "",
			},
			{
				Value: "Active",
				Key:   shopstore.PRODUCT_STATUS_ACTIVE,
			},
			{
				Value: "Disabled",
				Key:   shopstore.PRODUCT_STATUS_DISABLED,
			},
			{
				Value: "Draft",
				Key:   shopstore.PRODUCT_STATUS_DRAFT,
			},
		},
	})

	price := form.NewField(form.FieldOptions{
		Label: "Price",
		Name:  "product_price",
		Type:  form.FORM_FIELD_TYPE_NUMBER,
		Value: data.formPrice,
		Help:  `The price of the product.`,
	})

	quantity := form.NewField(form.FieldOptions{
		Label: "Quantity",
		Name:  "product_quantity",
		Type:  form.FORM_FIELD_TYPE_NUMBER,
		Value: data.formQuantity,
		Help:  `The quantity of the product that is available to purchase.`,
	})

	memo := form.NewField(form.FieldOptions{
		Label: "Admin Notes",
		Name:  "product_memo",
		Type:  form.FORM_FIELD_TYPE_TEXTAREA,
		Value: data.formMemo,
		Help:  "Admin notes for this product. These notes will not be visible to the public.",
	})

	return []form.FieldInterface{
		status,
		price,
		quantity,
		memo,
	}
}

func (c *productUpdateController) saveProductSettings(data productUpdateControllerData) (d productUpdateControllerData, errorMessage string) {
	data.formMemo = utils.Req(c.opts.GetRequest(), "product_memo", "")
	data.formPrice = utils.Req(c.opts.GetRequest(), "product_price", "")
	data.formQuantity = utils.Req(c.opts.GetRequest(), "product_quantity", "")
	data.formStatus = utils.Req(c.opts.GetRequest(), "product_status", "")

	if data.formStatus == "" {
		data.formErrorMessage = "Status is required"
		return data, ""
	}

	if data.formPrice == "" {
		data.formErrorMessage = "Price is required"
		return data, ""
	}

	if data.formQuantity == "" {
		data.formErrorMessage = "Quantity is required"
		return data, ""
	}

	if !govalidator.IsFloat(data.formPrice) {
		data.formErrorMessage = "Price must be numeric"
		return data, ""
	}

	if !govalidator.IsInt(data.formQuantity) {
		data.formErrorMessage = "Quantity must be numeric"
		return data, ""
	}

	price, _ := strconv.ParseFloat(data.formPrice, 64)

	if price < 0 {
		data.formErrorMessage = "Price cannot be negative"
		return data, ""
	}

	quantity, _ := strconv.ParseInt(data.formQuantity, 10, 64)

	if quantity < 0 {
		data.formErrorMessage = "Quantity cannot be negative"
		return data, ""
	}

	data.product.SetMemo(data.formMemo)
	data.product.SetQuantity(data.formQuantity)
	data.product.SetPrice(data.formPrice)
	data.product.SetStatus(data.formStatus)

	err := c.opts.GetStore().ProductUpdate(context.Background(), data.product)

	if err != nil {
		c.opts.GetLogger().Error("At productUpdateController > prepareDataAndValidate", "error", err.Error())
		data.formErrorMessage = "System error. Saving details failed"
		return data, ""
	}

	data.formSuccessMessage = "Product settings saved successfully"

	return data, ""
}

// func ReqArrayOfMaps(r *http.Request, key string, defaultValue []map[string]string) []map[string]string {
// 	all := utils.ReqAll(r)

// 	reqArrayOfMaps := []map[string]string{}

// 	if all == nil {
// 		return reqArrayOfMaps
// 	}

// 	mapIndexMap := map[string]map[string]string{}

// 	// Iterate through all the parameters
// 	for k, v := range all {
// 		if !strings.HasPrefix(k, key+"[") {
// 			continue
// 		}
// 		if !strings.HasSuffix(k, "]") {
// 			continue
// 		}
// 		if !strings.Contains(k, "][") {
// 			continue
// 		}
// 		mapValue := v[0]

// 		str := strings.TrimSuffix(strings.TrimPrefix(k, key+"["), "]")
// 		split := strings.Split(str, "][")
// 		if len(split) != 2 {
// 			// Handle invalid format
// 			continue
// 		}

// 		index, key := split[0], split[1]

// 		if lo.HasKey(mapIndexMap, index) {
// 			mapIndexMap[index][key] = mapValue
// 		} else {
// 			mapIndexMap[index] = map[string]string{
// 				key: mapValue,
// 			}
// 		}
// 	}

// 	for _, v := range mapIndexMap {
// 		if v == nil {
// 			continue
// 		}
// 		reqArrayOfMaps = append(reqArrayOfMaps, v)
// 	}

// 	return reqArrayOfMaps
// }

func (c *productUpdateController) saveProductMedia(data productUpdateControllerData) (d productUpdateControllerData, errorMessage string) {
	media := req.Maps(c.opts.GetRequest(), "product_media", []map[string]string{})

	cfmt.Infoln(media)

	productMedia := lo.Map(media, func(m map[string]string, index int) map[string]string {
		id := strings.TrimSpace(m["id"])
		title := strings.TrimSpace(m["title"])
		url := strings.TrimSpace(m["url"])
		entry := map[string]string{}
		entry["id"] = id
		entry["title"] = title
		entry["url"] = url
		return entry
	})

	data.formMedia = productMedia

	cfmt.Successln(data.formMedia)

	if data.action == "add" {
		data.formMedia = append(data.formMedia, map[string]string{
			"id":    uid.HumanUid(),
			"title": "",
			"url":   "",
		})
		return data, ""
	}
	return data, ""
}

func (c *productUpdateController) saveProductMetadata(data productUpdateControllerData) (d productUpdateControllerData, errorMessage string) {
	metas := req.Maps(c.opts.GetRequest(), "product_meta", []map[string]string{})

	cfmt.Infoln(metas)

	productMetas := map[string]string{}

	lo.ForEach(metas, func(meta map[string]string, index int) {
		metaKey := strings.TrimSpace(meta["key"])
		metaValue := strings.TrimSpace(meta["value"])
		if metaKey == "" {
			return
		}
		productMetas[metaKey] = metaValue
	})

	data.formMetas = productMetas

	cfmt.Infoln(data.formMetas)

	if data.formMetas == nil {
		data.formErrorMessage = "Metadata is required"
		return data, ""
	}

	data.product.SetMetas(data.formMetas)

	err := c.opts.GetStore().ProductUpdate(context.Background(), data.product)

	if err != nil {
		c.opts.GetLogger().Error("At shopstoreadmin > productUpdateController > prepareDataAndValidate", "error", err.Error())
		data.formErrorMessage = "System error. Saving metas failed"
		return data, ""
	}

	data.formSuccessMessage = "Metadata saved successfully"

	return data, ""
}

func (c *productUpdateController) saveProductContent(data productUpdateControllerData) (d productUpdateControllerData, errorMessage string) {
	data.formDescription = utils.Req(c.opts.GetRequest(), "product_description", "")
	data.formShortDescription = utils.Req(c.opts.GetRequest(), "product_short_description", "")
	data.formTitle = utils.Req(c.opts.GetRequest(), "product_title", "")

	if data.formTitle == "" {
		data.formErrorMessage = "Title is required"
		return data, ""
	}

	data.product.SetDescription(data.formDescription)
	data.product.SetShortDescription(data.formShortDescription)
	data.product.SetTitle(data.formTitle)

	err := c.opts.GetStore().ProductUpdate(context.Background(), data.product)

	if err != nil {
		c.opts.GetLogger().Error("At productUpdateController > prepareDataAndValidate", "error", err.Error())
		data.formErrorMessage = "System error. Saving details failed"
		return data, ""
	}

	data.formSuccessMessage = "Product contents saved successfully"

	return data, ""
}

func (c *productUpdateController) prepareDataAndValidate() (data productUpdateControllerData, errorMessage string) {
	data.request = c.opts.GetRequest()
	data.action = utils.Req(c.opts.GetRequest(), "action", "")
	data.productID = utils.Req(c.opts.GetRequest(), "product_id", "")
	data.view = utils.Req(c.opts.GetRequest(), "view", "")

	if data.productID == "" {
		return data, "Product ID is required"
	}

	if data.view == "" {
		data.view = viewContent
	}

	product, err := c.opts.GetStore().ProductFindByID(context.Background(), data.productID)

	if err != nil {
		c.opts.GetLogger().Error("At shopstoreadmin > productUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, "Product not found"
	}

	if product == nil {
		return data, "Product not found"
	}

	data.product = product

	metas, err := product.Metas()

	if err != nil {
		c.opts.GetLogger().Error("At shopstoreadmin > productUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, "Product metas not found"
	}

	media, err := c.opts.GetStore().MediaList(context.Background(), shopstore.NewMediaQuery().
		SetEntityID(product.ID()))

	if err != nil {
		c.opts.GetLogger().Error("At shopstoreadmin > productUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, "Product media not found"
	}

	data.media = media

	data.formMemo = data.product.Memo()
	data.formStatus = data.product.Status()
	data.formTitle = data.product.Title()
	data.formDescription = data.product.Description()
	data.formPrice = data.product.Price()
	data.formQuantity = data.product.Quantity()
	data.formMetas = metas
	data.formMedia = lo.Map(data.media, func(media shopstore.MediaInterface, index int) map[string]string {
		return map[string]string{
			"id":   media.ID(),
			"name": media.Title(),
			"url":  media.URL(),
		}
	})

	if c.opts.GetRequest().Method != http.MethodPost {
		return data, ""
	}

	if data.view == viewContent {
		return c.saveProductContent(data)
	}

	if data.view == viewMedia {
		return c.saveProductMedia(data)
	}

	if data.view == viewMetadata {
		return c.saveProductMetadata(data)
	}

	if data.view == viewSettings {
		return c.saveProductSettings(data)
	}

	return data, "view is required"

	// if data.action == "update-details" {
	// 	return c.saveProductDetails(data)
	// }

	// if data.action == "update-metadata" {
	// 	return c.saveProductMetadata(data)
	// }

	// return data, "action is required"
}

type productUpdateControllerData struct {
	request   *http.Request
	action    string
	productID string
	product   shopstore.ProductInterface
	media     []shopstore.MediaInterface
	view      string

	formErrorMessage     string
	formSuccessMessage   string
	formDescription      string
	formMemo             string
	formMedia            []map[string]string
	formMetas            map[string]string
	formQuantity         string
	formPrice            string
	formShortDescription string
	formStatus           string
	formTitle            string
}
