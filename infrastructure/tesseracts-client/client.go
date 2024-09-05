package tesseracts_client

import (
	"github.com/otiai10/gosseract/v2"
	"go-ocr/infrastructure/config"
)

func NewClient() *gosseract.Client {
	languagesAvailable := config.Conf.TesseractsConfig.Languages
	client := gosseract.NewClient()
	client.Languages = languagesAvailable
	return client
}
