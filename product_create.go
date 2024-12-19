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

func productCreate(opts UiOptionsInterface) pageInterface {
	return &productCreateController{
		opts: opts,
	}
}

type productCreateController struct {
	opts UiOptionsInterface
}

// ===========================================================================
// == INTERFACE IMPLEMENTATION
// ===========================================================================

func (controller productCreateController) ToTag() hb.TagInterface {
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

func (controller productCreateController) ToHTML() string {
	return controller.ToTag().ToHTML()
}

// ===========================================================================
// == METHODS
// ===========================================================================

func (controller *productCreateController) modal(data productCreateControllerData) hb.TagInterface {
	submitUrl := url(controller.opts.GetRequest(), pathProductCreate, nil)

	formGroupTitle := bs.FormGroup().
		Class("mb-3").
		Child(bs.FormLabel("Title")).
		Child(bs.FormInput().Name("product_title").Value(data.formTitle))

	modalID := "ModalproductCreate"
	modalBackdropClass := "ModalBackdrop"

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("New product Create").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalproductCreate').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

	buttonSend := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Create & Edit").
		Class("btn btn-primary float-end").
		HxInclude("#" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalproductCreate").
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

func (c *productCreateController) prepareDataAndValidate() (data productCreateControllerData, errorMessage string) {
	data.formTitle = strings.TrimSpace(utils.Req(c.opts.GetRequest(), "product_title", ""))

	if c.opts.GetRequest().Method != http.MethodPost {
		return data, ""
	}

	if data.formTitle == "" {
		return data, "product title is required"
	}

	product := shopstore.NewProduct()
	product.SetTitle(data.formTitle)

	err := c.opts.GetStore().ProductCreate(c.opts.GetRequest().Context(), product)

	if err != nil {
		c.opts.GetLogger().Error("Error. At shopstoreadmin > productCreateController > prepareDataAndValidate", "err", err.Error())
		return data, "Creating product failed. Please contact an administrator."
	}

	data.successMessage = "product created successfully."

	return data, ""

}

// ===========================================================================
// == DATA
// ===========================================================================

type productCreateControllerData struct {
	formTitle      string
	successMessage string
	//errorMessage   string
}
