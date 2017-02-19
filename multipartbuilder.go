// Package multipartbuilder provides multipart reader builder.
//
// Usage:
//   builder := multipartbuilder.New().
//     WriteField("field", "value").
//     WriteField("field", "another value").
//     WriteFields(map[string][]string{
//       "field":         []string{"even", "more", "values"},
//       "another_field": []string{"another value"},
//     }).
//     SlurpReader("reader", "file.bin", strings.NewReader("foo bar")).
//     SlurpFile("file", "path/to/file.bin").
//   builder.WriteField("or", "don't use chaining, doesn't matter")
//
// Then, you may either get Content-Type and body reader:
//   contentType, bodyReader, err := builder.Build()
//   if err != nil {
//     panic(err.Error()) // handle error somehow
//   }
//   resp, err := http.Post("https://test.com/", contentType, bodyReader)
//
// Or build request right away:
//   req, err := builder.BuildRequest("POST", "https://test.com/")
//   if err != nil {
//     panic(err.Error()) // handle error somehow
//   }
//   // modify request, if needed:
//   req.AddCookie(...)
//   // and, finally, execute it:
//   resp, err := http.DefaultClient.Do(req)
package multipartbuilder
