package shopstoreadmin

import (
	"context"
	"net/http"

	// "project/config"
	// "project/controllers/admin/shop/shared"

	// "project/internal/helpers"
	// "project/internal/layouts"
	// "project/internal/links"
	"strings"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/shopstore"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const ActionModalDiscountFilterShow = "modal_discount_filter_show"

// ===========================================================================
// == CONSTRUCTOR
// ===========================================================================

func discountManager(opts UiOptionsInterface) pageInterface {
	return &discountManagerController{
		opts: opts,
	}
}

// ===========================================================================
// == CONTROLLER
// ===========================================================================

type discountManagerController struct {
	opts UiOptionsInterface
}

// ===========================================================================
// == INTERFACE IMPLEMENTATION
// ===========================================================================

func (c *discountManagerController) ToTag() hb.TagInterface {
	data, errorMessage := c.prepareData(c.opts.GetRequest())

	if errorMessage != "" {
		return hb.Div().HTML(errorMessage)
	}

	if data.action == ActionModalDiscountFilterShow {
		return c.onModalDiscountFilterShow(data)
	}

	c.opts.GetLayout().SetBody(c.page(data).ToHTML())
	c.opts.GetLayout().SetScriptURLs([]string{
		cdn.Htmx_2_0_0(),
		cdn.Sweetalert2_10(),
	})

	c.opts.GetLayout().SetStyles([]string{
		`
.htmx-indicator {
    display: none;
}
.htmx-request .htmx-indicator {
    display: inline-block;
}
		`,
	})

	return hb.Raw(c.opts.GetLayout().Render(c.opts.GetResponseWriter(), c.opts.GetRequest()))
}

func (controller *discountManagerController) ToHTML() string {
	return controller.ToTag().ToHTML()
}

// ===========================================================================
// == METHODS
// ===========================================================================

func (controller *discountManagerController) onModalDiscountFilterShow(data discountManagerControllerData) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.Heading5().
		Text("Filters").
		Style(`margin:0px;padding:0px;`)

	buttonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonOk := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Apply").
		Class("btn btn-primary float-end").
		OnClick(`FormFilters.submit();` + modalCloseScript)

	fieldStatus := form.NewField(form.FieldOptions{
		Label: "Status",
		Name:  "filter_status",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Help:  `The status of the discount.`,
		Value: data.formStatus,
		Options: []form.FieldOption{
			{
				Value: "",
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

	fieldTitle := form.NewField(form.FieldOptions{
		Label: "Title",
		Name:  "filter_title",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formTitle,
		Help:  `Filter by title.`,
	})

	fieldCreatedFrom := form.NewField(form.FieldOptions{
		Label: "Created From",
		Name:  "filter_created_from",
		Type:  form.FORM_FIELD_TYPE_DATE,
		Value: data.formCreatedFrom,
		Help:  `Filter by creation date.`,
	})

	fieldCreatedTo := form.NewField(form.FieldOptions{
		Label: "Created To",
		Name:  "filter_created_to",
		Type:  form.FORM_FIELD_TYPE_DATE,
		Value: data.formCreatedTo,
		Help:  `Filter by creation date.`,
	})

	fieldUpdatedFrom := form.NewField(form.FieldOptions{
		Label: "Updated From",
		Name:  "filter_updated_from",
		Type:  form.FORM_FIELD_TYPE_DATE,
		Value: data.formUpdatedFrom,
		Help:  `Filter by update date.`,
	})

	fieldUpdatedTo := form.NewField(form.FieldOptions{
		Label: "Updated To",
		Name:  "filter_updated_to",
		Type:  form.FORM_FIELD_TYPE_DATE,
		Value: data.formUpdatedTo,
		Help:  `Filter by update date.`,
	})

	fieldDiscountID := form.NewField(form.FieldOptions{
		Label: "Discount ID",
		Name:  "filter_discount_id",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formDiscountID,
		Help:  `Find discount by reference number (ID).`,
	})

	fieldController := form.NewField(form.FieldOptions{
		Label:     "Controller",
		Name:      "controller",
		Type:      form.FORM_FIELD_TYPE_STRING,
		Value:     pathDiscounts,
		Help:      `Required to display the correct page.`,
		Invisible: true,
	})

	filterForm := form.NewForm(form.FormOptions{
		ID:     "FormFilters",
		Method: http.MethodGet,
		Fields: []form.FieldInterface{
			fieldStatus,
			fieldTitle,
			fieldCreatedFrom,
			fieldCreatedTo,
			fieldUpdatedFrom,
			fieldUpdatedTo,
			fieldDiscountID,
			fieldController,
		},
	}).Build()

	modal := bs.Modal().
		ID("ModalMessage").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						title,
						buttonModalClose,
					}),

					bs.ModalBody().
						Child(filterForm),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonOk),
				}),
			}),
		})

	backdrop := hb.Div().
		ID("ModalBackdrop").
		Class("modal-backdrop fade show").
		Style("display:block;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})

}

func (controller *discountManagerController) page(data discountManagerControllerData) hb.TagInterface {
	discountURL := url(controller.opts.GetRequest(), pathDiscounts, nil)
	discountCreateURL := url(controller.opts.GetRequest(), pathDiscountCreate, nil)

	breadcrumbs := breadcrumbs(controller.opts.GetRequest(), []breadcrumb{
		{
			Name: "Discount Manager",
			URL:  discountURL,
		},
	})

	buttonDiscountNew := hb.Button().
		Class("btn btn-primary float-end").
		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("New Discount").
		HxGet(discountCreateURL).
		HxTarget("body").
		HxSwap("beforeend")

	title := hb.Heading1().
		HTML("Shop. Discount Manager").
		Child(buttonDiscountNew)

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(header(controller.opts)).
		Child(hb.HR()).
		Child(title).
		Child(controller.tableDiscounts(data))
}

func (controller *discountManagerController) tableDiscounts(data discountManagerControllerData) hb.TagInterface {
	type columnStruct struct {
		Children    []columnStruct
		Label       string
		Sortable    bool
		SortableKey string
	}
	columns := []columnStruct{
		{
			Children: []columnStruct{
				{
					Label:       "Title",
					Sortable:    true,
					SortableKey: "title",
				},
				{
					Label:       "Reference",
					Sortable:    true,
					SortableKey: "id",
				},
			},
		},
		{
			Label:       "Status",
			Sortable:    true,
			SortableKey: "status",
		},
		{
			Label:       "Created",
			Sortable:    true,
			SortableKey: "created_at",
		},
		{
			Label:       "Modified",
			Sortable:    true,
			SortableKey: "updated_at",
		},
		{
			Label:    "Actions",
			Sortable: false,
		},
	}
	table := hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().
				Child(hb.TR().Children(lo.Map(columns, func(column columnStruct, _ int) hb.TagInterface {
					children := []columnStruct{}
					if len(column.Children) == 0 {
						children = append(children, column)
					} else {
						children = append(children, column.Children...)
					}

					links := lo.Map(children, func(column columnStruct, _ int) string {
						return hb.Span().
							ChildIf(column.Sortable, controller.sortableColumnLabel(data, column.Label, column.SortableKey)).
							StyleIf(column.Sortable, "cursor: pointer;").
							TextIf(!column.Sortable, column.Label).
							ToHTML()
					})

					return hb.TH().
						HTML(strings.Join(links, ", ")).
						Style("width: 1px;")
				}))),
			hb.Tbody().Children(lo.Map(data.discountList, func(discount shopstore.DiscountInterface, _ int) hb.TagInterface {
				discountDeleteURL := url(controller.opts.GetRequest(), pathDiscountDelete, map[string]string{"discount_id": discount.ID()})
				discountUpdateURL := url(controller.opts.GetRequest(), pathDiscountUpdate, map[string]string{"discount_id": discount.ID()})

				discountLink := hb.Hyperlink().
					Text(discount.Title()).
					Href(discountUpdateURL)

				status := hb.Span().
					Style(`font-weight: bold;`).
					StyleIf(discount.Status() == shopstore.PRODUCT_STATUS_ACTIVE, `color:green;`).
					StyleIf(discount.Status() == shopstore.PRODUCT_STATUS_DISABLED, `color:silver;`).
					StyleIf(discount.Status() == shopstore.PRODUCT_STATUS_DRAFT, `color:blue;`).
					HTML(discount.Status())

				buttonEdit := hb.Hyperlink().
					Class("btn btn-primary me-2").
					Child(hb.I().Class("bi bi-pencil-square")).
					Title("Edit").
					Href(discountUpdateURL)

				buttonDelete := hb.Hyperlink().
					Class("btn btn-danger").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete").
					HxGet(discountDeleteURL).
					HxTarget("body").
					HxSwap("beforeend")

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().
						Child(hb.Div().Child(discountLink)).
						Child(hb.Div().
							Style("font-size: 11px;").
							HTML("Ref: ").
							HTML(discount.ID())),
					hb.TD().
						Child(status),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(discount.CreatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(discount.UpdatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(buttonEdit).
						Child(buttonDelete),
				})
			})),
		})

	// cfmt.Successln("Table: ", table)

	return hb.Wrap().Children([]hb.TagInterface{
		controller.tableFilter(data),
		table,
		controller.tablePagination(data, int(data.discountCount), data.pageInt, data.perPage),
	})
}

func (controller *discountManagerController) sortableColumnLabel(data discountManagerControllerData, columnLabel string, columnSortKey string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnSortKey)

	changeDirection := sb.ASC

	if isSelected {
		changeDirection = lo.If(data.sortOrder == sb.ASC, sb.DESC).Else(sb.ASC)
	}

	discountsURL := url(controller.opts.GetRequest(), pathDiscounts, map[string]string{
		"page":        "0",
		"by":          columnSortKey,
		"sort":        changeDirection,
		"date_from":   data.formCreatedFrom,
		"date_to":     data.formCreatedTo,
		"status":      data.formStatus,
		"discount_id": data.formDiscountID,
	})

	return hb.Hyperlink().
		HTML(columnLabel).
		Child(controller.sortingIndicator(columnSortKey, data.sortBy, changeDirection)).
		Href(discountsURL)
}

func (controller *discountManagerController) sortingIndicator(columnSortKey string, sortByColumnKey string, sortOrder string) hb.TagInterface {
	isSelected := strings.EqualFold(sortByColumnKey, columnSortKey)

	direction := lo.If(isSelected && sortOrder == "asc", "up").
		ElseIf(isSelected && sortOrder == "desc", "down").
		Else("none")

	sortingIndicator := hb.Span().
		Class("sorting").
		HTMLIf(direction == "up", "&#8595;").
		HTMLIf(direction == "down", "&#8593;").
		HTMLIf(direction != "down" && direction != "up", "")

	return sortingIndicator
}

func (controller *discountManagerController) tableFilter(data discountManagerControllerData) hb.TagInterface {
	discountsURL := url(controller.opts.GetRequest(), pathDiscounts, map[string]string{
		"action":              ActionModalDiscountFilterShow,
		"filter_title":        data.formTitle,
		"filter_status":       data.formStatus,
		"filter_discount_id":  data.formDiscountID,
		"filter_created_from": data.formCreatedFrom,
		"filter_created_to":   data.formCreatedTo,
		"filter_updated_from": data.formUpdatedFrom,
		"filter_updated_to":   data.formUpdatedTo,
	})

	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(discountsURL).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing discounts").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formTitle != "" {
		description = append(description, hb.Span().Text("and email: "+data.formTitle).ToHTML())
	}

	if data.formDiscountID != "" {
		description = append(description, hb.Span().Text("and ID: "+data.formDiscountID).ToHTML())
	}

	if data.formCreatedFrom != "" && data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created between: "+data.formCreatedFrom+" and "+data.formCreatedTo).ToHTML())
	} else if data.formCreatedFrom != "" {
		description = append(description, hb.Span().Text("and created after: "+data.formCreatedFrom).ToHTML())
	} else if data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created before: "+data.formCreatedTo).ToHTML())
	}

	return hb.Div().
		Class("card bg-light mb-3").
		Style("").
		Children([]hb.TagInterface{
			hb.Div().Class("card-body").
				Child(buttonFilter).
				Child(hb.Span().
					HTML(strings.Join(description, " "))),
		})
}

func (controller *discountManagerController) tablePagination(data discountManagerControllerData, count int, page int, perPage int) hb.TagInterface {
	discountsURL := url(controller.opts.GetRequest(), pathDiscounts, map[string]string{
		"filter_status":       data.formStatus,
		"filter_created_from": data.formCreatedFrom,
		"filter_created_to":   data.formCreatedTo,
		"filter_discount_id":  data.formDiscountID,
		"filter_updated_from": data.formUpdatedFrom,
		"filter_updated_to":   data.formUpdatedTo,
		"by":                  data.sortBy,
		"order":               data.sortOrder,
	})

	discountsURL = lo.Ternary(strings.Contains(discountsURL, "?"), discountsURL+"&page=", discountsURL+"?page=") // page must be last

	pagination := bs.Pagination(bs.PaginationOptions{
		NumberItems:       count,
		CurrentPageNumber: page,
		PagesToShow:       5,
		PerPage:           perPage,
		URL:               discountsURL,
	})

	return hb.Div().
		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
		HTML(pagination)
}

func (controller *discountManagerController) prepareData(r *http.Request) (data discountManagerControllerData, errorMessage string) {
	var err error
	data.request = r
	data.action = utils.Req(r, "action", "")
	data.page = utils.Req(r, "page", "0")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(utils.Req(r, "per_page", "10"))
	data.sortOrder = utils.Req(r, "sort", sb.DESC)
	data.sortBy = utils.Req(r, "by", shopstore.COLUMN_CREATED_AT)
	data.formDiscountID = utils.Req(r, "filter_discount_id", "")
	data.formTitle = utils.Req(r, "filter_title", "")
	data.formStatus = utils.Req(r, "filter_status", "")
	data.formCreatedFrom = utils.Req(r, "filter_created_from", "")
	data.formCreatedTo = utils.Req(r, "filter_created_to", "")
	data.formUpdatedFrom = utils.Req(r, "filter_updated_from", "")
	data.formUpdatedTo = utils.Req(r, "filter_updated_to", "")

	discountList, discountCount, err := controller.fetchDiscountList(data)

	if err != nil {
		controller.opts.GetLogger().Error("At shopstore > discountManagerController > prepareData", "error", err.Error())
		return data, "error retrieving discounts"
	}

	data.discountList = discountList
	data.discountCount = discountCount

	return data, ""
}

func (controller *discountManagerController) fetchDiscountList(data discountManagerControllerData) ([]shopstore.DiscountInterface, int64, error) {
	discountIDs := []string{}

	if !lo.Contains([]string{sb.DESC, sb.ASC}, data.sortOrder) {
		data.sortOrder = sb.ASC
	}

	if !lo.Contains([]string{
		shopstore.COLUMN_CREATED_AT,
		shopstore.COLUMN_TITLE,
		shopstore.COLUMN_ID,
		shopstore.COLUMN_STATUS,
		shopstore.COLUMN_UPDATED_AT,
	}, data.sortBy) {
		data.sortBy = shopstore.COLUMN_CREATED_AT
	}

	query := shopstore.NewDiscountQuery().
		SetOffset(data.pageInt * data.perPage).
		SetLimit(data.perPage).
		SetOrderBy(data.sortBy).
		SetSortDirection(data.sortOrder)

	if len(discountIDs) > 0 {
		query.SetIDIn(discountIDs)
	}

	if data.formStatus != "" {
		query.SetStatus(data.formStatus)
	}

	if data.formCreatedFrom != "" {
		query.SetCreatedAtGte(data.formCreatedFrom + " 00:00:00")
	}

	if data.formCreatedTo != "" {
		query.SetCreatedAtLte(data.formCreatedTo + " 23:59:59")
	}

	discountList, err := controller.opts.GetStore().DiscountList(context.Background(), query)

	if err != nil {
		controller.opts.GetLogger().Error("At shopstore > discountManagerController > prepareData", "error", err.Error())
		return []shopstore.DiscountInterface{}, 0, err
	}

	discountCount, err := controller.opts.GetStore().DiscountCount(context.Background(), query)

	if err != nil {
		controller.opts.GetLogger().Error("At shopstore > discountManagerController > prepareData", "error", err.Error())
		return []shopstore.DiscountInterface{}, 0, err
	}

	return discountList, discountCount, nil
}

type discountManagerControllerData struct {
	request         *http.Request
	action          string
	page            string
	pageInt         int
	perPage         int
	sortOrder       string
	sortBy          string
	formStatus      string
	formTitle       string
	formCreatedFrom string
	formCreatedTo   string
	formUpdatedFrom string
	formUpdatedTo   string
	formDiscountID  string
	discountList    []shopstore.DiscountInterface
	discountCount   int64
}
