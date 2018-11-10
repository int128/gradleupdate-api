// This file is automatically generated by qtc from "base.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line base.qtpl:1
package templates

//line base.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line base.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line base.qtpl:1
func StreamHeader(qw422016 *qt422016.Writer) {
	//line base.qtpl:1
	qw422016.N().S(`
  <link href="https://maxcdn.bootstrapcdn.com/bootswatch/3.3.7/cosmo/bootstrap.min.css"
        rel="stylesheet"
        integrity="sha384-h21C2fcDk/eFsW9sC9h0dhokq5pDinLNklTKoxIZRUn3+hvmgQSffLLQ4G4l2eEr"
        crossorigin="anonymous"/>
  <link href="/static/app.css" rel="stylesheet"/>
`)
//line base.qtpl:7
}

//line base.qtpl:7
func WriteHeader(qq422016 qtio422016.Writer) {
	//line base.qtpl:7
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line base.qtpl:7
	StreamHeader(qw422016)
	//line base.qtpl:7
	qt422016.ReleaseWriter(qw422016)
//line base.qtpl:7
}

//line base.qtpl:7
func Header() string {
	//line base.qtpl:7
	qb422016 := qt422016.AcquireByteBuffer()
	//line base.qtpl:7
	WriteHeader(qb422016)
	//line base.qtpl:7
	qs422016 := string(qb422016.B)
	//line base.qtpl:7
	qt422016.ReleaseByteBuffer(qb422016)
	//line base.qtpl:7
	return qs422016
//line base.qtpl:7
}
