package transport

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/Kolyan4ik99/fileServer/internal/service"
	"github.com/gin-gonic/gin"
)

type FileSaverInterface interface {
	Save(c *gin.Context)
	GetByFileName(c *gin.Context)
}

type FileSaver struct {
	fileService service.FileSaverInterface
}

func NewFileSaver(fileService service.FileSaverInterface) *FileSaver {
	return &FileSaver{fileService: fileService}
}

// Save godoc
// @Summary      Upload file
// @Description  Upload file by name
// @Accept       multipart/form-data
// @Tags         file
// @Param        file  formData  file  true  "file"
// @Success      201
// @Failure      400,500
// @Router       /file [post]
func (f *FileSaver) Save(c *gin.Context) {

	fileForm, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	open, err := fileForm.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	reader := bufio.NewReader(open)
	defer c.Request.Body.Close()

	err = f.fileService.Save(context.Background(), fileForm.Filename, reader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusCreated)
}

// GetByFileName godoc
// @Summary      Download file
// @Description  Download file by name if exist
// @Tags         file
// @Param        file_name   query     string  true  "file_name"
// @Success      200
// @Failure      404,500
// @Router       /file [get]
func (f *FileSaver) GetByFileName(c *gin.Context) {
	fileName, err := parseFileName(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	reader, err := f.fileService.GetByFileName(context.Background(), fileName)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func parseFileName(c *gin.Context) (string, error) {
	fileName := c.Query("file_name")
	if fileName == "" {
		return "", errors.New("bad file_name")
	}
	return fileName, nil
}
