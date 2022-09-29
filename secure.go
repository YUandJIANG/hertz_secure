package secure

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type (
	Option func(o *options)

	// options is a struct for specifying configuration options for the secure.
	options struct {
		// AllowedHosts is a list of fully qualified domain names that are allowed.
		// Default is empty list, which allows any and all host names.
		allowedHosts []string
		// If WithSSLRedirect is set to true, then only allow https requests.
		// Default is false.
		WithSSLRedirect bool
		// If SSLTemporaryRedirect is true, the a 302 will be used while redirecting.
		// Default is false (301).
		sslTemporaryRedirect bool
		// SSLHost is the host name that is used to redirect http requests to https.
		// Default is "", which indicates to use the same host.
		sslHost string
		// STSSeconds is the max-age of the Strict-Transport-Security header.
		// Default is 0, which would NOT include the header.
		stsSeconds int64
		// If STSIncludeSubdomains is set to true, the `includeSubdomains` will
		// be appended to the Strict-Transport-Security header. Default is false.
		stsIncludeSubdomains bool
		// If FrameDeny is set to true, adds the X-Frame-Options header with
		// the value of `DENY`. Default is false.
		frameDeny bool
		// CustomFrameOptionsValue allows the X-Frame-Options header value
		// to be set with a custom value. This overrides the FrameDeny option.
		customFrameOptionsValue string
		// If ContentTypeNosniff is true, adds the X-Content-Type-Options header
		// with the value `nosniff`. Default is false.
		contentTypeNosniff bool
		// If BrowserXssFilter is true, adds the X-XSS-Protection header with
		// the value `1; mode=block`. Default is false.
		browserXssFilter bool
		// ContentSecurityPolicy allows the Content-Security-Policy header value
		// to be set with a custom value. Default is "".
		contentSecurityPolicy string
		// HTTP header "Referrer-Policy" governs which referrer information, sent in the Referrer header, should be included with requests made.
		referrerPolicy string
		// When true, the whole security policy applied by the middleware is disabled completely.
		isDevelopment bool
		// Handlers for when an error occurs (ie bad host).
		badHostHandler app.HandlerFunc
		// Prevent Internet Explorer from executing downloads in your site’s context
		ieNoOpen bool
		// Feature Policy is a new header that allows a site to control which features and APIs can be used in the browser.
		featurePolicy string
		// If DontRedirectIPV4Hostnames is true, requests to hostnames that are IPV4
		// addresses aren't redirected. This is to allow load balancer health checks
		// to succeed.
		dontRedirectIPV4Hostnames bool

		// If the request is insecure, treat it as secure if any of the headers in this dict are set to their corresponding value
		// This is useful when your app is running behind a secure proxy that forwards requests to your app over http (such as on Heroku).
		sslProxyHeaders map[string]string
	}
)

// WithAllowedHosts is a list of fully qualified domain names that are allowed.
// Default is empty list, which allows any and all host names.
func WithAllowedHosts(ss []string) Option {
	return func(o *options) {
		o.allowedHosts = ss
	}
}

// WithSSLRedirect when WithSSLRedirect is set to true, then only allow https requests.
// Default is false.
func WithSSLRedirect(b bool) Option {
	return func(o *options) {
		o.WithSSLRedirect = b
	}
}

// WithSSLTemporaryRedirect when SSLTemporaryRedirect is true, the a 302 will be used while redirecting.
// Default is false (301).
func WithSSLTemporaryRedirect(b bool) Option {
	return func(o *options) {
		o.sslTemporaryRedirect = b
	}
}

// WithSSLHost is the host name that is used to redirect http requests to https.
// Default is "", which indicates to use the same host.
func WithSSLHost(s string) Option {
	return func(o *options) {
		o.sslHost = s
	}
}

// WithSTSSecond is the max-age of the Strict-Transport-Security header.
// Default is 0, which would NOT include the header.
func WithSTSSecond(sec int64) Option {
	return func(o *options) {
		o.stsSeconds = sec
	}
}

// WithSTSIncludeSubdomains when STSIncludeSubdomains is set to true, the `includeSubdomains` will
// be appended to the Strict-Transport-Security header. Default is false.
func WithSTSIncludeSubdomains(b bool) Option {
	return func(o *options) {
		o.stsIncludeSubdomains = b
	}
}

// WithFrameDeny when FrameDeny is set to true, adds the X-Frame-Options header with
// the value of `DENY`. Default is false.
func WithFrameDeny(b bool) Option {
	return func(o *options) {
		o.frameDeny = b
	}
}

// WithCustomFrameOptionsValue allows the X-Frame-Options header value
// to be set with a custom value. This overrides the FrameDeny option.
func WithCustomFrameOptionsValue(s string) Option {
	return func(o *options) {
		o.customFrameOptionsValue = s
	}
}

// WithContentTypeNosniff when ContentTypeNosniff is true, adds the X-Content-Type-Options header
// with the value `nosniff`. Default is false.
func WithContentTypeNosniff(b bool) Option {
	return func(o *options) {
		o.contentTypeNosniff = b
	}
}

// WithBrowserXssFilter when BrowserXssFilter is true, adds the X-XSS-Protection header with
// the value `1; mode=block`. Default is false.
func WithBrowserXssFilter(b bool) Option {
	return func(o *options) {
		o.browserXssFilter = b
	}
}

// WithContentSecurityPolicy  allows the Content-Security-Policy header value
// to be set with a custom value. Default is "".
func WithContentSecurityPolicy(s string) Option {
	return func(o *options) {
		o.contentSecurityPolicy = s
	}
}

// WithReferrerPolicy use to set HTTP header "Referrer-Policy" governs which referrer information,
// sent in the Referrer header,/should be included with requests made.
func WithReferrerPolicy(s string) Option {
	return func(o *options) {
		o.referrerPolicy = s
	}
}

// WithIsDevelopment when true, the whole security policy applied by the middleware is disabled completely.
func WithIsDevelopment(b bool) Option {
	return func(o *options) {
		o.isDevelopment = b
	}
}

// WithIENoOpen prevents Internet Explorer from executing downloads in your site’s context
func WithIENoOpen(b bool) Option {
	return func(o *options) {
		o.ieNoOpen = b
	}
}

// WithBadHostHandler use to when an error occurs (ie bad host).
func WithBadHostHandler(handler app.HandlerFunc) Option {
	return func(o *options) {
		o.badHostHandler = handler
	}
}

// WithFeaturePolicy  is a new header that allows a site to control which features and APIs can be used in the browser.
func WithFeaturePolicy(s string) Option {
	return func(o *options) {
		o.featurePolicy = s
	}
}

// WithDontRedirectIPV4Hostnames when DontRedirectIPV4Hostnames is true, requests to hostnames that are IPV4
// addresses aren't redirected. This is to allow load balancer health checks
// to succeed.
func WithDontRedirectIPV4Hostnames(b bool) Option {
	return func(o *options) {
		o.dontRedirectIPV4Hostnames = b
	}
}

// WithSSLProxyHeaders If the request is insecure, treat it as secure if any of the headers in this dict are set to their corresponding value
// This is useful when your app is running behind a secure proxy that forwards requests to your app over http (such as on Heroku).
func WithSSLProxyHeaders(m map[string]string) Option {
	return func(o *options) {
		o.sslProxyHeaders = m
	}
}

// Default returns a Configuration with strict security settings.
// ```
//		WithSSLRedirect:           true
//		IsDevelopment:         false
//		STSSeconds:            315360000
//		STSIncludeSubdomains:  true
//		FrameDeny:             true
//		ContentTypeNosniff:    true
//		BrowserXssFilter:      true
//		ContentSecurityPolicy: "default-src 'self'"
//		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
// ```
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

func (o *options) Apply(opts []Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// New creates an instance of the secure middleware using the specified configuration.
// router.Use(secure.N)
func New(opts ...Option) app.HandlerFunc {
	policy := newPolicy(opts)
	return func(ctx context.Context, c *app.RequestContext) {
		if !policy.applyToContext(ctx, c) {
			return
		}
		c.Next(ctx)
	}
}
