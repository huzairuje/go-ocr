package primitive

import "errors"

const (
	ProcessOcrSuccess                = "processing file ocr succeeded"
	SuccessGetOcr                    = "success get record ocr"
	ParamIdIsZeroOrNullString        = "param id given value is either zero or empty"
	RecordOCrNotFound                = "record data ocr not found"
	QueryIsSuspicious                = "the query parameter given value is suspicious"
	ErrorBindBodyRequest             = "error bind body from request"
	SomethingWrongWithTheBodyRequest = "oops, something wrong with body request, please recheck!"
	SomethingWentWrong               = "oops, something went wrong!"
	ErrOcrNotFound                   = "ocr not found"
)

var ErrorArticleNotFound = errors.New(ErrOcrNotFound)
