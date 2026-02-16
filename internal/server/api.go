package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tom-draper/log-analyzer/pkg/parse"
)

type Data struct {
	Extraction []parse.Extraction `json:"extraction"`
	Locations  map[string]string  `json:"locations"`
	Config     *parse.Config      `json:"config"`
}

func Start(data *Data, port string, dashboardFS fs.FS) {
	sub, _ := fs.Sub(dashboardFS, "dashboard/dist")
	fileServer := http.FileServer(http.FS(sub))

	r := chi.NewRouter()
	// Return lines data when requested by dashboard on load
	r.Get("/data", func(w http.ResponseWriter, r *http.Request) {
		jsonString, err := json.Marshal(data)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
		}
		w.Write(jsonString)
	})
	// Serve embedded dashboard static files
	r.Handle("/*", fileServer)

	srv := &http.Server{Addr: "127.0.0.1:" + port, Handler: r}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	fmt.Printf("Dashboard running at http://localhost:%s/\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "server error: %s\n", err)
		os.Exit(1)
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
