package shopstoreadmin

import (
	"context"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/vld"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/crud/v2"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore"
	"github.com/gouniverse/utils"
	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// ===========================================================================
// == CONSTRUCTOR
// ===========================================================================

func discountUpdate(opts UiOptionsInterface) pageInterface {
	return &discountUpdateController{
		opts: opts,
	}
}

// ===========================================================================
// == CONTROLLER
// ===========================================================================

type discountUpdateController struct {
	opts UiOptionsInterface
}

// ===========================================================================
// == INTERFACE IMPLEMENTATION
// ===========================================================================

func (c *discountUpdateController) ToTag() hb.TagInterface {
	data, errorMessage := c.prepareDataAndValidate()

	if errorMessage != "" {
		return hb.Div().Class("alert alert-danger").Child(hb.Text(errorMessage))
	}

	if c.opts.GetRequest().Method == http.MethodPost {
		return c.form(data)
	}

	c.opts.GetLayout().SetTitle("Edit Discount | Shop")
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

func (controller *discountUpdateController) ToHTML() string {
	return controller.ToTag().ToHTML()
}

// ===========================================================================
// == METHODS
// ===========================================================================

func (c *discountUpdateController) page(data discountUpdateControllerData) hb.TagInterface {
	discountManegerURL := url(c.opts.GetRequest(), pathDiscounts, map[string]string{})

	discountUpdateURL := url(c.opts.GetRequest(), pathDiscountUpdate, map[string]string{
		"discountID": data.discountID,
	})

	breadcrumbs := breadcrumbs(c.opts.GetRequest(), []breadcrumb{
		{
			Name: "Discount Manager",
			URL:  discountManegerURL,
		},
		{
			Name: "Edit Discount",
			URL:  discountUpdateURL,
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(discountManegerURL)

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		Child(hb.Div().ID("ButtonSaveIndicator").Class("spinner-border spinner-border-sm ms-2 htmx-indicator")).
		HxIndicator("#ButtonSaveIndicator").
		HxInclude("#FormDiscountUpdate").
		HxPost(discountUpdateURL).
		HxTarget("#FormDiscountUpdate")

	title := hb.Heading1().
		Class("mb-3").
		Text("Edit Discount: ").
		Text(data.discount.Title()).
		Child(buttonCancel)

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTMLIf(data.view == viewContent, "Discount Contents").
					HTMLIf(data.view == viewMetadata, "Discount Metadata").
					HTMLIf(data.view == viewSettings, "Discount Settings").
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
		// Child(c.tabs(data)).
		Child(card).
		Child(hb.BR()).
		Child(hb.BR()).
		Child(hb.BR()).
		Child(hb.BR())
}

// func (controller *discountUpdateController) tabs(data discountUpdateControllerData) hb.TagInterface {
// 	viewContentURL := url(controller.opts.GetRequest(), pathDiscountUpdate, map[string]string{
// 		"discount_id": data.discountID,
// 		"view":        viewContent,
// 	})

// 	viewMedia := url(controller.opts.GetRequest(), pathDiscountUpdate, map[string]string{
// 		"discount_id": data.discountID,
// 		"view":        viewMedia,
// 	})

// 	viewMetadataURL := url(controller.opts.GetRequest(), pathDiscountUpdate, map[string]string{
// 		"discount_id": data.discountID,
// 		"view":        viewMetadata,
// 	})

// 	viewSettingsURL := url(controller.opts.GetRequest(), pathDiscountUpdate, map[string]string{
// 		"discount_id": data.discountID,
// 		"view":        viewSettings,
// 	})

// 	tabs := bs.NavTabs().
// 		Class("mb-3").
// 		Child(bs.NavItem().
// 			Child(bs.NavLink().
// 				ClassIf(data.view == viewContent, "active").
// 				Href(viewContentURL).
// 				HTML("Content"))).
// 		Child(bs.NavItem().
// 			Child(bs.NavLink().
// 				ClassIf(data.view == viewMedia, "active").
// 				Href(viewMedia).
// 				HTML("Media"))).
// 		Child(bs.NavItem().
// 			Child(bs.NavLink().
// 				ClassIf(data.view == viewMetadata, "active").
// 				Href(viewMetadataURL).
// 				HTML("Metas"))).
// 		Child(bs.NavItem().
// 			Child(bs.NavLink().
// 				ClassIf(data.view == viewSettings, "active").
// 				Href(viewSettingsURL).
// 				HTML("Settings")))

// 	return tabs
// }

func (c *discountUpdateController) form(data discountUpdateControllerData) hb.TagInterface {
	formDiscountUpdate := form.NewForm(form.FormOptions{
		ID: "FormDiscountUpdate",
	})

	formDiscountUpdate.SetFields(c.formSettingsFields(data))

	if data.formErrorMessage != "" {
		formDiscountUpdate.AddField(form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
		}))
	}

	if data.formSuccessMessage != "" {
		formDiscountUpdate.AddField(form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}).ToHTML(),
		}))
	}

	discountID := form.NewField(form.FieldOptions{
		Label:     "Discount ID",
		Name:      "discount_id",
		Type:      form.FORM_FIELD_TYPE_STRING,
		Value:     data.discountID,
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

	action := form.NewField(form.FieldOptions{
		Label:     "Action",
		Name:      "action",
		Type:      form.FORM_FIELD_TYPE_HIDDEN,
		Value:     actionSave,
		Readonly:  true,
		Invisible: true,
	})

	formDiscountUpdate.AddField(discountID)
	formDiscountUpdate.AddField(view)
	formDiscountUpdate.AddField(action)

	return formDiscountUpdate.Build()

}

func (c *discountUpdateController) formSettingsFields(data discountUpdateControllerData) []form.FieldInterface {
	title := form.NewField(form.FieldOptions{
		Label: "Title",
		Name:  "discount_title",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formTitle,
		Help:  `The title of the discount.`,
	})

	description := form.NewField(form.FieldOptions{
		Label: "Description",
		Name:  "discount_description",
		Type:  form.FORM_FIELD_TYPE_HTMLAREA,
		Value: data.formDescription,
		Help:  `The description of the discount.`,
	})

	status := form.NewField(form.FieldOptions{
		Label: "Status",
		Name:  "discount_status",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formStatus,
		Help:  `The status of the discount.`,
		Options: []form.FieldOption{
			{
				Value: "- not selected -",
				Key:   "",
			},
			{
				Value: "Active",
				Key:   shopstore.DISCOUNT_STATUS_ACTIVE,
			},
			{
				Value: "Inactive",
				Key:   shopstore.DISCOUNT_STATUS_INACTIVE,
			},
			{
				Value: "Draft",
				Key:   shopstore.DISCOUNT_STATUS_DRAFT,
			},
		},
	})

	discountType := form.NewField(form.FieldOptions{
		Label: "Type of Discount",
		Name:  "discount_type",
		Help:  `The type of the discount.`,
		Type:  crud.FORM_FIELD_TYPE_SELECT,
		Value: data.formType,
		Options: []form.FieldOption{
			{
				Key:   "",
				Value: "",
			},
			{
				Key:   shopstore.DISCOUNT_TYPE_AMOUNT,
				Value: shopstore.DISCOUNT_TYPE_AMOUNT,
			},
			{
				Key:   shopstore.DISCOUNT_TYPE_PERCENT,
				Value: shopstore.DISCOUNT_TYPE_PERCENT,
			},
		},
	})

	amount := form.NewField(form.FieldOptions{
		Label: "Amount / Percentage",
		Name:  "discount_amount",
		Type:  form.FORM_FIELD_TYPE_NUMBER,
		Value: data.formAmount,
		Help:  `The amount / percentage of the discount. Depends on the discount type.`,
	})

	discountCode := form.NewField(form.FieldOptions{
		Label: "Discount Code",
		Name:  "discount_code",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formCode,
		Help:  `The code of the discount. It must be unique.`,
	})

	startsAt := form.NewField(form.FieldOptions{
		Label: "Time Starts (UTC)",
		Name:  "discount_starts_at",
		Type:  form.FORM_FIELD_TYPE_DATETIME,
		Value: lo.Ternary(data.formStartsAt == "", "", carbon.Parse(data.formStartsAt, carbon.UTC).ToDateTimeString()),
	})

	endsAt := form.NewField(form.FieldOptions{
		Label: "Time Ends (UTC)",
		Name:  "discount_ends_at",
		Type:  form.FORM_FIELD_TYPE_DATETIME,
		Value: lo.Ternary(data.formEndsAt == "", "", carbon.Parse(data.formEndsAt, carbon.UTC).ToDateTimeString()),
	})

	memo := form.NewField(form.FieldOptions{
		Label: "Admin Notes",
		Name:  "discount_memo",
		Type:  form.FORM_FIELD_TYPE_TEXTAREA,
		Value: data.formMemo,
		Help:  "Admin notes for this discount. These notes will not be visible to the public.",
	})

	return []form.FieldInterface{
		status,
		title,
		description,
		discountCode,
		discountType,
		amount,
		startsAt,
		endsAt,
		memo,
	}
}

func (c *discountUpdateController) saveDiscount(data discountUpdateControllerData) (d discountUpdateControllerData, errorMessage string) {
	data.formAmount = utils.Req(c.opts.GetRequest(), "discount_amount", "")
	data.formCode = utils.Req(c.opts.GetRequest(), "discount_code", "")
	data.formDescription = utils.Req(c.opts.GetRequest(), "discount_description", "")
	data.formEndsAt = utils.Req(c.opts.GetRequest(), "discount_ends_at", "")
	data.formMemo = utils.Req(c.opts.GetRequest(), "discount_memo", "")
	data.formStatus = utils.Req(c.opts.GetRequest(), "discount_status", "")
	data.formStartsAt = utils.Req(c.opts.GetRequest(), "discount_starts_at", "")
	data.formTitle = utils.Req(c.opts.GetRequest(), "discount_title", "")
	data.formType = utils.Req(c.opts.GetRequest(), "discount_type", "")

	if data.formStatus == "" {
		data.formErrorMessage = "Status is required"
		return data, ""
	}

	if data.formTitle == "" {
		data.formErrorMessage = "Title is required"
		return data, ""
	}

	if data.formCode == "" {
		data.formErrorMessage = "Code is required"
		return data, ""
	}

	if data.formType == "" {
		data.formErrorMessage = "Type is required"
		return data, ""
	}

	if data.formAmount == "" {
		data.formErrorMessage = "Amount is required"
		return data, ""
	}

	if !govalidator.IsFloat(data.formAmount) {
		data.formErrorMessage = "Amount must be numeric"
		return data, ""
	}

	price, _ := strconv.ParseFloat(data.formAmount, 64)

	if price < 0 {
		data.formErrorMessage = "Amount cannot be negative"
		return data, ""
	}

	if data.formStartsAt == "" {
		data.formErrorMessage = "Starts at is required"
		return data, ""
	}

	if data.formEndsAt == "" {
		data.formErrorMessage = "Ends at is required"
		return data, ""
	}

	if len(data.formStartsAt) == 16 {
		data.formStartsAt = data.formStartsAt + ":00"
	}

	if len(data.formEndsAt) == 16 {
		data.formEndsAt = data.formEndsAt + ":59"
	}

	cfmt.Infoln(data.formStartsAt, data.formEndsAt)

	if !vld.IsDateTime(data.formStartsAt) {
		data.formErrorMessage = "Starts at must be a valid date"
		return data, ""
	}

	if !vld.IsDateTime(data.formEndsAt) {
		data.formErrorMessage = "Ends at must be a valid date"
		return data, ""
	}

	if data.formStartsAt > data.formEndsAt {
		data.formErrorMessage = "Starts at must be before ends at"
		return data, ""
	}

	data.discount.SetCode(data.formCode)
	data.discount.SetDescription(data.formDescription)
	data.discount.SetEndsAt(data.formEndsAt)
	data.discount.SetMemo(data.formMemo)
	data.discount.SetStartsAt(data.formStartsAt)
	data.discount.SetStatus(data.formStatus)
	data.discount.SetType(data.formType)
	data.discount.SetTitle(data.formTitle)

	err := c.opts.GetStore().DiscountUpdate(context.Background(), data.discount)

	if err != nil {
		c.opts.GetLogger().Error("At discountUpdateController > prepareDataAndValidate", "error", err.Error())
		data.formErrorMessage = "System error. Saving details failed"
		return data, ""
	}

	data.formSuccessMessage = "Discount settings saved successfully"

	return data, ""
}

// func (c *discountUpdateController) saveDiscountMedia(data discountUpdateControllerData) (d discountUpdateControllerData, errorMessage string) {
// 	media := req.Maps(c.opts.GetRequest(), "discount_media", []map[string]string{})

// 	cfmt.Infoln(media)

// 	discountMedia := lo.Map(media, func(m map[string]string, index int) map[string]string {
// 		id := strings.TrimSpace(m["id"])
// 		title := strings.TrimSpace(m["title"])
// 		url := strings.TrimSpace(m["url"])
// 		entry := map[string]string{}
// 		entry["id"] = id
// 		entry["title"] = title
// 		entry["url"] = url
// 		return entry
// 	})

// 	data.formMedia = discountMedia

// 	cfmt.Successln(data.formMedia)

// 	if data.action == "add" {
// 		data.formMedia = append(data.formMedia, map[string]string{
// 			"id":    uid.HumanUid(),
// 			"title": "",
// 			"url":   "",
// 		})
// 		return data, ""
// 	}
// 	return data, ""
// }

func (c *discountUpdateController) prepareDataAndValidate() (data discountUpdateControllerData, errorMessage string) {
	data.request = c.opts.GetRequest()
	data.action = utils.Req(c.opts.GetRequest(), "action", "")
	data.discountID = utils.Req(c.opts.GetRequest(), "discount_id", "")

	if data.discountID == "" {
		return data, "Discount ID is required"
	}

	discount, err := c.opts.GetStore().DiscountFindByID(context.Background(), data.discountID)

	if err != nil {
		c.opts.GetLogger().Error("At shopstoreadmin > discountUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, "Discount not found"
	}

	if discount == nil {
		return data, "Discount not found"
	}

	data.discount = discount

	data.formStatus = data.discount.Status()
	data.formTitle = data.discount.Title()
	data.formDescription = data.discount.Description()
	data.formCode = data.discount.Code()
	data.formType = data.discount.Type()
	data.formAmount = cast.ToString(data.discount.Amount())
	data.formStartsAt = data.discount.StartsAt()
	data.formEndsAt = data.discount.EndsAt()
	data.formMemo = data.discount.Memo()
	data.formMetas, err = data.discount.Metas()

	if err != nil {
		c.opts.GetLogger().Error("At shopstoreadmin > discountUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, "Discount not found"
	}

	if data.action == actionSave {
		return c.saveDiscount(data)
	}

	return data, ""
}

type discountUpdateControllerData struct {
	request    *http.Request
	action     string
	discountID string
	discount   shopstore.DiscountInterface
	media      []shopstore.MediaInterface
	view       string

	formErrorMessage   string
	formSuccessMessage string
	formAmount         string
	formCode           string
	formDescription    string
	formEndsAt         string
	formMemo           string
	formMetas          map[string]string
	formStartsAt       string
	formStatus         string
	formTitle          string
	formType           string
}
