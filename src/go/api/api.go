package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

func New() error {

	// Mise en place de l'API
	api := rest.NewApi()

	api.Use(&rest.AccessLogApacheMiddleware{})
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/test", func(w rest.ResponseWriter, req *rest.Request) {
			w.WriteJson(map[string]string{"Body": "TEST"})
		}),
	)

	if err != nil {
		return err
	}

	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	// Gestion des pages web statique
	http.Handle("/", http.StripPrefix("/",
		http.FileServer(http.Dir("/home/lmartin/accelerometer/src/web"))))

	return nil
}
