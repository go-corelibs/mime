// Copyright (c) 2024  The Go-Enjin Authors
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

package mime

import (
	"testing"

	"github.com/gabriel-vasile/mimetype"
	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("GetExtension", t, func() {
		Convey("additional mime types", func() {
			Convey("njn", func() {
				mime, ok := GetExtension("njn")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/enjin; charset=utf-8")
			})
			Convey("org", func() {
				mime, ok := GetExtension("org")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/org-mode; charset=utf-8")
			})
			Convey("md", func() {
				mime, ok := GetExtension("md")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/markdown; charset=utf-8")
			})
		})

		Convey("default mime types", func() {
			Convey("txt", func() {
				mime, ok := GetExtension("txt")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/plain; charset=utf-8")
			})
			Convey("html", func() {
				mime, ok := GetExtension("html")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/html; charset=utf-8")
			})
			Convey("css", func() {
				mime, ok := GetExtension("css")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/css; charset=utf-8")
			})
			Convey("scss", func() {
				mime, ok := GetExtension("scss")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/x-scss; charset=utf-8")
			})
			Convey("json", func() {
				mime, ok := GetExtension("json")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "application/json; charset=utf-8")
			})
			Convey("js", func() {
				mime, ok := GetExtension("js")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/javascript; charset=utf-8")
			})
		})

		Convey("os mime types", func() {
			Convey("zip", func() {
				mime, ok := GetExtension("zip")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "application/zip")
			})
			Convey("txt", func() {
				mime, ok := GetExtension("txt")
				So(ok, ShouldBeTrue)
				So(mime, ShouldEqual, "text/plain; charset=utf-8")
			})
		})
	})

	Convey("SetExtension", t, func() {
		mime, ok := GetExtension("not-a-thing")
		So(ok, ShouldBeFalse)
		So(mime, ShouldBeEmpty)
		SetExtension("not-a-thing", "application/nope")
		mime, ok = GetExtension("not-a-thing")
		So(ok, ShouldBeTrue)
		So(mime, ShouldEqual, "application/nope")
		SetExtension("not-a-thing", "")
		mime, ok = GetExtension("not-a-thing")
		So(ok, ShouldBeFalse)
		So(mime, ShouldBeEmpty)
	})

	Convey("GetCharset", t, func() {
		charset, ok := GetCharset("text/enjin")
		So(ok, ShouldBeTrue)
		So(charset, ShouldEqual, "utf-8")
		charset, ok = GetCharset("text/nope")
		So(ok, ShouldBeFalse)
		So(charset, ShouldBeEmpty)
	})

	Convey("SetCharset", t, func() {
		charset, ok := GetCharset("nope/nope")
		So(ok, ShouldBeFalse)
		So(charset, ShouldBeEmpty)
		SetCharset("nope/nope", "utf-8")
		charset, ok = GetCharset("nope/nope")
		So(ok, ShouldBeTrue)
		So(charset, ShouldEqual, "utf-8")
		SetCharset("nope/nope", "")
		charset, ok = GetCharset("nope/nope")
		So(ok, ShouldBeFalse)
		So(charset, ShouldBeEmpty)
	})

	Convey("PruneCharset", t, func() {
		So(PruneCharset(""), ShouldBeEmpty)
		So(PruneCharset("nope/plain"), ShouldEqual, "nope/plain")
		So(PruneCharset("nope/plain; charset=utf-8"), ShouldEqual, "nope/plain")
	})

	Convey("RegisterTextType", t, func() {
		So(RegisterTextType("", "", nil), ShouldNotBeNil)
		So(RegisterTextType("not/a-thing", "nope", nil), ShouldBeNil)
		var ok bool
		var charset, mime string
		charset, ok = GetCharset("not/a-thing")
		So(ok, ShouldBeTrue)
		So(charset, ShouldEqual, "utf-8")
		mime, ok = GetExtension("nope")
		So(ok, ShouldBeTrue)
		So(mime, ShouldEqual, "not/a-thing; charset=utf-8")
		mt := mimetype.Lookup("not/a-thing; charset=utf-8")
		So(mt, ShouldNotBeNil)
		pt := mt.Parent()
		So(pt, ShouldNotBeNil)
		So(pt.Is("text/plain"), ShouldBeTrue)
		mt = mimetype.Lookup("not/a-thing")
		So(mt, ShouldNotBeNil)
		pt = mt.Parent()
		So(pt, ShouldNotBeNil)
		So(pt.Is("text/plain"), ShouldBeTrue)
		So(RegisterTextType("bad mime", "bad", nil), ShouldNotBeNil)
		So(RegisterTextType("good/mime", "good", func(raw []byte, limit uint32) bool {
			return false
		}), ShouldBeNil)
	})

	Convey("IsPlainText", t, func() {
		Convey("internally registered types", func() {
			So(IsPlainText("text/enjin"), ShouldBeTrue)
			So(IsPlainText("text/org-mode"), ShouldBeTrue)
			So(IsPlainText("text/markdown"), ShouldBeTrue)
		})

		Convey("externally registered types", func() {
			So(IsPlainText("application/zip"), ShouldBeFalse)
			// unset text/plain so that the first actual mimetype lookup
			// check can pass
			SetCharset("text/plain", "")
			So(IsPlainText("text/plain"), ShouldBeTrue)
			// unset text/html so that the parent mimetype lookup
			// check can pass
			SetCharset("text/html", "")
			So(IsPlainText("text/html"), ShouldBeTrue)
		})
	})

	Convey("FromPathOnly", t, func() {
		So(FromPathOnly("file.txt"), ShouldEqual, "text/plain; charset=utf-8")
		So(FromPathOnly("file.html.tmpl"), ShouldEqual, "text/html; charset=utf-8")
	})

	Convey("Mime", t, func() {
		So(Mime("."), ShouldEqual, "inode/directory")
		So(Mime("./testdata/README.md"), ShouldEqual, "text/markdown; charset=utf-8")
		So(Mime("./testdata/empty-png"), ShouldEqual, "image/png")
	})

	Convey("PlainTextDetector", t, func() {
		So(PlainTextDetector([]byte("plain text"), 1024), ShouldBeTrue)
	})

}
