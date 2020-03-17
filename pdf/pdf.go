package pdf

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
)

const pdfExt = ".pdf"

// Meta represents pdf-file's meta info
type Meta struct {
	Filename string
	Info     pdfInfo
}

type pdfInfo []struct {
	Field string
	Value interface{}
}

// File represents pdf file
type File struct {
	Filename string
	File     io.ReadSeeker
}

// ReadPDFFromRequest reads pdf file from request
func ReadPDFFromRequest(r *http.Request) (*File, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	if path.Ext(handler.Filename) != pdfExt {
		return nil, errors.New(`error: file format is not .pdf`)
	}

	return &File{handler.Filename, io.ReadSeeker(file)}, nil
}

// GetMeta returns pdf file's meta info
func (file *File) GetMeta() (*Meta, error) {
	conf := pdfcpu.NewDefaultConfiguration()

	ctx, err := pdfcpu.Read(file.File, conf)
	if err != nil {
		return nil, err
	}

	if err := validate.XRefTable(ctx.XRefTable); err != nil {
		return nil, err
	}

	var meta = &Meta{
		Filename: file.Filename,
		Info: pdfInfo{
			{"Title", ctx.Title},
			{"Subject", ctx.Subject},
			{"Keywords", ctx.Keywords},
			{"Author", ctx.Author},
			{"Creator", ctx.Creator},
			{"Producer", ctx.Producer},
			{"CreationDate", ctx.CreationDate},
			{"ModDate", ctx.ModDate},
			{"Tagged", ctx.Tagged},
			{"Pages", ctx.PageCount},
			{"EncryptedUsingAES", ctx.EncryptUsingAES},
			{"Encrypted", ctx.Encrypt != nil},
			{"PageSize", getPageSize(ctx)},
			{"FileSize", *ctx.Size},
			{"Optimized", ctx.Optimized},
			{"PDFVersion", ctx.Version().String()},
			{"Hybrid", ctx.Read.Hybrid},
			{"Linearized", ctx.Read.Linearized},
			{"UsingXRefStreams", ctx.Read.UsingXRefStreams},
			{"UsingObjectStreams", ctx.Read.UsingObjectStreams},
			{"Watermarked", ctx.Watermarked},
			{"Permissions", getPermissions(ctx)},
		},
	}

	return meta, nil
}

func getPermissions(ctx *pdfcpu.Context) string {
	perms := pdfcpu.Permissions(ctx)
	return strings.Join(perms, ", ")
}

func getPageSize(ctx *pdfcpu.Context) string {
	var results = make([]string, 0)
	pd, err := ctx.PageDims()
	if err != nil {
		return ""
	}

	m := map[pdfcpu.Dim]bool{}
	for _, d := range pd {
		m[d] = true
	}

	for d := range m {
		var orientation string

		if d.Landscape() {
			orientation = "landscape"
		} else if d.Portrait() {
			orientation = "letter"
		} else {
			orientation = "unknown"
		}

		results = append(results,
			fmt.Sprintf("%.2fx%.2f (%s orientation, %.2f ratio)",
				d.Width,
				d.Height,
				orientation,
				d.AspectRatio(),
			),
		)
	}

	return strings.Join(results, ", ")
}
