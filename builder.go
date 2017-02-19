package multipartbuilder

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Builder represents HTTP multipart request builder.
//
// It is not thread-safe.
type Builder struct {
	body *bytes.Buffer
	form *multipart.Writer
	errs []error
}

// New constructs new multipart Builder.
func New() *Builder {
	body := bytes.NewBuffer(nil)
	return &Builder{
		body: body,
		form: multipart.NewWriter(body),
	}
}

// WriteField writes field.
//
// This method returns current Builder for chaining.
func (b *Builder) WriteField(name string, value string) *Builder {
	if err := b.form.WriteField(name, value); err != nil {
		b.errs = append(b.errs, fmt.Errorf("multipartbuilder: failed to write field %s=%s: %s", name, value, err.Error()))
	}
	return b
}

// WriteFields writes multiple fields.
// It is intended to work with net/url.Values.
//
// This method returns current Builder for chaining.
func (b *Builder) WriteFields(fields map[string][]string) *Builder {
	for name, values := range fields {
		for _, value := range values {
			b.WriteField(name, value)
		}
	}
	return b
}

// SlurpFile reads filePath as fieldName.
//
// This method returns current Builder for chaining.
func (b *Builder) SlurpFile(fieldName, filePath string) *Builder {
	f, err := os.Open(filePath)
	if err != nil {
		b.errs = append(b.errs, fmt.Errorf("multipartbuilder: failed to open file %s for field %s: %s", filePath, fieldName, err.Error()))
		return b
	}
	defer f.Close()

	return b.SlurpReader(fieldName, filepath.Base(filePath), f)
}

// SlurpReader reads reader as fieldName with fileName.
//
// This method returns current Builder for chaining.
func (b *Builder) SlurpReader(fieldName, fileName string, reader io.Reader) *Builder {
	w, err := b.form.CreateFormFile(fieldName, fileName)
	if err != nil {
		b.errs = append(b.errs, fmt.Errorf("multipartbuilder: failed to create form file %s (%s): %s", fieldName, fileName, err.Error()))
		return b
	}

	if _, err = io.Copy(w, reader); err != nil {
		b.errs = append(b.errs, fmt.Errorf("multipartbuilder: failed to copy reader for form file %s (%s): %s", fieldName, fileName, err.Error()))
	}
	return b
}

// Build finalizes builder and returns Content-Type and multipart body reader.
//
// It does not check, if is called multiple times, so be careful.
func (b *Builder) Build() (string, io.Reader, error) {
	if err := multiError(b.errs); err != nil {
		return "", nil, err
	}

	if err := b.form.Close(); err != nil {
		return "", nil, fmt.Errorf("multipartbuilder: failed to close multipart writer: %s", err.Error())
	}

	return b.form.FormDataContentType(), b.body, nil
}

// BuildRequest is a convenience method for creating HTTP request from builder.
// It is just a wrapper for .Build() method.
//
// It does not check, if is called multiple times, so be careful.
func (b *Builder) BuildRequest(method string, url string) (*http.Request, error) {
	ctype, body, err := b.Build()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("multipartbuilder: failed to create HTTP request for %s %s: %s", method, url, err.Error())
	}

	req.Header.Set("Content-Type", ctype)
	return req, nil
}
