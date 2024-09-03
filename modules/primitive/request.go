package primitive

type OcrRequest struct {
	Image string `form:"-"`
	Type  string `json:"type"`
}
