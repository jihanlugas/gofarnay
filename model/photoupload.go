package model

import (
	"gofarnay/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	ORDER = 1
)

type Photoupload struct {
	PhotouploadId  int       `json:"photouploadId"`
	RefType        int       `json:"refType"`
	RefId          int       `json:"refId"`
	FilePath       string    `json:"filePath"`
	FilePathResize string    `json:"filePathResize"`
	FolderName     string    `json:"folderName"`
	FileName       string    `json:"fileName"`
	AltName        string    `json:"altName"`
	FileExt        string    `json:"fileExt"`
	Size           int       `json:"size"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// GetFileDetail  return (filepath, filename, extension)
func GetFileDetail(name string) (string, string, string) {
	re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
	match := re.FindStringSubmatch(name)
	return match[1], match[2], match[3]
}

func (p *Photoupload) UploadPhoto(r *http.Request, formName string) error {
	u := Uploadinc{RefType: p.RefType}
	if err := u.GetUploadincToUpload(); err != nil {
		return err
	}
	if err := u.RunningInc(); err != nil {
		return err
	}
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile(formName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, p.FileName, p.FileExt = GetFileDetail(handler.Filename)
	p.FilePath = u.FolderName

	fileLocation := filepath.Join(p.FilePath, p.FileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, file); err != nil {
		return err
	}

	//tempFile, err := ioutil.TempFile(p.FilePath, "*.png")
	//if err != nil {
	//	return err
	//}
	//defer tempFile.Close()
	//
	//fileBytes, err := ioutil.ReadAll(file)
	//if err != nil {
	//	return err
	//}
	//tempFile.Write(fileBytes)









	if err := p.CreatePhotoupload(); err != nil {
		return err
	}
	return nil
}

func (p *Photoupload) CreatePhotoupload() error {
	db := config.DbConn()
	defer db.Close()

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	res, err := db.Exec("INSERT INTO photouploads(ref_type, ref_id, file_path, file_path_resize, folder_name, file_name, alt_file, file_ext, size, width, height, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)", p.RefType, p.RefId, p.FilePath, p.FilePathResize, p.FolderName, p.FileName, p.AltName, p.FileExt, p.Size, p.Width, p.Height, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		return err
	}

	p.PhotouploadId = int(lid)

	return nil
}
