package trac

import (
	"encoding/json"
	"fmt"
	"time"
)

// CustomType is used for non-standard types.
// JSON-RPC has no formalized type system, so a class-hint system is used for
// input and output of non-standard types:
//
// {"__jsonclass__": ["datetime", "YYYY-MM-DDTHH:MM:SS"]} => DateTime (UTC)
//
// {"__jsonclass__": ["binary", "<base64-encoded>"]} => Binary
type CustomType struct {
	Kv [2]string `json:"__jsonclass__"`
}

// PageInfo represents page information.
type PageInfo struct {
	Name         string
	Author       string
	Version      int
	LastModified time.Time
	Comment      string
}

// UnmarshalJSON deserializes PageInfo.
func (pi *PageInfo) UnmarshalJSON(in []byte) error {
	type Alias PageInfo
	tmp := struct {
		*Alias
		LastModified CustomType
	}{
		Alias: (*Alias)(pi),
	}
	if err := json.Unmarshal(in, &tmp); err != nil {
		return err
	}
	lm, err := time.Parse(timeFormat, tmp.LastModified.Kv[1])
	if err != nil {
		return err
	}
	pi.LastModified = lm
	return nil
}

// Page represents a Wiki page.
type Page struct {
	Info PageInfo
	Wiki string
	HTML string
}

// Wiki represents WikiRPC.
type Wiki struct {
	client *Client
}

// Page returns the latest version of the Wiki page; both raw text and HTML.
func (w *Wiki) Page(pagename string) (Page, error) {
	var p = Page{}
	pg, err := w.client.Query("wiki.getPage", pagename)
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal(pg.Result, &p.Wiki); err != nil {
		return p, err
	}

	h, err := w.client.Query("wiki.getPageHTML", pagename)
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal(h.Result, &p.HTML); err != nil {
		return p, err
	}

	info, err := w.PageInfo(pagename)
	if err != nil {
		return p, err
	}
	p.Info = info

	return p, nil
}

// PageInfo returns information about the given page.
func (w *Wiki) PageInfo(pagename string) (PageInfo, error) {
	var pi = PageInfo{}
	r, err := w.client.Query("wiki.getPageInfo", pagename)
	if err != nil {
		return pi, err
	}
	if err := json.Unmarshal(r.Result, &pi); err != nil {
		return pi, err
	}
	return pi, nil
}

// RPCVersion returns the version of the Trac API.
func (w *Wiki) RPCVersion() (int, error) {
	var ver int
	r, err := w.client.Query("wiki.getRPCVersionSupported")
	if err != nil {
		return ver, err
	}
	if err := json.Unmarshal(r.Result, &ver); err != nil {
		return ver, err
	}
	return ver, nil
}

// PageVersion is not implemented.
func (w *Wiki) PageVersion(pagename string, version int) error {
	return fmt.Errorf("Not implemented")
}

// RecentChanges is not implemented.
func (w *Wiki) RecentChanges(since time.Time) error {
	return fmt.Errorf("Not implemented")
}

// Pages returns a list of all pages. The result is an array of utf8 pagenames.
func (w *Wiki) Pages() ([]string, error) {
	return w.client.All("wiki.getAllPages")
}

// PageInfoVersion is not implemented.
func (w *Wiki) PageInfoVersion(pagename string) ([]string, error) {
	return nil, fmt.Errorf("Not implemented")
}
