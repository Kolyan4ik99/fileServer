package repository

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Kolyan4ik99/fileServer/internal/model"
)

type FileSaverInterface interface {
	Save(fileName string, splitLimit int, reader io.Reader) error
	GetByFileName(fileName string) (io.Reader, error)
}

type FileType map[string][]model.FileInfo

type FileSaver struct {
	maxId     int
	dir       string
	filesById FileType
}

func NewFileSaver(folderForFiles string) (*FileSaver, error) {
	file := &FileSaver{}

	err := file.fillFilesStruct(folderForFiles)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *FileSaver) Save(fileName string, splitLimit int, reader io.Reader) error {
	err := f.fileSplit(fileName, splitLimit, reader)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (f *FileSaver) fileSplit(fileName string, splitLimit int, reader io.Reader) error {
	f.filesById[fileName] = make([]model.FileInfo, 0)
	readBytes := make([]byte, splitLimit)
	fmt.Println(len(readBytes))
	for i := 1; true; i++ {
		read, err := reader.Read(readBytes)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("read size =[%d]", read)

		createFileName := fmt.Sprintf("%s_%d_%d_%d", fileName, 1, i, read)
		fmt.Println(createFileName)
		if read < splitLimit {
			nullBytes := make([]byte, splitLimit)

			// Вот тут бы функцию Copy
			for j := 0; j < read; j++ {
				nullBytes[j] = readBytes[j]
			}

			readBytes = nullBytes
		}

		f.filesById[fileName] = append(f.filesById[fileName], model.FileInfo{
			Name: createFileName,
			Id:   i,
			Size: read,
		})
		err = os.WriteFile(fmt.Sprintf("%s/%s", f.dir, createFileName), readBytes, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FileSaver) GetByFileName(fileName string) (io.Reader, error) {
	files, exist := f.filesById[fileName]
	if !exist {
		return nil, errors.New("bad file name")
	}

	newFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	var tmp model.TMPFILE
	tmp = files
	sort.Sort(tmp)
	defer os.Remove(fileName)

	for i := 0; i < len(files); i++ {
		readBytes := make([]byte, files[i].Size)
		file, err := os.OpenFile(fmt.Sprintf("%s/%s", f.dir, files[i].Name), os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
		}

		_, err = file.Read(readBytes)
		if err != nil {
			return nil, err
		}

		_, err = newFile.Write(readBytes)
		if err != nil {
			return nil, err
		}
	}

	_, err = newFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(newFile), nil
}

func (f *FileSaver) fillFilesStruct(dirName string) error {
	files := make([]string, 0)
	f.filesById = make(FileType)
	f.dir = dirName
	err := filepath.Walk(dirName, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			matched, err := regexp.MatchString("_[1-9][0-9]{0,}_[1-9][0-9]{0,}_[1-9][0-9]{0,}", info.Name())
			if err != nil {
				return err
			}
			if matched {
				files = append(files, info.Name())
			}
		}

		return nil
	})

	for i := range files {
		partsFile := strings.Split(files[i], "_")
		fileName := partsFile[0]

		fileId, err := strconv.Atoi(partsFile[2])
		if err != nil {
			return err
		}
		size, err := strconv.Atoi(partsFile[3])
		if err != nil {
			return err
		}

		_, exist := f.filesById[fileName]
		if !exist {
			f.filesById[fileName] = make([]model.FileInfo, 0)
		}

		f.filesById[fileName] = append(f.filesById[fileName], model.FileInfo{
			Name: files[i],
			Id:   fileId,
			Size: size,
		})
	}
	return err
}
