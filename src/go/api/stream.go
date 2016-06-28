package api

import (
	//"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"ohohleo/accelerometer/input"
)

func NewStream(ch chan input.Accelerometer) error {

	// Mise en place de l'API Stream
	api_stream := rest.NewApi()

	api_stream.Use(&rest.AccessLogApacheMiddleware{})
	api_stream.Use(rest.DefaultDevStack...)

	stream, err := rest.MakeRouter(
		rest.Get("/accelerometer", Stream(ch)),
	)

	if err != nil {
		return err
	}

	api_stream.SetApp(stream)

	http.Handle("/stream/", http.StripPrefix("/stream", api_stream.MakeHandler()))

	return nil
}

func Stream(ch chan input.Accelerometer) func(w rest.ResponseWriter, r *rest.Request) {

	return func(w rest.ResponseWriter, r *rest.Request) {

		w.(http.ResponseWriter).Header().Set("Content-Type", "text/event-stream")

		w.(http.ResponseWriter).Write([]byte("data:it works!\n\n"))
		w.(http.Flusher).Flush()

		for {

			accelerometer := <-ch

			w.(http.ResponseWriter).Write([]byte("data:" + accelerometer.String() + "\n\n"))

			// Flush the buffer to client
			w.(http.Flusher).Flush()
		}
	}
}
