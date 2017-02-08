package style

import (
	"os"
	"regexp"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"github.com/Denisrudov/base64"
	"errors"
)

const UrlMatch string = `url\(["']? ?(.*\.jpg|png|gif|svg?) ?["']?\)`

type file struct {
	info         os.FileInfo
	file         *os.File
	urlReg       *regexp.Regexp
	filePath     string
	maxImageSize int64
}

func (fp *file) SetMaxImage(size int64) {
	fp.maxImageSize = size
}

/*
Create file file instance by path of a file file
 */
func NewFile(path string) (fp *file, err error) {
	reg, err := regexp.Compile(UrlMatch)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return &file{}, err
	}
	fpath, _ := filepath.Split(absPath)

	fp = &file{
		urlReg:       reg,
		filePath:     fpath,
		maxImageSize: 1024 * 1024,
	}

	err = fp.Open(path)
	return fp, err
}

/*
 Open file for encode
 */
func (fp *file) Open(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	defer file.Close()
	if err != nil {
		return err
	}
	fp.file = file
	fp.info, err = file.Stat()
	return err
}

/*
   Encode all urls in file to base64
   Works only for local files
 */
func (fp *file) Encode() error {

	content, err := ioutil.ReadFile(fp.file.Name())
	if err != nil {
		return err
	}

	updatedContent := fp.urlReg.ReplaceAllFunc(content, func(content []byte) []byte {

		subMatches := fp.urlReg.FindStringSubmatch(string(content))

		encoded, err := fp.encodeFile(subMatches[1])

		if err != nil {
			return content
		}

		return []byte(fmt.Sprintf("url(\"%s\")", encoded))
	})

	return ioutil.WriteFile(fp.file.Name(), updatedContent, 0)
}

func (fp *file) encodeFile(i string) (string, error) {

	imagePath := filepath.Join(fp.filePath, i)

	imFile, err := os.Open(imagePath)
	defer imFile.Close()
	if err != nil {
		return "", err
	}
	imInfo, err := imFile.Stat()
	if err != nil {
		return "", err
	}

	if imInfo.Size() > fp.maxImageSize {
		return "", errors.New("Big file")
	}

	imageFile, err := base64.NewImageFile(imagePath)

	if err != nil {
		return "", err
	}

	encodedContent, mimeType, err := imageFile.EncodeBase64()
	if err != nil {
		return "", err
	}

	encoded := fmt.Sprintf("data:%s;base64,%s", mimeType, encodedContent)

	return encoded, err
}
