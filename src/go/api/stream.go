package api

import (
	//"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"ohohleo/accelerometer/input"
)

func NewStream(ch chan input.AccelGyro) error {

	// Mise en place de l'API Stream
	apiStream := rest.NewApi()

	apiStream.Use(&rest.AccessLogApacheMiddleware{})
	apiStream.Use(rest.DefaultDevStack...)

	stream, err := rest.MakeRouter(
		rest.Get("/accelerometer", Stream(ch)),
	)

	if err != nil {
		return err
	}

	apiStream.SetApp(stream)

	http.Handle("/stream/", http.StripPrefix("/stream", apiStream.MakeHandler()))

	return nil
}

func Stream(ch chan input.AccelGyro) func(w rest.ResponseWriter, r *rest.Request) {

	return func(w rest.ResponseWriter, r *rest.Request) {

		w.(http.ResponseWriter).Header().Set("Content-Type", "text/event-stream")

		w.(http.ResponseWriter).Write([]byte("data:it works!\n\n"))
		w.(http.Flusher).Flush()

		for {

			accelerometer := <-ch

			log.Printf("%+v", accelerometer)

			w.(http.ResponseWriter).Write([]byte("data:" + accelerometer.String() + "\n\n"))

			// Flush the buffer to client
			w.(http.Flusher).Flush()
		}
	}
}
