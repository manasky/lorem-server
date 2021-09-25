package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"image/jpeg"
	"io/ioutil"
	"log"
	"lorem/image"
	"lorem/manager"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Options struct {
	CacheFiles bool

	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

const (
	minSize    = 8
	maxSize    = 2000
	normalSize = 500
)

type API struct {
	mngr *manager.Manager
	pr   image.Processor
	opt  *Options
}

func New(manager *manager.Manager, imageProcessor image.Processor, options *Options) *API {
	if options == nil {
		options = &Options{
			CacheFiles: true,
			MinWidth:   minSize,
			MinHeight:  minSize,
			MaxWidth:   maxSize,
			MaxHeight:  maxSize,
		}
	} else {
		if options.MinWidth == 0 {
			options.MinWidth = minSize
		}

		if options.MinHeight == 0 {
			options.MinHeight = minSize
		}

		if options.MaxWidth == 0 {
			options.MaxWidth = maxSize
		}

		if options.MaxHeight == 0 {
			options.MaxHeight = maxSize
		}
	}

	return &API{
		mngr: manager,
		pr:   imageProcessor,
		opt:  options,
	}
}

func (a *API) SizeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	file := a.mngr.Pick(vars["category"])
	if file == "" {
		handleError(w, http.StatusNotFound, "category not found")
		return
	}

	params := r.URL.Query()
	var err error
	var width int
	if params.Get("w") != "" {
		width, err = strconv.Atoi(params.Get("w"))
		if err != nil || width < a.opt.MinWidth || width > a.opt.MaxWidth {
			handleError(w, http.StatusBadRequest, "invalid width size")
			return
		}
	}

	var height int
	if params.Get("h") != "" {
		height, err = strconv.Atoi(params.Get("h"))
		if err != nil || height < a.opt.MinHeight || height > a.opt.MaxHeight {
			handleError(w, http.StatusBadRequest, "invalid height size")
			return
		}
	}

	if width == 0 && height == 0 {
		width = normalSize
		height = normalSize
	}

	if a.opt.CacheFiles {
		cp := cachePath(file, width, height)
		if _, err := os.Stat(cp); err == nil {
			log.Printf("cached file found: %s", file)
			err = handleCacheFile(w, cp)
			if err == nil {
				return
			}
			log.Printf("error while writing cached file: %s", file)
		}
	}

	img, err := image.Decode(file)
	if err != nil {
		log.Printf("unable to decode image: %s : %s", file, err)
		handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if height != 0 && width != 0 {
		img = a.pr.CropCenter(*img, width, height)
	} else {
		img = a.pr.Resize(*img, width, height)
	}

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Printf("unable to encode image: %s", err)
		return
	}

	if a.opt.CacheFiles {
		go func() {
			cp := cachePath(file, width, height)
			err := cacheFile(cp, buffer.Bytes())
			if err != nil {
				log.Printf("error while caching file: %s : %v", cp, err)
				return
			}
			log.Printf("file cached: %s", cp)
		}()
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Printf("unable to write image: %v", err)
		return
	}
}

func (a *API) NotFound(w http.ResponseWriter, r *http.Request) {
	handleError(w, http.StatusNotFound, "not found")
}

func handleCacheFile(w http.ResponseWriter, file string) error {
	s, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(s)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(b.Bytes())))
	if _, err := w.Write(b.Bytes()); err != nil {
		return err
	}
	return nil
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

func cacheFile(cp string, i []byte) error {
	dir := path(cp)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(cp)
	if err != nil {
		return err
	}

	_, err = f.Write(i)
	if err != nil {
		return err
	}

	return nil
}

func cachePath(file string, width, height int) string {
	tmp := strings.Split(file, "/")
	tmp[len(tmp)-1] = fmt.Sprintf("%s/%dx%d/%s", manager.CacheDir, width, height, tmp[len(tmp)-1])
	return strings.Join(tmp, "/")
}

func path(file string) string {
	tmp := strings.Split(file, "/")
	return strings.Join(tmp[:len(tmp)-1], "/")
}
