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
type Builder interface {

	// WriteField writes field.
	//
	// This method returns current Builder for chaining.
	WriteField(name string, value string) Builder

	// WriteFields writes multiple fields.
	// It is intended to work with net/url.Values.
	//
	// This method returns current Builder for chaining.
	WriteFields(fields map[string][]string) Builder

	// SlurpFile reads filePath as fieldName.
	//
	// This method returns current Builder for chaining.
	SlurpFile(fieldName, filePath string) Builder

	// SlurpReader reads reader as fieldName with fileName.
	//
	// This method returns current Builder for chaining.
	SlurpReader(fieldName, fileName string, reader io.Reader) Builder

	// Build finalizes builder and returns Content-Type and multipart body reader.
	//
	// It does not check, if is called multiple times, so be careful.
	Build() (contentType string, body io.Reader, err error)

	// BuildRequest is a convenience method for creating HTTP request from builder.
	// It is just a wrapper for .Build() method.
	//
	// It does not check, if is called multiple times, so be careful.
	BuildRequest(method string, url string) (*http.Request, error)
}

func New() Builder {
	body := bytes.NewBuffer(nil)
	return &builder{
		body: body,
		form: multipart.NewWriter(body),
	}
}

// ----------------------------------------------------------------------------

type builder struct {
	body *bytes.Buffer
	form *multipart.Writer
	errs []error
}

func (b *builder) WriteField(name string, value string) Builder {
	if err := b.form.WriteField(name, value); err != nil {
		b.errs = append(b.errs, fmt.Errorf("multipartbuilder: failed to write field %s=%s: %s", name, value, err.Error()))
	}
	return b
}

func (b *builder) WriteFields(fields map[string][]string) Builder {
	for name, values := range fields {
		for _, value := range values {
			b.WriteField(name, value)
		}
	}
	return b
}

func (b *builder) SlurpFile(fieldName, filePath string) Builder {
	f, err := os.Open(filePath)
	if err != nil {
		b.errs = append(b.errs, fmt.Errorf("multipartbuilder: failed to open file %s for field %s: %s", filePath, fieldName, err.Error()))
		return b
	}
	defer f.Close()

	return b.SlurpReader(fieldName, filepath.Base(filePath), f)
}

func (b *builder) SlurpReader(fieldName, fileName string, reader io.Reader) Builder {
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

func (b *builder) Build() (string, io.Reader, error) {
	if err := multiError(b.errs); err != nil {
		return "", nil, err
	}

	if err := b.form.Close(); err != nil {
		return "", nil, fmt.Errorf("multipartbuilder: failed to close multipart writer: %s", err.Error())
	}

	return b.form.FormDataContentType(), b.body, nil
}

func (b *builder) BuildRequest(method string, url string) (*http.Request, error) {
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
