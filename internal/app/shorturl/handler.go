package shorturl

import (
	"encoding/json"
	"fmt"
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/entity"
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/errors"
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/resource"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
)

type CreateShortUrlReq struct {
	Url           string `json:"url" validate:"required"`
	ExpireSeconds int64  `json:"expire_seconds" validate:"min=0"`
}

type CreateShortUrlRsp struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	ShortUrl string `json:"short_url"`
}

type GetUrlInfoRsp struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Info *entity.UrlDetailInfo `json:"info"`
}

func (app *App) createShortUrl(w http.ResponseWriter, r *http.Request) {
	var req CreateShortUrlReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWithErr(w, &errors.StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("json decode err: %v", req),
		})
		return
	}
	if err := validator.New().Struct(req); err != nil {
		responseWithErr(w, &errors.StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("req: %+v validate err: %v", req, err),
		})
		return
	}
	defer r.Body.Close()
	log.Printf("req: %+v\n", req)
	ret, err := resource.Storage.Shorten(req.Url, req.ExpireSeconds)
	if err != nil {
		responseWithErr(w, &errors.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		})
		return
	}
	responseWithJson(w, http.StatusCreated, &CreateShortUrlRsp{
		Code:     0,
		Msg:      "ok",
		ShortUrl: ret,
	})
}

func (app *App) getUrlInfo(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("sid")
	log.Printf("sid: %+v", sid)
	if info, err := resource.Storage.ShortLinkInfo(sid); err != nil {
		responseWithErr(w, &errors.StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		})
		return
	} else {
		responseWithJson(w, http.StatusOK, &GetUrlInfoRsp{
			Code: 0,
			Msg:  "ok",
			Info: info,
		})
	}
}

func (app *App) redirect(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["sid"]
	log.Printf("redirect...,sid: %v", sid)
	url, err := resource.Storage.UnShorten(sid)
	if err != nil {
		responseWithErr(w, err)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func (app *App) counter(writer http.ResponseWriter, request *http.Request) {
	log.Println("counter...")
}

func responseWithErr(w http.ResponseWriter, err error) {
	switch val := err.(type) {
	case errors.Error:
		log.Printf("Http %d - %s\n", val.Status(), val.Error())
		responseWithJson(w, val.Status(), val.Error())
	default:
		responseWithJson(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	rspBody, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(rspBody)
}
