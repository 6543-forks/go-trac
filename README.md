# go-trac #

go-trac is a Go client library for accessing Trac data via [Trac XML-RPC Plugin](https://trac-hacks.org/wiki/XmlRpcPlugin)

## Usage ##

```go
import "github.com/ics/go-trac/pkg/trac"

tracURL := "https://user:pass@trac.example.com/login/jsonrpc"

trc := trac.NewClient(tracURL, nil)
t, err := trc.Ticket.Get(1)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%#v\n", t)
```

**Documentation:** [![GoDoc](https://godoc.org/github.com/ics/go-trac/github?status.svg)](https://godoc.org/github.com/ics/go-trac/pkg/trac)
