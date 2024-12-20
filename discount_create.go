package shopstoreadmin

import (
	"net/http"
	"strings"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore"
	"github.com/gouniverse/utils"
)

// ===========================================================================
// == CONSTRUCTOR
// ===========================================================================

func discountCreate(opts UiOptionsInterface) pageInterface {
	return &discountCreateController{
		opts: opts,
	}
}

// ===========================================================================
// == CONTROLLER
// ===========================================================================

type discountCreateController struct {
	opts UiOptionsInterface
}

// ===========================================================================
// == INTERFACE IMPLEMENTATION
// ===========================================================================

func (controller discountCreateController) ToTag() hb.TagInterface {
	data, errorMessage := controller.prepareDataAndValidate()

	if errorMessage != "" {
		return hb.Swal(hb.SwalOptions{
			Icon: "error",
			Text: errorMessage,
		})
	}

	if data.successMessage != "" {
		return hb.Wrap().
			Child(hb.Swal(hb.SwalOptions{
				Icon: "success",
				Text: data.successMessage,
			})).
			Child(hb.Script("setTimeout(() => {window.location.href = window.location.href}, 2000)"))
	}

	return controller.modal(data)
}

func (controller discountCreateController) ToHTML() string {
	return controller.ToTag().ToHTML()
}

// ===========================================================================
// == METHODS
// ===========================================================================

func (controller *discountCreateController) modal(data discountCreateControllerData) hb.TagInterface {
	submitUrl := url(controller.opts.GetRequest(), pathDiscountCreate, nil)

	formGroupTitle := bs.FormGroup().
		Class("mb-3").
		Child(bs.FormLabel("Title")).
		Child(bs.FormInput().Name("discount_title").Value(data.formTitle))

	modalID := "ModaldiscountCreate"
	modalBackdropClass := "ModalBackdrop"

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("New discount Create").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModaldiscountCreate').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

	buttonSend := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Create & Edit").
		Class("btn btn-primary float-end").
		Child(hb.Div().ID("ButtonSaveIndicator").Class("spinner-border spinner-border-sm ms-2 htmx-indicator")).
		HxIndicator("#ButtonSaveIndicator").
		HxInclude("#" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModaldiscountCreate").
		HxTarget("body").
		HxSwap("beforeend")

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Close").
		Class("btn btn-secondary float-start").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Child(hb.Script(jsCloseFn)).
		Child(bs.ModalDialog().
			Child(bs.ModalContent().
				Child(
					bs.ModalHeader().
						Child(modalHeading).
						Child(modalClose)).
				Child(
					bs.ModalBody().
						Child(formGroupTitle)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(buttonCancel).
					Child(buttonSend)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})
}

func (c *discountCreateController) prepareDataAndValidate() (data discountCreateControllerData, errorMessage string) {
	data.formTitle = strings.TrimSpace(utils.Req(c.opts.GetRequest(), "discount_title", ""))

	if c.opts.GetRequest().Method != http.MethodPost {
		return data, ""
	}

	if data.formTitle == "" {
		return data, "discount title is required"
	}

	discount := shopstore.NewDiscount()
	discount.SetTitle(data.formTitle)

	err := c.opts.GetStore().DiscountCreate(c.opts.GetRequest().Context(), discount)

	if err != nil {
		c.opts.GetLogger().Error("Error. At shopstoreadmin > discountCreateController > prepareDataAndValidate", "err", err.Error())
		return data, "Creating discount failed. Please contact an administrator."
	}

	data.successMessage = "discount created successfully."

	return data, ""

}

// ===========================================================================
// == DATA
// ===========================================================================

type discountCreateControllerData struct {
	formTitle      string
	successMessage string
	//errorMessage   string
}
