package helm

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Upload struct {
	TmpDir    string `json:"chart_dir"`
	ChartName []byte `json:"chart_name"`
	checkMark helmCheck
}

type helmCheck struct {
	templates int
	helpers   int
	values    int
	chart     int
}

func (u *Upload) UploadInit() error {
	u.TmpDir = fmt.Sprintf("%d-%d", time.Now().Unix(), rand.Intn(1000))
	u.ChartName = []byte("Chart")
	u.checkMark = helmCheck{0, 0, 0, 0}
	if err := os.MkdirAll(u.TmpDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (u *Upload) UploadClose() {
	_ = os.RemoveAll(u.TmpDir)
}

type UploadSaveError struct {
	S string
}

func (e *UploadSaveError) Error() string {
	return e.S
}

func (u *Upload) UploadSave(c *gin.Context) *UploadSaveError {
	// 多文件上传
	form, _ := c.MultipartForm()
	files := form.File["files"]

	is_files_model := false
	// save Chart file
	for _, file := range files {
		fileName := strings.Split(file.Filename, ".")

		// zip 包模式
		if fileName[1] == "zip" {
			u.ChartName = []byte(fileName[0])
			if err := c.SaveUploadedFile(file, u.TmpDir+"/"+file.Filename); err != nil {
				return &UploadSaveError{"save zip failed"}
			}

			if !isZip(u.TmpDir + "/" + file.Filename) {
				return &UploadSaveError{"zip package error"}
			}

			if err := unzip(u.TmpDir+"/"+file.Filename, u.TmpDir+"/"); err != nil {
				log.Printf("unzip package error: %s", err)
				return &UploadSaveError{"unzip package error"}

			}

			continue
		}

		// 文件夹模式
		if filepath.Dir(file.Filename) != "." {
			u.ChartName = []byte(strings.Split(file.Filename, "/")[0])

			if err := unDir(file, u.TmpDir+"/"); err != nil {
				log.Printf("decode Dir  error: %s", err)
				return &UploadSaveError{"decode Dir  error"}
			}

			continue
		}

		// 文件模式
		ChartDir := u.TmpDir + "/Chart"
		if err := os.MkdirAll(ChartDir+"/templates", os.ModePerm); err != nil {
			log.Printf("mkdir templates failed: %s", err)
			return &UploadSaveError{"decode failure"}
		}

		is_files_model = true
		switch fileName[0] {
		case "values":
			if err := c.SaveUploadedFile(file, ChartDir+"/"+file.Filename); err != nil {
				return &UploadSaveError{"decode values failure"}
				//return error{ return "decode values failure"}
			}
			u.checkMark.values++
		case "Chart":
			if err := c.SaveUploadedFile(file, ChartDir+"/"+file.Filename); err != nil {
				return &UploadSaveError{"decode Chart failure"}
			}
			u.checkMark.chart++
		case "_helpers":
			if err := c.SaveUploadedFile(file, ChartDir+"/templates/"+file.Filename); err != nil {
				return &UploadSaveError{"decode _helpers failure"}
			}
			u.checkMark.helpers++
		default:

			if err := c.SaveUploadedFile(file, ChartDir+"/templates/"+file.Filename); err != nil {
				return &UploadSaveError{"decode template yaml failure"}
			}
			u.checkMark.templates++

		}
	}

	// check Chart file
	if is_files_model && u.checkMark.values*u.checkMark.chart*u.checkMark.helpers*u.checkMark.templates < 1 {
		log.Printf("values: %d, chart: %d, helpers : %d, templates: %d",
			u.checkMark.values,
			u.checkMark.chart,
			u.checkMark.helpers,
			u.checkMark.templates)

		return &UploadSaveError{"values/chart/helpers/template is require"}
	}

	return nil

}

func isZip(zipPath string) bool {

	f, err := os.Open(zipPath)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if n, err := f.Read(buf); err != nil || n < 4 {
		return false
	}
	return bytes.Equal(buf, []byte("PK\x03\x04"))
}

func unzip(archive, target string) error {

	reader, err := zip.OpenReader(archive)

	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)

		// mkdir dir
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		// read file
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		// creat target file
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err = io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil

}

func unDir(file *multipart.FileHeader, target string) error {

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	path := filepath.Join(target, file.Filename)

	dirSp := filepath.Dir(path)

	os.MkdirAll(dirSp, os.ModePerm)

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
