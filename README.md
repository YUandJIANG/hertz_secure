# Secure (This is a community driven project)

`Secure` middleware for hertz framework.

This repo is forked from [secure](https://github.com/gin-contrib/secure) and adapted for hertz.

## Install

```bash
go get github.com/hertz-contrib/secure
```

### [Default example](example/default/main.go)

Default configuration for users to set security configuration directly using secure middleware

#### Sample Code

```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/secure"
)

func main() {
	h := server.New(server.WithHostPorts(":8080"))

	h.Use(secure.Default(
		secure.WithSSLHost("ssl.example.com"),
		secure.WithAllowedHosts([]string{"example.com", "ssl.example.com"}),
	))
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.String(200, "pong")
	})
	h.Spin()
}

```

### [Custom example](example/custom/main.go)

User passed in custom configuration items

#### Function Signature

```go
func New(opts ...Option) app.HandlerFunc
```

#### Sample Code

```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/secure"
)

func main() {
	h := server.Default(
		server.WithHostPorts("127.0.0.1:8080"),
	)
	h.Use(secure.New(
		secure.WithAllowedHosts([]string{"example.com", "ssl.example.com"}),
		secure.WithSSLRedirect(true),
		secure.WithSSLHost("ssl.example.com"),
		secure.WithSTSSecond(315360000),
		secure.WithSTSIncludeSubdomains(true),
		secure.WithFrameDeny(true),
		secure.WithContentTypeNosniff(true),
		secure.WithBrowserXssFilter(true),
		secure.WithContentSecurityPolicy("default-src 'self'"),
		secure.WithIENoOpen(true),
		secure.WithReferrerPolicy("strict-origin-when-cross-origin"),
		secure.WithSSLProxyHeaders(map[string]string{"X-Forwarded-Proto": "https"}),
	))

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, "pong")
	})
	h.Spin()
}
```

## Default Configuration

```go
    func Default(opts ...Option) app.HandlerFunc {
        options := []Option{
            WithSSLRedirect(true),
            WithIsDevelopment(false),
            WithSTSSecond(315360000),
            WithFrameDeny(true),
            WithContentTypeNosniff(true),
            WithBrowserXssFilter(true),
            WithContentSecurityPolicy("default-src 'self'"),
            WithIENoOpen(true),
            WithSSLProxyHeaders(map[string]string{"X-Forwarded-Proto": "https"}),
        }
        options = append(options, opts...)
        return New(options...)
    }
```

## Option

| options                       | Parameters        | Description                                                  |
| ----------------------------- | ----------------- | ------------------------------------------------------------ |
| WithSSLRedirect               | bool              | If `WithSSLRedirect` is set to true, then only allow https requests |
| WithIsDevelopment             | bool              | When true, the whole security policy applied by the middleware is disabled completely. |
| WithSTSSecond                 | int64             | Default is 315360000, which would NOT include the header.    |
| WithSTSIncludeSubdomains      | bool              | If `WithSTSIncludeSubdomains` is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false. |
| WithFrameDeny                 | bool              | If `WithFrameDeny` is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false |
| WithContentTypeNosniff        | bool              | If `WithContentTypeNosniff` is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false. |
| WithBrowserXssFilter          | bool              | If `WithBrowserXssFilter` is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false. |
| WithContentSecurityPolicy     | []string          | `WithContentSecurityPolicy` allows the Content-Security-Policy header value to be set with a custom value. Default is "". |
| WithIENoOpen                  | bool              | Prevent Internet Explorer from executing downloads in your site’s context |
| WithSSLProxyHeaders           | map[string]string | This is useful when your app is running behind a secure proxy that forwards requests to your app over http (such as on Heroku). |
| WithAllowedHosts              | []string          | `WithAllowedHosts` is a list of fully qualified domain names that are allowed.Default is empty list, which allows any and all host names. |
| WithSSLTemporaryRedirect      | bool              | If `WithSSLTemporaryRedirect` is true, the a 302 will be used while redirecting. Default is false (301). |
| WithSSLHost                   | string            | `WithSSLHost` is the host name that is used to redirect http requests to https. Default is "", which indicates to use the same host. |
| WithCustomFrameOptionsValue   | string            | `WithCustomFrameOptionsValue` allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. |
| WithReferrerPolicy            | string            | HTTP header "Referrer-Policy" governs which referrer information, sent in the Referrer header, should be included with requests made. |
| WithBadHostHandler            | app.HandlerFunc   | Handlers for when an error occurs (ie bad host).             |
| WithFeaturePolicy             | string            | Feature Policy is a new header that allows a site to control which features and APIs can be used in the browser. |
| WithDontRedirectIPV4Hostnames | bool              | If `WithDontRedirectIPV4Hostnames` is true, requests to hostnames that are IPV4 addresses aren't redirected. This is to allow load balancer health checks  to succeed. |

## License

This project is under Apache License. See the [LICENSE](LICENSE) file for the full license text.