package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tom-draper/log-analyzer/pkg/parse"
)

type Data struct {
	Extraction []parse.Extraction `json:"extraction"`
	Locations  map[string]string  `json:"locations"`
	Config     *parse.Config      `json:"config"`
}

func Start(data *Data) {
	r := chi.NewRouter()
	// Serve dashboard index.html
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filePath := filepath.Join(workDir, "dashboard", "dist", "index.html")
		http.ServeFile(w, r, filePath)
	})
	// Return any asset used by index.html
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "*")
		workDir, _ := os.Getwd()
		filePath := filepath.Join(workDir, "dashboard", "dist", filename)
		http.ServeFile(w, r, filePath)
	})
	// Return lines data when requested by dashboard on load
	r.Get("/data", func(w http.ResponseWriter, r *http.Request) {
		jsonString, err := json.Marshal(data)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
		}
		w.Write(jsonString)
	})

	fmt.Println("Dashboard running at http://localhost:3000/")
	http.ListenAndServe("127.0.0.1:3000", r)
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
