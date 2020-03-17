package pdf

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var testMetas = []struct {
	ExpectErr bool
	Meta      Meta
}{
	{
		ExpectErr: false,
		Meta: Meta{
			Filename: "cpp.pdf",
			Info: pdfInfo{
				{"Title", "C++ International Standard"},
				{"Subject", "C++ International Standard"},
				{"Keywords", ""},
				{"Author", ""},
				{"Creator", "Stefanus Du Toit"},
				{"Producer", "pdfTeX-1.40.13"},
				{"CreationDate", "D:20130515171744-04'00'"},
				{"ModDate", "D:20130515180619-07'00'"},
				{"Tagged", false},
				{"Pages", 1374},
				{"EncryptedUsingAES", true},
				{"Encrypted", false},
				{"PageSize", "612.00x792.00 (letter orientation, 0.77 ratio)"},
				{"FileSize", 65046},
				{"Optimized", false},
				{"PDFVersion", "1.5"},
				{"Hybrid", false},
				{"Linearized", true},
				{"UsingXRefStreams", true},
				{"UsingObjectStreams", true},
				{"Watermarked", false},
				{"Permissions", "Full access"},
			},
		},
	},
	{
		ExpectErr: false,
		Meta: Meta{
			Filename: "org.pdf",
			Info: pdfInfo{
				{"Title", ""},
				{"Subject", ""},
				{"Keywords", ""},
				{"Author", ""},
				{"Creator", "TeX"},
				{"Producer", "pdfTeX-1.40.19"},
				{"CreationDate", "D:20200213080504Z"},
				{"ModDate", "D:20200213080504Z"},
				{"Tagged", false},
				{"Pages", 305},
				{"EncryptedUsingAES", true},
				{"Encrypted", false},
				{"PageSize", "612.00x792.00 (letter orientation, 0.77 ratio)"},
				{"FileSize", 6204},
				{"Optimized", false},
				{"PDFVersion", "1.5"},
				{"Hybrid", false},
				{"Linearized", false},
				{"UsingXRefStreams", true},
				{"UsingObjectStreams", true},
				{"Watermarked", false},
				{"Permissions", "Full access"},
			},
		},
	},
	{
		ExpectErr: false,
		Meta: Meta{
			Filename: "about_metadata.pdf",
			Info: pdfInfo{
				{"Title", "About Metadata"},
				{"Subject", "By simple definition, metadata is data about data. Metadata is structured information that explains, describes, or locates the original primary data, or that otherwise makes using the original primary data more efficient. A wide variety of industries use metadata, but for the purposes of digital imaging, there are currently only a few technical structures or schema that are being employed. A schema is a set of properties and their defined meanings, such as the type of value (date, size, URL, or any useful designation). \r"},
				{"Keywords", ""},
				{"Author", "Adobe Systems Incorporated"},
				{"Creator", "Adobe InDesign CS (3.0.1)"},
				{"Producer", "Adobe PDF Library 6.0"},
				{"CreationDate", "D:20040924145630Z"},
				{"ModDate", "D:20040929132118-07'00'"},
				{"Tagged", true},
				{"Pages", 7},
				{"EncryptedUsingAES", true},
				{"Encrypted", false},
				{"PageSize", "612.00x792.00 (letter orientation, 0.77 ratio)"},
				{"FileSize", 343},
				{"Optimized", false},
				{"PDFVersion", "1.4"},
				{"Hybrid", false},
				{"Linearized", true},
				{"UsingXRefStreams", false},
				{"UsingObjectStreams", false},
				{"Watermarked", false},
				{"Permissions", "Full access"},
			},
		},
	},
	{
		ExpectErr: true,
		Meta: Meta{
			Filename: "bad.pdf",
		},
	},
}

func TestGetMeta(t *testing.T) {
	for _, expectedMeta := range testMetas {
		t.Run(expectedMeta.Meta.Filename, func(t *testing.T) {
			file, err := os.Open(filepath.Join("pdf_samples", expectedMeta.Meta.Filename))
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			pdfFile := File{
				Filename: expectedMeta.Meta.Filename,
				File:     file,
			}

			actualMeta, err := pdfFile.GetMeta()
			if (err != nil) != expectedMeta.ExpectErr {
				t.Fatalf(`expected err = %v, but got err = %v`, expectedMeta.ExpectErr, err != nil)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(actualMeta.Info, expectedMeta.Meta.Info) {
				t.Fatalf(`expected != actual for file: %s`, expectedMeta.Meta.Filename)
			}
		})
	}
}
