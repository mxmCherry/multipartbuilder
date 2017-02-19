# multipartbuilder [![Build Status](https://travis-ci.org/mxmCherry/multipartbuilder.svg?branch=master)](https://travis-ci.org/mxmCherry/multipartbuilder)

Simple HTTP multipart request builder for Go (Golang)

# Usage

```go
builder := multipartbuilder.New().
	WriteField("field", "value").
	WriteField("field", "another value").
	WriteFields(map[string][]string{
		"field":         []string{"even", "more", "values"},
		"another_field": []string{"another value"},
	}).
	SlurpReader("reader", "file.bin", strings.NewReader("foo bar")).
	SlurpFile("file", "path/to/file.bin").
builder.WriteField("or", "don't use chaining, doesn't matter")
```

Then, you may either get Content-Type and body reader:

```go
	contentType, bodyReader, err := builder.Build()
	if err != nil {
		panic(err.Error()) // handle error somehow
	}
	resp, err := http.Post("https://test.com/", contentType, bodyReader)
```

Or build request right away:

```go
	req, err := builder.BuildRequest("POST", "https://test.com/")
	if err != nil {
		panic(err.Error()) // handle error somehow
	}
	// modify request, if needed:
	req.AddCookie(...)
	// and, finally, execute it:
	resp, err := http.DefaultClient.Do(req)
```
