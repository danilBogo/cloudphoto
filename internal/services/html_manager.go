package services

import (
	"bytes"
	"cloudphoto/internal/constants"
	"fmt"
	"html/template"
	"os"
	"path"
)

type HtmlManager struct {
	currentDir string
}

func NewHtmlManager() (*HtmlManager, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &HtmlManager{currentDir: currentDir}, nil
}

type Photo struct {
	URL   string
	Title string
}

func (hm HtmlManager) GetErrorHtml() ([]byte, error) {
	result, err := os.ReadFile(path.Join(hm.currentDir, constants.FolderName, constants.ErrorHtml))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (hm HtmlManager) GetIndexHtml(count int) ([]byte, error) {
	htmlTemplate, err := os.ReadFile(path.Join(hm.currentDir, constants.FolderName, constants.IndexHtml))
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("index").Parse(string(htmlTemplate))
	var buffer bytes.Buffer
	numbers := make([]int, count)
	for i := range numbers {
		numbers[i] = i + 1
	}

	err = tmpl.Execute(&buffer, numbers)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (hm HtmlManager) GetAlbumHtml(photos []Photo) ([]byte, error) {
	htmlTemplate, err := os.ReadFile(path.Join(hm.currentDir, constants.FolderName, constants.AlbumHtml))
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("index").Parse(string(htmlTemplate))
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, photos)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func GetAlbumName(index int) string {
	return fmt.Sprintf("album%v.html", index)
}
