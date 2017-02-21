// Package multipartbuilder provides simple streaming multipart builder.
//
// Usage:
//
//   builder := New()
//   builder.AddField("field", "value")
//
//   // or use chaining:
//   builder.
//     AddReader("reader", strings.NewReader("Some reader")).
//     AddFile("file", "path/to/file.bin")
//
//   // finalize builder (it should not be used anymore after this);
//   // any errors will be returned on bodyReader usage (Read/Close):
//   contentType, bodyReader := builder.Build()
//
//   // for proper cleanup, returned bodyReader should be used at least once,
//   // so at least close it (multiple closes are fine):
//   defer bodyReader.Close()
//
//   // finally, use built reader:
//   resp, err := http.Post("https://test.com/", contentType, bodyReader)
//   if err != nil {
//     // handle error
//   }
package multipartbuilder
