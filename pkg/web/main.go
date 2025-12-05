package web

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Headers map[string]string

type Response struct {
	Status      int
	ContentType string
	Content     io.Reader
	Headers     Headers
}

func (response *Response) Write(rw http.ResponseWriter) {
	if response != nil {
		if response.ContentType != "" {
			rw.Header().Set("Content-Type", response.ContentType)
		}
		for k, v := range response.Headers {
			rw.Header().Set(k, v)
		}
		rw.WriteHeader(response.Status)
		_, err := io.Copy(rw, response.Content)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

type Action func(r *http.Request) *Response

func (a Action) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	response := a(r)
	response.Write(rw)
}

func Error(status int, err error, headers Headers) *Response {
	return &Response{
		Status:  status,
		Content: bytes.NewBufferString(err.Error()),
		Headers: headers,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func ErrorJSON(status int, err error, headers Headers) *Response {
	errResp := errorResponse{
		Error: err.Error(),
	}

	b, err := json.Marshal(errResp)
	if err != nil {
		return Error(http.StatusInternalServerError, err, headers)
	}

	return &Response{
		Status:      status,
		ContentType: "application/json",
		Content:     bytes.NewBuffer(b),
		Headers:     headers,
	}
}

func Data(status int, content []byte, headers Headers) *Response {
	return &Response{
		Status:  status,
		Content: bytes.NewBuffer(content),
		Headers: headers,
	}
}

func DataJSON(status int, v any, headers Headers) *Response {

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ErrorJSON(http.StatusInternalServerError, err, headers)
	}

	return &Response{
		Status:      status,
		ContentType: "application/json",
		Content:     bytes.NewBuffer(b),
		Headers:     headers,
	}
}

func Empty(status int) *Response {
	return Data(status, []byte(""), nil)
}

func Redirect(location string) *Response {
	return &Response{
		Status:  http.StatusFound,
		Content: bytes.NewBuffer(nil),
		Headers: map[string]string{
			"Location": location,
		},
	}
}

func HTML(status int, t *template.Template, template string, data any, headers Headers) *Response {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, template, data); err != nil {
		log.Println(err)
		return Empty(http.StatusInternalServerError)
	}
	return &Response{
		Status:      status,
		ContentType: "text/html",
		Content:     &buf,
		Headers:     headers,
	}
}

func TemplateParseFSRecursive(
	templates fs.FS,
	ext string,
	nonRootTemplateNames bool,
	funcMap template.FuncMap) (*template.Template, error) {

	root := template.New("")
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ext) {
			if err != nil {
				return err
			}
			b, err := fs.ReadFile(templates, path)
			if err != nil {
				return err
			}
			name := ""
			if nonRootTemplateNames {
				parts := strings.Split(path, string(os.PathSeparator))
				name = strings.Join(parts[1:], string(os.PathSeparator))
			}
			t := root.New(name).Funcs(funcMap)
			_, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return root, err

}
