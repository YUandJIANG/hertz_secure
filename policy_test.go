// The MIT License (MIT)
//
// Copyright (c) 2016 Bo-Yi Wu
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package secure

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

const (
	testResponse = "bar"
)

func newServer(cfg Config) *route.Engine {
	opts := config.NewOptions([]config.Option{})
	engine := route.NewEngine(opts)
	engine.Use(New(cfg))
	engine.GET("/foo", func(_ context.Context, c *app.RequestContext) {
		c.String(200, testResponse)
	})
	return engine
}

func TestNoConfig(t *testing.T) {
	engine := newServer(Config{})
	w := ut.PerformRequest(engine, "GET", "http://example.com/foo", nil)
	result := w.Result()
	assert.DeepEqual(t, 200, result.StatusCode())
	assert.DeepEqual(t, "bar", string(result.Body()))
}

func TestDefaultConfig(t *testing.T) {
	engine := newServer(DefaultConfig())
	w := ut.PerformRequest(engine, "GET", "https://www.example.com/foo", nil)
	result := w.Result()
	assert.Assert(t, http.StatusOK == result.StatusCode())
	assert.DeepEqual(t, "bar", string(result.Body()))
	res := ut.PerformRequest(engine, "Get", "http://www.example.com/foo", nil).Result()

	assert.Assert(t, http.StatusMovedPermanently == res.StatusCode())
	assert.DeepEqual(t, "https://www.example.com/foo", res.Header.Get("Location"))
}

func TestNoAllowHosts(t *testing.T) {
	engine := newServer(Config{
		AllowedHosts: []string{},
	})
	result := performRequest(engine, "http://www.example.com/foo")
	assert.Assert(t, http.StatusOK == result.StatusCode())
	assert.DeepEqual(t, "bar", string(result.Body()))
}

func TestGoodSingleAllowHosts(t *testing.T) {
	router := newServer(Config{
		AllowedHosts: []string{"www.example.com"},
	})

	w := performRequest(router, "http://www.example.com/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "bar", string(w.Body()))
}

func TestGoodMulitipleAllowHosts(t *testing.T) {
	router := newServer(Config{
		AllowedHosts: []string{"www.example.com", "sub.example.com"},
	})

	w := performRequest(router, "http://sub.example.com/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "bar", string(w.Body()))
}

func TestBadSingleAllowHosts(t *testing.T) {
	router := newServer(Config{
		AllowedHosts: []string{"sub.example.com"},
	})

	w := performRequest(router, "http://www.example.com/foo")

	assert.DeepEqual(t, http.StatusForbidden, w.StatusCode())
}

func TestGoodMultipleAllowHosts(t *testing.T) {
	router := newServer(Config{
		AllowedHosts: []string{"www.example.com", "sub.example.com"},
	})

	w := performRequest(router, "http://sub.example.com/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "bar", string(w.Body()))
}

func TestBadMultipleAllowHosts(t *testing.T) {
	router := newServer(Config{
		AllowedHosts: []string{"www.example.com", "sub.example.com"},
	})

	w := performRequest(router, "http://www3.example.com/foo")

	assert.DeepEqual(t, http.StatusForbidden, w.StatusCode())
}

func TestAllowHostsInDevMode(t *testing.T) {
	router := newServer(Config{
		AllowedHosts:  []string{"www.example.com", "sub.example.com"},
		IsDevelopment: true,
	})

	w := performRequest(router, "http://www3.example.com/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
}

func TestBadHostHandler(t *testing.T) {
	badHandler := func(_ context.Context, c *app.RequestContext) {
		c.String(http.StatusInternalServerError, "BadHost")
		c.Abort()
	}

	router := newServer(Config{
		AllowedHosts:   []string{"www.example.com", "sub.example.com"},
		BadHostHandler: badHandler,
	})

	w := performRequest(router, "http://www3.example.com/foo")

	assert.DeepEqual(t, http.StatusInternalServerError, w.StatusCode())
	assert.DeepEqual(t, "BadHost", string(w.Body()))
}

func TestSSL(t *testing.T) {
	router := newServer(Config{
		SSLRedirect: true,
	})

	w := performRequest(router, "https://www.example.com/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "bar", string(w.Body()))
}

func TestSSLInDevMode(t *testing.T) {
	router := newServer(Config{
		SSLRedirect:   true,
		IsDevelopment: true,
	})

	w := performRequest(router, "http://www.example.com/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "bar", string(w.Body()))
}

func TestBasicSSL(t *testing.T) {
	router := newServer(Config{
		SSLRedirect: true,
	})

	w := performRequest(router, "http://www.example.com/foo")

	assert.DeepEqual(t, http.StatusMovedPermanently, w.StatusCode())
	assert.DeepEqual(t, "https://www.example.com/foo", w.Header.Get("Location"))
}

func TestDontRedirectIPV4Hostnames(t *testing.T) {
	router := newServer(Config{
		SSLRedirect:               true,
		DontRedirectIPV4Hostnames: true,
	})

	w1 := performRequest(router, "http://www.example.com/foo")
	assert.DeepEqual(t, http.StatusMovedPermanently, w1.StatusCode())

	w2 := performRequest(router, "http://127.0.0.1/foo")
	assert.DeepEqual(t, http.StatusOK, w2.StatusCode())
}

func TestBasicSSLWithHost(t *testing.T) {
	router := newServer(Config{
		SSLRedirect: true,
		SSLHost:     "secure.example.com",
	})

	w := performRequest(router, "http://www.example.com/foo")

	assert.DeepEqual(t, http.StatusMovedPermanently, w.StatusCode())
	assert.DeepEqual(t, "https://secure.example.com/foo", w.Header.Get("Location"))
}

func TestBadProxySSL(t *testing.T) {
	var req protocol.Request
	req.Header.Add("X-Forwarded-Proto", "https")
	engine := newServer(Config{SSLRedirect: true})
	w := ut.PerformRequest(engine, "GET", "http://www.example.com/foo", nil, ut.Header{
		Key:   "X-Forwarded-Proto",
		Value: "https",
	})
	resp := w.Result()
	assert.DeepEqual(t, http.StatusMovedPermanently, resp.StatusCode())
	assert.DeepEqual(t, "https://www.example.com/foo", w.Header().Get("Location"))
}

func TestProxySSLWithHeaderOption(t *testing.T) {
	h := server.New(server.WithHostPorts("127.0.0.1:8001"))
	h.Use(New(Config{
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Arbitrary-Header": "arbitrary-value"},
	}))
	h.GET("/foo", func(_ context.Context, c *app.RequestContext) {
		c.String(200, testResponse)
	})
	go h.Spin()
	time.Sleep(200 * time.Millisecond)
	client := http.Client{}
	req1, _ := http.NewRequest(consts.MethodGet, "http://127.0.0.1:8001/foo", nil)
	req1.Host = "www.example.com"
	req1.URL.Scheme = "http"
	req1.Header.Add("X-Arbitrary-Header", "arbitrary-value")
	resp, _ := client.Do(req1)
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode)
}

func TestProxySSLWithWrongHeaderValue(t *testing.T) {
	engine := newServer(Config{
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Arbitrary-Header": "arbitrary-value"},
	})

	resp := performRequest(engine, "http://www.example.com/foo", ut.Header{
		Key:   "X-Arbitrary-Header",
		Value: "wrong-value",
	})
	assert.DeepEqual(t, http.StatusMovedPermanently, resp.StatusCode())
	assert.DeepEqual(t, "https://www.example.com/foo", resp.Header.Get("Location"))
}

func TestStsHeader(t *testing.T) {
	router := newServer(Config{
		STSSeconds: 315360000,
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "max-age=315360000", w.Header.Get("Strict-Transport-Security"))
}

func TestStsHeaderInDevMode(t *testing.T) {
	router := newServer(Config{
		STSSeconds:    315360000,
		IsDevelopment: true,
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "", w.Header.Get("Strict-Transport-Security"))
}

func TestStsHeaderWithSubdomain(t *testing.T) {
	router := newServer(Config{
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "max-age=315360000; includeSubdomains", w.Header.Get("Strict-Transport-Security"))
}

func TestFrameDeny(t *testing.T) {
	router := newServer(Config{
		FrameDeny: true,
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "DENY", w.Header.Get("X-Frame-Options"))
}

func TestCustomFrameValue(t *testing.T) {
	router := newServer(Config{
		CustomFrameOptionsValue: "SAMEORIGIN",
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "SAMEORIGIN", w.Header.Get("X-Frame-Options"))
}

func TestCustomFrameValueWithDeny(t *testing.T) {
	router := newServer(Config{
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "SAMEORIGIN", w.Header.Get("X-Frame-Options"))
}

func TestContentNosniff(t *testing.T) {
	router := newServer(Config{
		ContentTypeNosniff: true,
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "nosniff", w.Header.Get("X-Content-Type-Options"))
}

func TestXSSProtection(t *testing.T) {
	router := newServer(Config{
		BrowserXssFilter: true,
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "1; mode=block", w.Header.Get("X-XSS-Protection"))
}

func TestReferrerPolicy(t *testing.T) {
	router := newServer(Config{
		ReferrerPolicy: "strict-origin-when-cross-origin",
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "strict-origin-when-cross-origin", w.Header.Get("Referrer-Policy"))
}

func TestFeaturePolicy(t *testing.T) {
	router := newServer(Config{
		FeaturePolicy: "vibrate 'none';",
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "vibrate 'none';", w.Header.Get("Feature-Policy"))
}

func TestCsp(t *testing.T) {
	router := newServer(Config{
		ContentSecurityPolicy: "default-src 'self'",
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "default-src 'self'", w.Header.Get("Content-Security-Policy"))
}

func TestInlineSecure(t *testing.T) {
	router := newServer(Config{
		FrameDeny: true,
	})

	w := performRequest(router, "/foo")

	assert.DeepEqual(t, http.StatusOK, w.StatusCode())
	assert.DeepEqual(t, "DENY", w.Header.Get("X-Frame-Options"))
}

func TestIsIpv4Host(t *testing.T) {
	assert.DeepEqual(t, isIPV4("127.0.0.1"), true)
	assert.DeepEqual(t, isIPV4("127.0.0.1:8080"), true)
	assert.DeepEqual(t, isIPV4("localhost"), false)
	assert.DeepEqual(t, isIPV4("localhost:8080"), false)
	assert.DeepEqual(t, isIPV4("example.com"), false)
	assert.DeepEqual(t, isIPV4("example.com:8080"), false)
}

func performRequest(engine *route.Engine, url string, header ...ut.Header) *protocol.Response {
	return ut.PerformRequest(engine, consts.MethodGet, url, nil).Result()
}
