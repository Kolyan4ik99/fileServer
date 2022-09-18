package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	fileService FileSaverInterface
}

func NewHandler(fileService FileSaverInterface) *Handler {
	return &Handler{fileService: fileService}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Handle(http.MethodPost, "/file", h.fileService.Save)
	router.Handle(http.MethodGet, "/file", h.fileService.GetByFileName)

	return router
}
