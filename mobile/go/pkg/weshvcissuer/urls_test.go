package weshvcissuer

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/core-laboratories/weshnet/v2/pkg/verifiablecredstypes"
	"github.com/stretchr/testify/require"
)

func TestIssuerRejectsUntrustedCallback(t *testing.T) {
	issuer, err := New(&Config{
		ServerRootURL: "https://issuer.example",
		IssuerSignKey: &[32]byte{1},
		Flow: &FlowConfig{
			Type:             verifiablecredstypes.FlowType_FlowTypeCode,
			CodeGenerator:    CodeGeneratorZero,
			CodeSenderClient: &PhoneCodeSenderMockService{},
		},
	})
	require.NoError(t, err)

	for _, callback := range []string{
		"https://attacker.example",
		"javascript:alert(1)",
		"//attacker.example",
	} {
		req := httptest.NewRequest("GET", "/challenge?state=state&redirect_uri="+url.QueryEscape(callback), nil)
		res := httptest.NewRecorder()
		issuer.challenge(res, req)
		require.Equal(t, 400, res.Code, callback)
		require.Contains(t, res.Body.String(), ErrFlowRedirectURINotAllowed.Error())
	}
}

func TestRedirectURIValidation(t *testing.T) {
	for _, callback := range []string{"tone://vc", "https://example.com/callback", "http://localhost:8080/callback", "http://127.0.0.1/callback"} {
		require.NoError(t, validateRedirectURI(callback), callback)
	}
	for _, callback := range []string{"", "https://example.com@attacker.example", "https://example.com/callback?next=evil", "http://attacker.example", "javascript:alert(1)", "file:///etc/passwd"} {
		require.Error(t, validateRedirectURI(callback), callback)
	}
}

func TestRedirectURLsEncodeQueryValues(t *testing.T) {
	context := "a&next=https://attacker.example"
	authenticate, err := url.Parse(makeAuthenticateURL("https://issuer.example/base", context))
	require.NoError(t, err)
	require.Equal(t, "/base/authenticate", authenticate.Path)
	require.Equal(t, context, authenticate.Query().Get("context"))
	require.Equal(t, "", authenticate.Query().Get("next"))

	callback, err := url.Parse(makeRedirectSuccessURI("tone://vc", "a&next=evil", []byte{1, 2, 3}))
	require.NoError(t, err)
	require.Equal(t, "tone", callback.Scheme)
	require.Equal(t, "vc", callback.Host)
	require.Equal(t, "/proof", callback.Path)
	require.Equal(t, "a&next=evil", callback.Query().Get("state"))
	require.False(t, strings.Contains(callback.RawQuery, "&next="))
}
