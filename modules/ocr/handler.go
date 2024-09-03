package ocr

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"go-ocr/infrastructure/httplib"
	"go-ocr/modules/primitive"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
)

type Http struct {
	serviceOcr ServiceInterface
}

func NewHttp(serviceOcr ServiceInterface) InterfaceHttp {
	return &Http{
		serviceOcr: serviceOcr,
	}
}

type InterfaceHttp interface {
	GroupOcr(group *gin.RouterGroup)
}

func (h *Http) GroupOcr(g *gin.RouterGroup) {
	g.POST("", h.FileUpload)
}

func (h *Http) FileUpload(c *gin.Context) {
	// Get uploaded file
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Create temporary file
	tempfile, err := os.CreateTemp("", "ocrserver"+"-")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	// Write uploaded file to the temporary file
	if _, err = io.Copy(tempfile, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(tempfile.Name())
	client.Languages = []string{"eng"}
	if langs := c.PostForm("languages"); langs != "" {
		client.Languages = strings.Split(langs, ",")
	}
	if whitelist := c.PostForm("whitelist"); whitelist != "" {
		client.SetWhitelist(whitelist)
	}

	var out string
	switch c.PostForm("format") {
	case "hocr":
		out, err = client.HOCRText()
		c.Set("EscapeHTML", false)
	default:
		out, err = client.Text()
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out = strings.Trim(out, "\n")

	fmt.Println("out:", out)

	httplib.SetSuccessResponse(c, http.StatusOK, primitive.ProcessOcrSuccess, out)
	return
}
