package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// OKResponse ...
type OKResponse struct {
	HTTPStatusCode int    `json:"-"`
	Data           string `json:"data"`
}

// Render ...
func (ok *OKResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, ok.HTTPStatusCode)
	return nil
}

// SuccessResponse ...
func SuccessResponse(data string) render.Renderer {
	return &OKResponse{
		HTTPStatusCode: 200,
		Data:           data,
	}
}

// ErrResponse ...
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render ...
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrNotFound ...
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

// ErrInvalidRequest ...
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrNodeResponse ...
func ErrNodeResponse(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Invalid node response.",
		ErrorText:      err.Error(),
	}
}

// ErrRender ...
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
