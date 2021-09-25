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

	a := api.New(m, &image.Imaging{}, nil)

	r := mux.NewRouter()
	r.HandleFunc("/{category}/{width:[0-9]+}/{height:[0-9]+}", a.SizeHandler).Methods(http.MethodGet)
	r.HandleFunc("/{width:[0-9]+}/{height:[0-9]+}", a.SizeHandler).Methods(http.MethodGet)
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
