package multipartbuilder_test

import (
	"io"
	"mime"
	"mime/multipart"
	"strings"

	. "github.com/mxmCherry/multipartbuilder"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Builder", func() {
	var subject *Builder

	BeforeEach(func() {
		subject = New()
	})

	Context("fields", func() {

		It("should add single field", func() {
			subject.AddField("field_name", "FIELD_VALUE")

			form := readForm(subject.Build())

			Expect(form.Value).To(Equal(map[string][]string{
				"field_name": {"FIELD_VALUE"},
			}))
		})

	}) // fields context

	Context("readers", func() {

		It("should add reader", func() {
			subject.AddReader("field_name", "file.bin", strings.NewReader("READER_CONTENTS"))

			form := readForm(subject.Build())

			Expect(form.File).To(HaveKey("field_name"))
			Expect(form.File["field_name"]).To(HaveLen(1))

			fileHeader := form.File["field_name"][0]
			Expect(fileHeader).NotTo(BeNil())
			Expect(fileHeader.Filename).To(Equal("file.bin"))

			Expect(read(fileHeader)).To(Equal("READER_CONTENTS"))
		})

	}) // readers context

	Context("files", func() {

		It("should add file", func() {
			subject.AddFile("field_name", "testdata/file.txt")

			form := readForm(subject.Build())

			Expect(form.File).To(HaveKey("field_name"))
			Expect(form.File["field_name"]).To(HaveLen(1))

			fileHeader := form.File["field_name"][0]
			Expect(fileHeader).NotTo(BeNil())
			Expect(fileHeader.Filename).To(Equal("file.txt"))

			Expect(read(fileHeader)).To(Equal("FILE_TXT_CONTENTS\n"))
		})

	}) // files context

	Context("errors", func() {

		It("should report errors", func() {
			subject.AddFile("file1", "testdata/inexisting.txt")

			_, reader := subject.Build()

			_, err := io.ReadAll(reader)
			Expect(err.Error()).To(ContainSubstring("multipartbuilder: failed to open file file1 (testdata/inexisting.txt)"))
		})

	}) // errors context

})

func readForm(ctype string, r io.Reader) *multipart.Form {
	const maxMemory = 1024 * 1024

	mediaType, params, err := mime.ParseMediaType(ctype)
	Expect(err).NotTo(HaveOccurred())
	Expect(strings.HasPrefix(mediaType, "multipart/"))

	mr := multipart.NewReader(r, params["boundary"])
	form, err := mr.ReadForm(maxMemory)
	Expect(err).NotTo(HaveOccurred())
	return form
}

func read(fh *multipart.FileHeader) string {
	r, err := fh.Open()
	Expect(err).NotTo(HaveOccurred())

	b, err := io.ReadAll(r)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}
