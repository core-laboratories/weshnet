package bertyvcissuer

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"strings"

	weshnet_vc "github.com/core-laboratories/weshnet/v2/pkg/vcissuer"
)

func makeAuthenticateURL(serverBaseRoot, flowCtxStr string) string {
	return makeIssuerURL(serverBaseRoot, weshnet_vc.PathAuthenticate, flowCtxStr)
}

func makeProofURL(serverBaseRoot, flowCtxStr string) string {
	return makeIssuerURL(serverBaseRoot, weshnet_vc.PathProof, flowCtxStr)
}

func makeRedirectSuccessURI(redirectURI, state string, credentials []byte) string {
	u, _ := url.Parse(redirectURI) // Redirect URIs are validated when the issuer is configured and on challenge.
	u.Path = strings.TrimRight(u.Path, "/") + weshnet_vc.PathProof
	query := u.Query()
	query.Set(weshnet_vc.ParamState, state)
	query.Set(weshnet_vc.ParamCredentials, base64.URLEncoding.EncodeToString(credentials))
	u.RawQuery = query.Encode()
	return u.String()
}

func makeIssuerURL(serverBaseRoot, path, flowCtxStr string) string {
	u, _ := url.Parse(serverBaseRoot) // The server root is validated when the issuer is configured.
	u.Path = strings.TrimRight(u.Path, "/") + path
	query := u.Query()
	query.Set(weshnet_vc.ParamContext, flowCtxStr)
	u.RawQuery = query.Encode()
	return u.String()
}

func validateIssuerRootURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u == nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid issuer server root URL")
	}
	return nil
}

func validateRedirectURI(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u == nil || u.Host == "" || u.User != nil || u.Opaque != "" || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("%w: invalid callback URL", ErrFlowRedirectURINotAllowed)
	}

	switch u.Scheme {
	case "tone", "https":
		return nil
	case "http":
		if host := u.Hostname(); host == "localhost" || net.ParseIP(host) != nil && net.ParseIP(host).IsLoopback() {
			return nil
		}
	}

	return fmt.Errorf("%w: unsupported callback scheme", ErrFlowRedirectURINotAllowed)
}
