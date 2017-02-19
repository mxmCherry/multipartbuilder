package multipartbuilder_test

import (
	"io"
	"io/ioutil"
	"strings"

	. "github.com/mxmCherry/multipartbuilder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Builder", func() {
	var subject *Builder

	BeforeEach(func() {
		subject = New()
	})

	Context("fields", func() {

		It("should write field", func() {
			subject.WriteField("field_name", "FIELD_VALUE")

			_, bodyReader, err := subject.Build()
			Expect(err).NotTo(HaveOccurred())

			body := read(bodyReader)
			Expect(body).To(ContainSubstring("field_name"))
			Expect(body).To(ContainSubstring("FIELD_VALUE"))
		})

		It("should write fields", func() {
			subject.WriteField("field_name", "FIELD_VALUE")
			subject.WriteFields(map[string][]string{
				"field_1": {"FIELD_1_VALUE_1"},
				"field_2": {"FIELD_2_VALUE_1", "FIELD_2_VALUE_2"},
			})

			_, bodyReader, err := subject.Build()
			Expect(err).NotTo(HaveOccurred())

			body := read(bodyReader)
			Expect(body).To(ContainSubstring("field_1"))
			Expect(body).To(ContainSubstring("FIELD_1_VALUE_1"))
			Expect(body).To(ContainSubstring("field_2"))
			Expect(body).To(ContainSubstring("FIELD_2_VALUE_1"))
			Expect(body).To(ContainSubstring("FIELD_2_VALUE_2"))
		})

	}) // fields context

	Context("slurping", func() {

		It("should slurp reader", func() {
			subject.SlurpReader("field_name", "file.bin", strings.NewReader("READER_CONTENTS"))

			_, bodyReader, err := subject.Build()
			Expect(err).NotTo(HaveOccurred())

			body := read(bodyReader)
			Expect(body).To(ContainSubstring("field_name"))
			Expect(body).To(ContainSubstring("file.bin"))
			Expect(body).To(ContainSubstring("READER_CONTENTS"))
		})

		It("should slurp file", func() {
			subject.SlurpFile("field_name", "testdata/file.txt")

			_, bodyReader, err := subject.Build()
			Expect(err).NotTo(HaveOccurred())

			body := read(bodyReader)
			Expect(body).To(ContainSubstring("field_name"))
			Expect(body).To(ContainSubstring("file.txt"))
			Expect(body).To(ContainSubstring("FILE_TXT_CONTENTS\n"))
		})

	}) // slurping context

	Context("errors", func() {

		It("should report errors", func() {
			subject.SlurpFile("file1", "testdata/inexisting.txt")
			subject.SlurpFile("file2", "testdata/inexisting.bin")

			contentType, bodyReader, err := subject.Build()
			Expect(err).To(HaveOccurred())
			Expect(contentType).To(BeEmpty())
			Expect(bodyReader).To(BeNil())

			Expect(err.Error()).To(ContainSubstring("multipartbuilder: failed to open file testdata/inexisting.txt for field file1"))
			Expect(err.Error()).To(ContainSubstring("multipartbuilder: failed to open file testdata/inexisting.bin for field file2"))
		})

	}) // errors context

	Context("request", func() {

		It("should build HTTP request", func() {
			subject.WriteField("field_name", "FIELD_VALUE")

			req, err := subject.BuildRequest("POST", "https://test.com/")
			Expect(err).NotTo(HaveOccurred())

			Expect(req.Method).To(Equal("POST"))
			Expect(req.URL.String()).To(Equal("https://test.com/"))
			Expect(req.Header.Get("Content-Type")).To(ContainSubstring("multipart"))

			body := read(req.Body)
			Expect(body).To(ContainSubstring("field_name"))
			Expect(body).To(ContainSubstring("FIELD_VALUE"))
		})

	}) // request context

})

func read(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}
