# Secure (This is a community driven project)

## Example
See the [example1](example/custom/main.go), [example2](example/default/main.go).

DefaultConfig returns a Configuration with strict security settings

[embedmd]:# (secure.go go /func DefaultConfig/ /^}$/)
```go
func DefaultConfig() Config {
	return Config{
		SSLRedirect:           true,
		IsDevelopment:         false,
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	}
}
```
[embedmd]:# (example/customize/main.go go)
```go
func main() {
	h := server.Default(
		server.WithHostPorts("127.0.0.1:8080"),
	)
	h.Use(secure.New(secure.Config{
		AllowedHosts:          []string{"example.com", "ssl.example.com"},
		SSLRedirect:           true,
		SSLHost:               "ssl.example.com",
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	}))

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.String(200, "pong")
	})
	h.Spin()
}
```