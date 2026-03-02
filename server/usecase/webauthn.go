package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
	"k8s.io/utils/ptr"
)

func (u *Usecase) BeginWebAuthnRegistration(ctx context.Context) (userID entity.UserID, creation *protocol.CredentialCreation, session *webauthn.SessionData, err error) {
	// Passkey = Userということにする
	user, err := u.createTemporaryUser(ctx)
	if err != nil {
		return entity.UserID{}, nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	existingCredentials, err := u.webAuthnCredentialRepository.ListWebAuthnCredentialsByUserID(ctx, user.ID)
	if err != nil {
		return entity.UserID{}, nil, nil, fmt.Errorf("failed to list WebAuthn credentials: %w", err)
	}

	var credentialExcludeList []protocol.CredentialDescriptor
	for cred := range existingCredentials {
		credentialExcludeList = append(credentialExcludeList, protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: []byte(cred.ID),
		})
	}

	options := []webauthn.RegistrationOption{
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementPreferred),
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			RequireResidentKey: ptr.To(false),
			ResidentKey:        protocol.ResidentKeyRequirementPreferred,
			// UserVerificationの有無でPRFの出力が変わるらしい
			// https://fidoalliance.org/specs/fido-v2.2-rd-20230321/fido-client-to-authenticator-protocol-v2.2-rd-20230321.html#sctn-hmac-secret-extension:~:text=If%20uv%20bit%20is%20set%20to%201%20in%20the%20response,%20let%20CredRandom%20be%20CredRandomWithUV.
			UserVerification: protocol.VerificationRequired,
		}),
		// Extensions
		webauthn.WithExtensions(
			map[string]any{
				"prf": map[string]any{
					"eval": map[string]any{
						"first": []byte(user.PRFSalt),
					},
				},
			},
		),
		webauthn.WithExclusions(credentialExcludeList),
	}

	webauthnUser := &entity.WebAuthnUser{
		User:        user,
		Credentials: []webauthn.Credential{},
	}

	creation, session, err = u.webAuthn.BeginRegistration(webauthnUser, options...)
	if err != nil {
		return entity.UserID{}, nil, nil, fmt.Errorf("failed to begin WebAuthn registration: %w", err)
	}

	return user.ID, creation, session, nil
}

func (u *Usecase) FinishWebAuthnRegistration(ctx context.Context, userID entity.UserID, creationData *protocol.ParsedCredentialCreationData, session *webauthn.SessionData) (entity.User, entity.WebAuthnCredential, error) {
	user, err := u.getUser(ctx, userID, true)
	if err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to get user: %w", err)
	}

	webauthnUser := &entity.WebAuthnUser{
		User:        user,
		Credentials: []webauthn.Credential{},
	}

	rawCred, err := u.webAuthn.CreateCredential(webauthnUser, *session, creationData)
	if err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to create WebAuthn credential: %w", err)
	}

	cred, err := entity.NewCredential(user, rawCred)
	if err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to create WebAuthn credential: %w", err)
	}

	if err := u.webAuthnCredentialRepository.StoreWebAuthnCredential(ctx, cred); err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to store WebAuthn credential: %w", err)
	}

	if err := u.userRepository.PromoteTemporaryUser(ctx, user.ID); err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to promote temporary user: %w", err)
	}

	return user, cred, nil
}

func (u *Usecase) BeginWebAuthnLogin(ctx context.Context, userID entity.UserID) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	user, err := u.getUser(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	options := []webauthn.LoginOption{
		webauthn.WithUserVerification(protocol.VerificationRequired),
		webauthn.WithAssertionExtensions(
			map[string]any{
				"prf": map[string]any{
					"eval": map[string]any{
						"first": []byte(user.PRFSalt),
					},
				},
			},
		),
	}

	credentials, err := u.webAuthnCredentialRepository.ListWebAuthnCredentialsByUserID(ctx, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list WebAuthn credentials: %w", err)
	}

	assertion, session, err := u.webAuthn.BeginLogin(&entity.WebAuthnUser{
		User: user,
		Credentials: slices.Collect(func(yield func(webauthn.Credential) bool) {
			for cred := range credentials {
				if !yield(*cred.Credential) {
					return
				}
			}
		}),
	}, options...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin WebAuthn login: %w", err)
	}
	return assertion, session, nil
}

func (u *Usecase) FinishWebAuthnLogin(ctx context.Context, userID entity.UserID, session webauthn.SessionData, assertion *protocol.ParsedCredentialAssertionData) (entity.User, entity.WebAuthnCredential, error) {
	user, err := u.getUser(ctx, userID)
	if err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to get user: %w", err)
	}
	creds, err := u.webAuthnCredentialRepository.ListWebAuthnCredentialsByUserID(ctx, user.ID)
	if err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to list WebAuthn credentials: %w", err)
	}

	cred, err := u.webAuthn.ValidateLogin(&entity.WebAuthnUser{
		User: user,
		Credentials: slices.Collect(func(yield func(webauthn.Credential) bool) {
			for cred := range creds {
				if !yield(*cred.Credential) {
					return
				}
			}
		}),
	}, session, assertion)
	if err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to finish WebAuthn login: %w", err)
	}

	if cred.Authenticator.CloneWarning {
		return entity.User{}, entity.WebAuthnCredential{}, errors.New("authenticator clone warning")
	}

	if err := u.webAuthnCredentialRepository.UpdateWebAuthnCredentialOnLogin(ctx, cred.ID, cred); err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to update WebAuthn credential on login: %w", err)
	}

	wc, err := u.webAuthnCredentialRepository.GetWebAuthnCredential(ctx, cred.ID)
	if err != nil {
		return entity.User{}, entity.WebAuthnCredential{}, fmt.Errorf("failed to get webauthn credential: %w", err)
	}
	wc.Credential = cred // 一応

	return user, wc, nil
}

func (u *Usecase) GetWebAuthnCredential(ctx context.Context, user entity.User, id entity.WebAuthnCredentialID) (entity.WebAuthnCredential, error) {
	credential, err := u.webAuthnCredentialRepository.GetWebAuthnCredential(ctx, id)
	if err != nil {
		return entity.WebAuthnCredential{}, fmt.Errorf("failed to get WebAuthn credential: %w", err)
	}

	if credential.UserID != user.ID {
		return entity.WebAuthnCredential{}, errors.New("invalid user")
	}

	return credential, nil
}

func (u *Usecase) DeleteWebAuthnCredential(ctx context.Context, user entity.User, id entity.WebAuthnCredentialID) error {
	if err := u.webAuthnCredentialRepository.DeleteWebAuthnCredential(ctx, user.ID, id); err != nil {
		return fmt.Errorf("failed to delete WebAuthn credential: %w", err)
	}
	return nil
}
