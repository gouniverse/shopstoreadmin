package shopstoreadmin

import (
	"github.com/gouniverse/bs"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore"
	"github.com/gouniverse/utils"
)

// ===========================================================================
// == CONSTRUCTOR
// ===========================================================================

func discountDelete(opts UiOptionsInterface) pageInterface {
	return &discountDeleteController{
		opts: opts,
	}
}

// ===========================================================================
// == CONTROLLER
// ===========================================================================

type discountDeleteController struct {
	opts UiOptionsInterface
}

// ===========================================================================
// == INTERFACE IMPLEMENTATION
// ===========================================================================

func (controller discountDeleteController) ToTag() hb.TagInterface {
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

func (controller discountDeleteController) ToHTML() string {
	return controller.ToTag().ToHTML()
}

// ===========================================================================
// == METHODS
// ===========================================================================

func (controller *discountDeleteController) modal(data discountDeleteControllerData) hb.TagInterface {
	submitUrl := url(controller.opts.GetRequest(), pathDiscountDelete, map[string]string{
		"discount_id": data.discountID,
	})

	modalID := "ModalDiscountDelete"
	modalBackdropClass := "ModalBackdrop"

	formGroupDiscountId := hb.Input().
		Type(hb.TYPE_HIDDEN).
		Name("discount_id").
		Value(data.discountID)

	buttonDelete := hb.Button().
		Type("button").
		Child(hb.I().Class("bi bi-trash me-2")).
		HTML("Delete").
		Class("btn btn-danger float-end").
		Child(hb.Div().ID("ButtonDeleteIndicator").Class("spinner-border spinner-border-sm ms-2 htmx-indicator")).
		HxIndicator("#ButtonDeleteIndicator").
		HxInclude("#Modal" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalDiscountDelete").
		HxTarget("body").
		HxSwap("beforeend")

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Delete Discount").Style(`margin:0px;`)

	modalClose := hb.Button().
		Type("button").
		Child(hb.I().Class("bi bi-chevron-left")).
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalDiscountDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

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
						Child(hb.Paragraph().Text("Are you sure you want to delete this discount?").Style(`margin-bottom:20px;color:red;`)).
						Child(hb.Paragraph().Text("This action cannot be undone.")).
						Child(formGroupDiscountId)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(
						hb.Button().HTML("Close").
							Class("btn btn-secondary float-start").
							Data("bs-dismiss", "modal").
							OnClick(modalCloseScript)).
					Child(buttonDelete)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().
		Children([]hb.TagInterface{
			modal,
			backdrop,
		})
}

func (c *discountDeleteController) prepareDataAndValidate() (data discountDeleteControllerData, errorMessage string) {
	data.discountID = utils.Req(c.opts.GetRequest(), "discount_id", "")

	if data.discountID == "" {
		return data, "discount id is required"
	}

	discount, err := c.opts.GetStore().DiscountFindByID(c.opts.GetRequest().Context(), data.discountID)

	if err != nil {
		c.opts.GetLogger().Error("Error. At shopstoreadmin > discountDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, "Discount not found"
	}

	if discount == nil {
		return data, "Discount not found"
	}

	data.discount = discount

	if c.opts.GetRequest().Method != "POST" {
		return data, ""
	}

	err = c.opts.GetStore().DiscountSoftDelete(c.opts.GetRequest().Context(), discount)

	if err != nil {
		c.opts.GetLogger().Error("Error. At shopstoreadmin > discountDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, "Deleting discount failed. Please contact an administrator."
	}

	data.successMessage = "discount deleted successfully."

	return data, ""

}

// =========================================================================
// == DATA
// =========================================================================

type discountDeleteControllerData struct {
	discountID     string
	discount       shopstore.DiscountInterface
	successMessage string
	//errorMessage   string
}
