package model

import (
	"database/sql"
	"fmt"
	"gofarnay/config"
	"os"
	"strconv"
	"time"
)

const RUNNING_LIMIT = 1000

type Uploadinc struct {
	UploadincId int       `json:"uploadincId"`
	RefType     int       `json:"refType"`
	FolderInc   int       `json:"folderInc"`
	FolderName  string    `json:"folderName"`
	Running     int       `json:"running"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (u *Uploadinc) GetUploadincToUpload() error {
	if err := u.GetUploadincReftype(); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("ErrNoRows")
			if err := u.CreateUploadinc(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if u.Running >= RUNNING_LIMIT {
		fmt.Println("RUNNING_LIMIT")
		if err := u.CreateUploadinc(); err != nil {
			return err
		}
	}

	return nil
}

func (u *Uploadinc)RunningInc() error {
	db := config.DbConn()
	defer db.Close()

	u.Running = u.Running + 1

	_, err := db.Exec("UPDATE uploadincs set running = ? where uploadinc_id = ?", u.Running, u.UploadincId)

	return err
}

func (u *Uploadinc) GetUploadincReftype() error {
	db := config.DbConn()
	defer db.Close()

	return db.QueryRow("SELECT uploadinc_id, ref_type, folder_inc, folder_name, running FROM uploadincs where ref_type = ? ORDER BY folder_inc DESC",
		u.RefType).Scan(&u.UploadincId, &u.RefType, &u.FolderInc, &u.FolderName, &u.Running )
}

func (u *Uploadinc) CreateUploadinc() error {
	db := config.DbConn()
	defer db.Close()

	now := time.Now()

	u.FolderInc = u.FolderInc + 1
	u.FolderName = "uploads/" + strconv.Itoa(u.RefType) + "/" + now.Format("2006-01-02") + "/"  + strconv.Itoa(u.FolderInc)
	u.CreatedAt = now
	u.UpdatedAt = now
	u.Running = 0

	res, err := db.Exec("INSERT INTO uploadincs(ref_type, folder_inc, folder_name, running, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)", u.RefType, u.FolderInc, u.FolderName, u.Running, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.UploadincId = int(lid)
	err = os.MkdirAll(u.FolderName, 0755)
	if err != nil {
		return err
	}

	return nil
}
