package primitive

type OcrRequest struct {
	Image       string `form:"-"`
	Type        string `form:"type" validate:"required"`
	HOCREnabled string `form:"hocrEnabled"`
}
