package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"lorem/api"
	"lorem/image"
	"lorem/manager"
	"net/http"
	"os"
	"path"
)

func init() {
	flags := pflag.NewFlagSet(path.Base(os.Args[0]), pflag.ContinueOnError)

	flags.String("host", "127.0.0.1:8080", "host:port for the HTTP server")
	flags.String("dir", "./images", "directory of images")
	flags.String("cache", "true", "cache processed files")
	flags.String("cdn", "", "cdn domain (works only on cache mode). leave empty to serve the files by app")
	flags.String("min-width", "8", "minimum supported width")
	flags.String("min-height", "8", "minimum supported height")
	flags.String("max-width", "2000", "maximum supported width")
	flags.String("max-height", "2000", "maximum supported height")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Panic("failed to parse arguments")
	}

	err = viper.BindPFlags(flags)
	if err != nil {
		log.Panic("failed to bind flags")
	}

	viper.AutomaticEnv()
}

func main() {
	m, err := manager.New(viper.GetString("dir"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%d items loaded", m.Total())

	a := api.New(m, &image.Imaging{}, &api.Options{
		CacheFiles: viper.GetBool("cache"),
		CDN: viper.GetString("cdn"),
		MinWidth: viper.GetInt("min-width"),
		MaxWidth: viper.GetInt("max-width"),
		MinHeight: viper.GetInt("min-height"),
		MaxHeight: viper.GetInt("max-height"),
	})

	r := mux.NewRouter()
	r.HandleFunc("/image/{category}", a.SizeHandler).Methods(http.MethodGet)
	r.HandleFunc("/image", a.SizeHandler).Methods(http.MethodGet)
	r.PathPrefix("/").HandlerFunc(a.NotFound)

	log.Printf("listening on %s", viper.GetString("host"))
	err = http.ListenAndServe(viper.GetString("host"), handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"}),
	)(r))
	if err != nil {
		log.Panic("HTTP server failed to start")
	}
}
