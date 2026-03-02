package webauthn

import (
	"context"
	"net/url"
	"slices"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/walnuts1018/PRFExample/server/definitions"
)

func NewWebAuthn(ctx context.Context, origin *url.URL, additionalOrigins ...string) (*webauthn.WebAuthn, error) {
	return webauthn.New(&webauthn.Config{
		RPDisplayName: definitions.ApplicationDisplayName,
		RPID:          origin.Hostname(),
		RPOrigins:     slices.Concat([]string{origin.String()}, additionalOrigins),
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationRequired,
		},
	})
}
