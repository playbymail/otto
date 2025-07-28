// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the `info` command.
package cli

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/maloquacious/wxx/models"
	"github.com/maloquacious/wxx/xmlio"
	"github.com/playbymail/otto/config"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strings"
)

var Command = &cobra.Command{
	Use:   "info",
	Short: "Show map information",
	Long:  `Info displays metadata from a map like  the Worldographer version, height, and width.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			fmt.Printf("info: %q\n", arg)
			if !strings.HasSuffix(arg, ".wxx") {
				fmt.Printf("\tnot a '.wxx' file\n")
				continue
			}
			sb, err := os.Stat(arg)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("\tdoes not exist\n")
				} else {
					fmt.Printf("\tunable to stat\n")
				}
				continue
			} else if sb.IsDir() {
				fmt.Printf("\tis a folder\n")
			} else if !sb.Mode().IsRegular() {
				fmt.Printf("\tis not a file\n")
			}
			fmt.Printf("\t%8d bytes on disk\n", sb.Size())
			input, err := os.ReadFile(arg)
			if err != nil {
				fmt.Printf("\tfailed to read\n")
				continue
			}

			// should be a gzip file
			input, err = unzip(input)
			if err != nil {
				fmt.Printf("\tnot gzip compressed\n")
			}
			fmt.Printf("\t%8d bytes compressed\n", sb.Size())
			fmt.Printf("\t%8d bytes uncompressed\n", len(input))

			// should be UTF-16/BE
			if len(input)%2 != 0 {
				fmt.Printf("\tnot utf-16/be encoded\n")
			}
			// verify the BOM
			if bytes.HasPrefix(input, []byte{0xfe, 0xff}) {
				fmt.Printf("\t%8d bytes utf-16/be encoded\n", len(input))
			} else if bytes.HasPrefix(input, []byte{0xff, 0xfe}) {
				fmt.Printf("\t%8d bytes utf-16/le encoded\n", len(input))
				continue
			} else {
				fmt.Printf("\tnot utf-16/be encoded\n")
				continue
			}

			// convert to UTF-8
			utf16Encoding := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
			input, err = io.ReadAll(transform.NewReader(bytes.NewReader(input), utf16Encoding.NewDecoder()))
			fmt.Printf("\t%8d bytes utf-8     encoded\n", len(input))

			// verify the xml header. the encoding may be wrong, but we'll accept it.
			xmlHeaderIndex, xmlHeaders := -1, []struct {
				heading  string
				version  string
				encoding string
			}{
				{heading: "<?xml version='1.0' encoding='utf-8'?>\n", version: "1.0", encoding: "utf-8"},
				{heading: "<?xml version='1.0' encoding='utf-16'?>\n", version: "1.0", encoding: "utf-16"},
				{heading: "<?xml version='1.1' encoding='utf-8'?>\n", version: "1.1", encoding: "utf-8"},
				{heading: "<?xml version='1.1' encoding='utf-16'?>\n", version: "1.1", encoding: "utf-16"},
			}
			for i, header := range xmlHeaders {
				if bytes.HasPrefix(input, []byte(header.heading)) {
					xmlHeaderIndex = i
					break
				}
			}
			if xmlHeaderIndex == -1 {
				fmt.Printf("\tmissing xml header\n")
				continue
			}
			fmt.Printf("\t%8s xml version\n", xmlHeaders[xmlHeaderIndex].version)
			fmt.Printf("\t%8s xml encoding\n", xmlHeaders[xmlHeaderIndex].encoding)
			fmt.Printf("\t%8d bytes xml data\n", len(input))

			// skip past the xml header so that we will be able to unmarshal the input
			input = input[len(xmlHeaders[xmlHeaderIndex].heading):]
			if !bytes.HasPrefix(input, []byte("<map ")) {
				fmt.Printf("\tmissing <map> element\n")
				continue
			}

			// read the map metadata
			xmlMetaData, err := readMapMetadata(input)
			if err != nil {
				fmt.Printf("\t%v\n", err)
				continue
			}
			if xmlMetaData.Release == "" && xmlMetaData.Version != "" && xmlMetaData.Schema == "" {
				// H2017 file
				fmt.Printf("\t%8s worldographer version\n", "H2017")
				fmt.Printf("\t%8s version\n", xmlMetaData.Version)
			} else if xmlMetaData.Release == "2025" && xmlMetaData.Version != "" && xmlMetaData.Schema != "" {
				// W2025 file
				fmt.Printf("\t%8s worldographer version\n", "W2025")
				fmt.Printf("\t%8s version\n", xmlMetaData.Version)
				fmt.Printf("\t%8s schema\n", xmlMetaData.Schema)
			} else {
				fmt.Printf("\tunknown metadata: %q %q %q\n", xmlMetaData.Release, xmlMetaData.Version, xmlMetaData.Schema)
				continue
			}

			_, err = xmlio.Read(input)
			if err != nil {
				fmt.Printf("\t%v\n", err)
				continue
			}
		}
		return nil
	},
}

func RegisterArgs(cfg *config.Config_t) error {
	return nil
}

type mapMetaData struct {
	Version string `xml:"version,attr"` // required
	Release string `xml:"release,attr"` // H2017 optional, W2025 required
	Schema  string `xml:"schema,attr"`  // H2017 optional, W2025 required
}

// readMapMetadata
func readMapMetadata(input []byte) (mapMetaData, error) {
	// sanity check, sweet sanity checks
	if !bytes.HasPrefix(input, []byte(`<map `)) {
		return mapMetaData{}, fmt.Errorf("<map> element missing")
	}
	// speed up the remaining sanity checks by extracting the map attributes.
	// we have to make the map element self-closing for this to work.
	endOfMap := bytes.IndexByte(input, '>')
	if endOfMap == -1 {
		return mapMetaData{}, fmt.Errorf("<map> not closed")
	}
	// initialize metadata with a copy of the source up to (but not including) the first closing '>'
	metadata := append(make([]byte, 0, endOfMap+1), input[:endOfMap]...)
	metadata = append(metadata, '/', '>')
	// read the version from the xml data
	var results mapMetaData
	if err := xml.Unmarshal(metadata, &results); err != nil {
		return mapMetaData{}, errors.Join(models.ErrInvalidMapMetadata, err)
	}
	return results, nil
}

func unzip(input []byte) ([]byte, error) {
	// create a new gzip reader to process the source
	gzr, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	defer func(gzr *gzip.Reader) {
		_ = gzr.Close() // ignore errors
	}(gzr)
	return io.ReadAll(gzr)
}
