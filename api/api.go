package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"image/jpeg"
	"log"
	"lorem/image"
	"lorem/manager"
	"net/http"
	"strconv"
)

type Options struct {
	MinWidth int
	MaxWidth int
	MinHeight int
	MaxHeight int
}

const (
	minSize = 8
	maxSize = 2000
)

type API struct {
	mngr *manager.Manager
	pr image.Processor
	opt *Options
}

func New(manager *manager.Manager, imageProcessor image.Processor, options *Options) *API {
	if options == nil {
		options = &Options{
			MinWidth: minSize,
			MinHeight: minSize,
			MaxWidth: maxSize,
			MaxHeight: maxSize,
		}
	}

	return &API{
		mngr: manager,
		pr: imageProcessor,
	}
}

func (a *API) SizeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)

	file := a.mngr.Pick(vars["category"])
	if file == "" {
		handleError(w, http.StatusNotFound, "category not found")
		return
	}

	width, err := strconv.ParseInt(vars["width"], 10, 32)
	if err != nil || width < 8 {
		handleError(w, http.StatusBadRequest, "invalid width size")
		return
	}

	height, err := strconv.ParseInt(vars["height"], 10, 32)
	if err != nil || height < 8 {
		handleError(w, http.StatusBadRequest, "invalid height size")
		return
	}

	img, err := image.Decode(file)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	img = a.pr.CropCenter(*img, int(width), int(height))


	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
		return
	}
}

func (a *API) NotFound(w http.ResponseWriter, r *http.Request) {
	handleError(w, http.StatusNotFound, "not found")
}

func handleError(w http.ResponseWriter, code int, msg string) {
	b, err := json.Marshal(&Error{
		Message: msg,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while marshaling response: %s", err)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
	return
}