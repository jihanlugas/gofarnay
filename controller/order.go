package controller

import (
	"gofarnay/model"
	"log"
	"net/http"
	"strconv"
)

func OrderForm(w http.ResponseWriter, r *http.Request) {
	var o model.Order

	o.Name = r.FormValue("name")
	o.Email = r.FormValue("email")
	o.Phone = r.FormValue("phone")
	o.Title = r.FormValue("title")
	o.Message = r.FormValue("message")

	if err := o.CreateOrder(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var p model.Photoupload

	p.RefType = model.ORDER
	p.RefId = o.OrderId

	if err := p.UploadPhoto(r, "file"); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	o.PhotoId = p.PhotouploadId
	if err := o.UpdatePhotoId(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithSuccess(w, http.StatusCreated, o)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	log.Println("GetOrders")
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	orders, err := model.GetOrders(start, count)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithSuccess(w, http.StatusOK, orders)


}
