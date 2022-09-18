package app

import (
	"github.com/Kolyan4ik99/fileServer/internal/repository"
	"github.com/Kolyan4ik99/fileServer/internal/service"
	"github.com/Kolyan4ik99/fileServer/internal/transport"

	_ "github.com/Kolyan4ik99/fileServer/docs"
)

// @title           File service
// @version         1.0
// @description     REST API for upload and download Files.
// @host      localhost:8080
// @BasePath  /

type FileServer struct {
	folder string
	addr   string
}

func NewFileServer(folder string, addr string) *FileServer {
	return &FileServer{folder: folder, addr: addr}
}

func (f *FileServer) Run() error {
	saverRepo, err := repository.NewFileSaver(f.folder)
	if err != nil {
		return err
	}

	saverService := service.NewFileSaver(saverRepo)
	saverTransport := transport.NewFileSaver(saverService)

	handler := transport.NewHandler(saverTransport)

	return handler.InitRoutes().Run(f.addr)
}
