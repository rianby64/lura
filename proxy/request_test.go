// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRequestGeneratePath(t *testing.T) {
	r := Request{
		Method: "GET",
		Params: map[string]string{
			"Supu": "42",
			"Tupu": "false",
			"Foo":  "bar",
			"Wild": "/level1/level2",
			"Wrng": "level1/level2",
			"Bad":  "",
			"2Bad": "/",
			"3Bad": "/",
			"4Bad": "///bad4",
			"Lvl2": "/level/2",
			"Lvl3": "/lvl/3",
		},
	}

	for i, testCase := range [][]string{
		{"/base/{{.4Bad}}?b={{.Foo}}", "/base/bad4?b=bar"},
		{"/a/{{.Supu}}", "/a/42"},
		{"/a?b={{.Tupu}}", "/a?b=false"},
		{"/a/{{.Supu}}/foo/{{.Foo}}", "/a/42/foo/bar"},
		{"/a", "/a"},
		{"/base/{{.Wild}}?b={{.Foo}}", "/base/level1/level2?b=bar"},
		{"/base/{{.Wrng}}?b={{.Foo}}", "/base/level1/level2?b=bar"},
		{"/base/{{.Bad}}?b={{.Foo}}", "/base/?b=bar"},
		{"/base/{{.2Bad}}?b={{.Foo}}", "/base/?b=bar"},
		{"/base/{{.2Bad}}/{{.3Bad}}?b={{.Foo}}", "/base//?b=bar"},
		{"/base/{{.Lvl2}}/{{.Lvl3}}?b={{.Foo}}", "/base/level/2/lvl/3?b=bar"},
	} {
		r.GeneratePath(testCase[0])
		if r.Path != testCase[1] {
			t.Errorf("%d: want %s, have %s", i, testCase[1], r.Path)
		}
	}
}

func TestRequest_Clone(t *testing.T) {
	r := Request{
		Method: "GET",
		Params: map[string]string{
			"Supu": "42",
			"Tupu": "false",
			"Foo":  "bar",
		},
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}
	clone := r.Clone()

	if len(r.Params) != len(clone.Params) {
		t.Errorf("wrong num of params. have: %d, want: %d", len(clone.Params), len(r.Params))
		return
	}
	for k, v := range r.Params {
		if res, ok := clone.Params[k]; !ok {
			t.Errorf("param %s not cloned", k)
		} else if res != v {
			t.Errorf("unexpected param %s. have: %s, want: %s", k, res, v)
		}
	}

	if len(r.Headers) != len(clone.Headers) {
		t.Errorf("wrong num of headers. have: %d, want: %d", len(clone.Headers), len(r.Headers))
		return
	}

	for k, vs := range r.Headers {
		if res, ok := clone.Headers[k]; !ok {
			t.Errorf("header %s not cloned", k)
		} else if len(res) != len(vs) {
			t.Errorf("unexpected header %s. have: %v, want: %v", k, res, vs)
		}
	}

	r.Headers["extra"] = []string{"supu"}

	if len(r.Headers) != len(clone.Headers) {
		t.Errorf("wrong num of headers. have: %d, want: %d", len(clone.Headers), len(r.Headers))
		return
	}

	for k, vs := range r.Headers {
		if res, ok := clone.Headers[k]; !ok {
			t.Errorf("header %s not cloned", k)
		} else if len(res) != len(vs) {
			t.Errorf("unexpected header %s. have: %v, want: %v", k, res, vs)
		}
	}
}

func TestCloneRequest(t *testing.T) {
	body := `{"a":1,"b":2}`
	r := Request{
		Method: "POST",
		Params: map[string]string{
			"Supu": "42",
			"Tupu": "false",
			"Foo":  "bar",
		},
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
	clone := CloneRequest(&r)

	if len(r.Params) != len(clone.Params) {
		t.Errorf("wrong num of params. have: %d, want: %d", len(clone.Params), len(r.Params))
		return
	}
	for k, v := range r.Params {
		if res, ok := clone.Params[k]; !ok {
			t.Errorf("param %s not cloned", k)
		} else if res != v {
			t.Errorf("unexpected param %s. have: %s, want: %s", k, res, v)
		}
	}

	if len(r.Headers) != len(clone.Headers) {
		t.Errorf("wrong num of headers. have: %d, want: %d", len(clone.Headers), len(r.Headers))
		return
	}

	for k, vs := range r.Headers {
		if res, ok := clone.Headers[k]; !ok {
			t.Errorf("header %s not cloned", k)
		} else if len(res) != len(vs) {
			t.Errorf("unexpected header %s. have: %v, want: %v", k, res, vs)
		}
	}

	r.Headers["extra"] = []string{"supu"}

	if _, ok := clone.Headers["extra"]; ok {
		t.Error("the cloned instance shares its headers with the original one")
	}

	delete(r.Params, "Supu")

	if _, ok := clone.Params["Supu"]; !ok {
		t.Error("the cloned instance shares its params with the original one")
	}

	rb, _ := io.ReadAll(r.Body)
	cb, _ := io.ReadAll(clone.Body)

	if !bytes.Equal(cb, rb) || body != string(rb) {
		t.Errorf("unexpected bodies. original: %s, returned: %s", string(rb), string(cb))
	}
}
