package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	cfg "github.com/slotopol/balance/config"
)

type AjaxErr struct {
	What string `json:"what" yaml:"what" xml:"what"`
	Code int    `json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty"`
	UID  uint64 `json:"uid,omitempty" yaml:"uid,omitempty" xml:"uid,omitempty,attr"`
}

func (err AjaxErr) Error() string {
	return fmt.Sprintf("what: %s, code: %d", err.What, err.Code)
}

var (
	ErrAlienServer = errors.New("server is not slotopol server")
)

func HttpGet[Tr any](path string, token string, arg url.Values) (ret Tr, status int, err error) {
	defer func() {
		if err != nil {
			log.Printf("error on api call '%s': %s", path, err.Error())
		}
	}()
	var fullurl = cfg.Credentials.Addr + path
	if arg != nil {
		var u, _ = url.ParseRequestURI(fullurl)
		var q = u.Query()
		for k, v := range arg {
			q[k] = v
		}
		fullurl = u.String()
	}
	var req *http.Request
	if req, err = http.NewRequest("GET", fullurl, nil); err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	return DoReq[Tr](req)
}

func HttpPost[Ta, Tr any](path string, token string, arg *Ta) (ret Tr, status int, err error) {
	defer func() {
		if err != nil {
			log.Printf("error on api call '%s': %s", path, err.Error())
		}
	}()
	var b []byte
	if arg != nil {
		if b, err = json.Marshal(arg); err != nil {
			return
		}
	}
	var req *http.Request
	if req, err = http.NewRequest("POST", cfg.Credentials.Addr+path, bytes.NewReader(b)); err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	return DoReq[Tr](req)
}

func DoReq[Tr any](req *http.Request) (ret Tr, status int, err error) {
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if vers := resp.Header.Get("Server"); !strings.HasPrefix(vers, "slotopol/") {
		err = ErrAlienServer
		return
	}

	status = resp.StatusCode

	var b []byte
	if b, err = io.ReadAll(resp.Body); err != nil {
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		var aerr AjaxErr
		if err = json.Unmarshal(b, &aerr); err != nil {
			return
		}
		err = &aerr
		return
	}
	if resp.StatusCode != http.StatusOK {
		return
	}
	if err = json.Unmarshal(b, &ret); err != nil {
		return
	}
	return
}
