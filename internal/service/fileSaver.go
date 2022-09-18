package service

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/Kolyan4ik99/fileServer/internal/repository"
)

type FileSaverInterface interface {
	Save(ctx context.Context, fileName string, reader io.Reader) error
	GetByFileName(ctx context.Context, fileName string) (io.Reader, error)
}

type FileSaver struct {
	repo repository.FileSaverInterface
}

func NewFileSaver(repo repository.FileSaverInterface) *FileSaver {
	return &FileSaver{repo: repo}
}

func (f *FileSaver) Save(ctx context.Context, fileName string, reader io.Reader) error {
	newFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer os.Remove(fileName)
	_, err = newFile.ReadFrom(reader)
	if err != nil {
		return err
	}

	_, err = newFile.Seek(0, 0)
	if err != nil {
		return err
	}

	return f.repo.Save(fileName, 1<<20, bufio.NewReader(newFile))
}

func (f *FileSaver) GetByFileName(ctx context.Context, fileName string) (io.Reader, error) {
	return f.repo.GetByFileName(fileName)
}
