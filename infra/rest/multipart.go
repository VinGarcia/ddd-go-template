package rest

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
)

type MultipartData map[string]io.Reader

type multipartFile struct {
	io.Reader
	name string
}

type multipartItem struct {
	io.Reader
	contentType string
}

func MultipartFile(data io.Reader, name string) io.Reader {
	return multipartFile{
		Reader: data,
		name:   name,
	}
}

func MultipartItem(data io.Reader, contentType string) io.Reader {
	return multipartItem{
		Reader:      data,
		contentType: contentType,
	}
}

type multipartStream struct {
	formClosed      bool
	multipartWriter *multipart.Writer
	parts           []formPart

	currentPartWriter io.Writer
	currentPartReader io.Reader

	buf *bytes.Buffer
}

type formPart struct {
	reader    io.Reader
	fieldname string
}

func newMultipartStream(data MultipartData) (_ *multipartStream, contentType string, err error) {
	var buffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&buffer)

	// These are used by the `write()` closure
	// to start sending the data only when requested
	stream := multipartStream{
		buf:             &buffer,
		multipartWriter: multipartWriter,
	}

	for key, reader := range data {
		stream.parts = append(stream.parts, formPart{
			reader:    reader,
			fieldname: key,
		})
	}

	return &stream, multipartWriter.FormDataContentType(), nil
}

func (m *multipartStream) Read(p []byte) (n int, err error) {
	if m.buf.Len() == 0 {
		err := m.loadFormBuffer(m.buf, int64(len(p)))
		if err != nil {
			return 0, err
		}
	}

	n, err = m.buf.Read(p)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (m *multipartStream) loadFormBuffer(buf *bytes.Buffer, numBytes int64) (err error) {
	if m.currentPartReader == nil && len(m.parts) == 0 {
		if !m.formClosed {
			// Write the closing of the form before returning io.EOF
			m.multipartWriter.Close()
			m.formClosed = true
			return nil
		}
		return io.EOF
	}

	if m.currentPartReader == nil {
		m.currentPartReader, m.currentPartWriter, err = m.loadNextPart()
		if err != nil {
			return err
		}
	}

	_, err = io.CopyN(m.currentPartWriter, m.currentPartReader, numBytes)
	if err == io.EOF {
		// If this part is finished:
		m.currentPartReader = nil
		m.currentPartWriter = nil
	} else if err != nil {
		return err
	}

	return nil
}

func (m *multipartStream) loadNextPart() (io.Reader, io.Writer, error) {
	p := m.parts[0]
	m.parts = m.parts[1:]
	writer, err := createFormPart(m.multipartWriter, p.reader, p.fieldname)
	return p.reader, writer, err
}

func createFormPart(w *multipart.Writer, reader io.Reader, fieldname string) (io.Writer, error) {
	switch reader := reader.(type) {
	case multipartFile:
		return w.CreateFormFile(fieldname, reader.name)
	case multipartItem:
		header := textproto.MIMEHeader{
			"Content-Disposition": []string{fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(fieldname))},
		}
		if reader.contentType != "" {
			header.Set("Content-Type", reader.contentType)
		}
		return w.CreatePart(header)
	default:
		header := textproto.MIMEHeader{
			"Content-Disposition": []string{fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(fieldname))},
		}
		return w.CreatePart(header)
	}
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
