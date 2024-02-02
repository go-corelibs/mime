// Copyright (c) 2024  The Go-CoreLibs Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package mime provides mime type system utilities
package mime

import (
	"errors"
	goMime "mime"
	"strings"

	"github.com/gabriel-vasile/mimetype"

	clPath "github.com/go-corelibs/path"
)

const (
	TextMimeType       = "text/plain"
	HtmlMimeType       = "text/html"
	CssMimeType        = "text/css"
	ScssMimeType       = "text/x-scss"
	JsonMimeType       = "application/json"
	JavaScriptMimeType = "text/javascript"
	BinaryMimeType     = "application/octet-stream"

	// DirectoryMimeType defines the mime type used for filesystem directories
	DirectoryMimeType = "inode/directory"
	// EnjinMimeType defines the mime type for the Go-Enjin project's `njn`
	// page format
	EnjinMimeType = "text/enjin"
	// OrgModeMimeType defines the mime type used by the Go-Enjin project to
	// identify the `org-mode` page format
	OrgModeMimeType = "text/org-mode"
	// MarkdownMimeType defines the mime type used by the Go-Enjin project to
	// identify the `markdown` page format
	MarkdownMimeType = "text/markdown"
)

const (
	// EnjinExtension defines the file extension associated with the
	// EnjinMimeType
	EnjinExtension = "njn"
	// OrgModeExtension defines the file extension associated with the
	// OrgModeMimeType
	OrgModeExtension = "org"
	// MarkdownExtension defines the file extension associated with the
	// MarkdownMimeType
	MarkdownExtension = "md"
)

var (
	gExtension = &lookup{m: map[string]string{
		"txt":  TextMimeType + "; charset=utf-8",
		"html": HtmlMimeType + "; charset=utf-8",
		"css":  CssMimeType + "; charset=utf-8",
		"scss": ScssMimeType + "; charset=utf-8",
		"json": JsonMimeType + "; charset=utf-8",
		"js":   JavaScriptMimeType + "; charset=utf-8",
	}}
	gCharset = &lookup{m: map[string]string{
		TextMimeType:       "utf-8",
		HtmlMimeType:       "utf-8",
		CssMimeType:        "utf-8",
		ScssMimeType:       "utf-8",
		JsonMimeType:       "utf-8",
		JavaScriptMimeType: "utf-8",
		EnjinMimeType:      "utf-8",
		OrgModeMimeType:    "utf-8",
		MarkdownMimeType:   "utf-8",
	}}
)

func init() {
	_ = RegisterTextType(EnjinMimeType, EnjinExtension, nil)
	_ = RegisterTextType(OrgModeMimeType, OrgModeExtension, nil)
	_ = RegisterTextType(MarkdownMimeType, MarkdownExtension, nil)
}

// GetExtension returns the mime type internally associated with this package
// using SetExtension, or if not present uses mime.TypeByExtension to lookup
// further
func GetExtension(extension string) (mime string, ok bool) {
	extension = strings.TrimPrefix(extension, ".")
	if mime, ok = gExtension.get(extension); !ok {
		mime = goMime.TypeByExtension("." + extension)
		ok = mime != ""
	}
	return
}

// SetExtension registers the given extension with the given mime type string.
// There can only be one mime type associated per extension and SetExtension
// will overwrite any existing value. If `mime` is empty, any internal
// association with the extension is cleared
func SetExtension(extension, mime string) {
	extension = strings.TrimPrefix(extension, ".")
	if mime == "" {
		gExtension.unset(extension)
		return
	}
	gExtension.set(extension, mime)
}

// GetCharset returns the `charset` internally associated with this package
func GetCharset(mime string) (charset string, ok bool) {
	charset, ok = gCharset.get(PruneCharset(mime))
	return
}

// SetCharset registers the given extension with the given charset string.
// There can only be one charset associated per extension and SetCharset
// will overwrite any existing value
func SetCharset(mime, charset string) {
	mime = PruneCharset(mime)
	if charset == "" {
		gCharset.unset(mime)
		return
	}
	gCharset.set(mime, charset)
}

// PruneCharset uses mime.ParseMediaType to parse the given mime string and
// returns only the media type value
func PruneCharset(mime string) (pruned string) {
	pruned, _, _ = goMime.ParseMediaType(mime)
	return
}

// RegisterTextType associates the given `mime` with the given `extension` and
// if the `detector` is not nil, registers the given `mime` with TextMimeType
// as it's parent within the github.com/gabriel-vasile/mimetype system
func RegisterTextType(mime, extension string, detector func(raw []byte, limit uint32) bool) (err error) {
	extension = strings.TrimPrefix(extension, ".")
	var mediatype string
	if mime == "" || extension == "" {
		err = errors.New("mime and extension arguments must not be empty")
		return
	} else if m, p, e := goMime.ParseMediaType(mime); e != nil {
		err = e
		return
	} else {
		p["charset"] = "utf-8"
		mime = goMime.FormatMediaType(m, p)
		mediatype = m
	}
	SetExtension(extension, mime)
	SetCharset(mediatype, "utf-8")
	for _, m := range []string{mediatype, mime} {
		if detector != nil {
			mimetype.Lookup(TextMimeType).Extend(detector, m, "."+extension)
		} else {
			mimetype.Lookup(TextMimeType).Extend(PlainTextDetector, m, "."+extension)
		}
	}
	err = goMime.AddExtensionType("."+extension, mediatype)
	return
}

// IsPlainText returns true if the given `mime` is of `text/plain` type.
// IsPlainText checks if it has an internally registered charset first
// and if so, returns true early. If the `mime` is not internally registered
// with a charset (via SetCharset), IsPlainText uses
// github.com/gabriel-vasile/mimetype.Lookup to check if the given `mime`
// is exactly TextMimeType or if any of the found mime type's parents are
// TextMimeType
func IsPlainText(mime string) (yes bool) {
	mime = PruneCharset(mime)
	if _, yes = GetCharset(mime); yes {
		return
	}
	if mt := mimetype.Lookup(mime); mt != nil {
		if yes = mt.Is(TextMimeType); yes {
			return
		}
		for check := mt; check != nil; check = check.Parent() {
			if yes = check.Is(TextMimeType); yes {
				return
			}
		}
	}
	return
}

// FromPathOnly checks the given `path` for any extensions using
// github.com/go-corelibs/path.ExtExt and uses GetExtension with any
// extensions found
func FromPathOnly(path string) (mime string) {
	if path != "" {
		if a, b := clPath.ExtExt(path); b != "" && a == "tmpl" {
			mime, _ = GetExtension(b)
		} else if a != "" {
			mime, _ = GetExtension(a)
		}
	}
	return
}

// Mime returns the MIME type string of a local filesystem directory or file.
// The specific type returned for directories is defined by the
// DirectoryMimeType constant
func Mime(path string) (mime string) {
	if clPath.IsDir(path) {
		mime = DirectoryMimeType
		return
	} else if clPath.IsFile(path) {
		if mime = FromPathOnly(path); mime != "" {
			return
		} else if mt, err := mimetype.DetectFile(path); err == nil {
			mime = mt.String()
		}
	}
	return
}

// PlainTextDetector is the default detector used when RegisterTextType is
// given a `nil` value for it's `detector` argument. PlainTextDetector always
// returns true
func PlainTextDetector(raw []byte, limit uint32) bool {
	// TODO: figure out a better way of detecting plain text things, for example: Mime("./LICENSE") returns "text/markdown" when it should be "text/plain"
	return true
}
