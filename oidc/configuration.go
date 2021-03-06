// Package oidc provides functions for interacting with OpenID Connect providers
package oidc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Configuration for connecting to OIDC provider
// Documentation and fields copied from
// https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
type Configuration struct {
	// REQUIRED. URL using the https scheme with no query or fragment component
	// that the OP asserts as its Issuer Identifier. If Issuer discovery is
	// supported (see Section 2), this value MUST be identical to the issuer value
	// returned by WebFinger. This also MUST be identical to the iss Claim value
	// in ID Tokens issued from this Issuer.
	Issuer string `json:"issuer"`
	// REQUIRED. URL of the OP's OAuth 2.0 Authorization Endpoint
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	// URL of the OP's OAuth 2.0 Token Endpoint. This is REQUIRED unless only
	// the Implicit Flow is used.
	TokenEndpoint string `json:"token_endpoint"`
	// RECOMMENDED. URL of the OP's UserInfo Endpoint. This URL MUST use the
	// https scheme and MAY contain port, path, and query parameter components.
	UserInfoEndpoint string `json:"userinfo_endpoint"`
	// REQUIRED. URL of the OP's JSON Web Key Set [JWK] document. This contains
	// the signing key(s) the RP uses to validate signatures from the OP.
	// The JWK Set MAY also contain the Server's encryption key(s), which are used
	// by RPs to encrypt requests to the Server. When both signing and encryption
	// keys are made available, a use (Key Use) parameter value is REQUIRED for
	// all keys in the referenced JWK Set to indicate each key's intended usage.
	// Although some algorithms allow the same key to be used for both signatures
	// and encryption, doing so is NOT RECOMMENDED, as it is less secure.
	// The JWK x5c parameter MAY be used to provide X.509 representations of keys
	// provided. When used, the bare key values MUST still be present and MUST
	// match those in the certificate.
	JWKSURI string `json:"jwks_uri"`
	// RECOMMENDED. URL of the OP's Dynamic Client Registration Endpoint
	RegistrationEndpoint string `json:"registration_endpoint"`
	// RECOMMENDED. JSON array containing a list of the OAuth 2.0 scope values
	// that this server supports. The server MUST support the openid scope value.
	// Servers MAY choose not to advertise some supported scope values even when
	// this parameter is used, although those defined in SHOULD be listed,
	// if supported.
	ScopesSupported []string `json:"scopes_supported"`
	// REQUIRED. JSON array containing a list of the OAuth 2.0 response_type
	// values that this OP supports. Dynamic OpenID Providers MUST support
	// the code, id_token, and the token id_token Response Type values.
	ResponseTypesSupported []string `json:"response_types_supported"`
	// OPTIONAL. JSON array containing a list of the OAuth 2.0 response_mode
	// values that this OP supports, as specified in OAuth 2.0 Multiple Response
	// Type Encoding Practices. If omitted, the default for Dynamic OpenID
	// Providers is ["query", "fragment"].
	ResponseModesSupported []string `json:"response_modes_supported"`
	// OPTIONAL. JSON array containing a list of the OAuth 2.0 Grant Type values
	// that this OP supports. Dynamic OpenID Providers MUST support
	// the authorization_code and implicit Grant Type values and MAY support
	// other Grant Types. If omitted, the default value is
	// ["authorization_code", "implicit"].
	GrantTypesSupported []string `json:"grant_types_supported"`
	// OPTIONAL. JSON array containing a list of the Authentication Context Class
	// References that this OP supports.
	ACRValuesSupported []string `json:"acr_values_supported"`
	// REQUIRED. JSON array containing a list of the Subject Identifier types
	// hat this OP supports. Valid types include pairwise and public.
	SubjectTypesSupported []string `json:"subject_types_supported"`
	// REQUIRED. JSON array containing a list of the JWS signing algorithms
	// (alg values) supported by the OP for the ID Token to encode the Claims
	// in a JWT. The algorithm RS256 MUST be included. The value none MAY be
	// supported, but MUST NOT be used unless the Response Type used returns
	// no ID Token from the Authorization Endpoint (such as when using the
	// Authorization Code Flow).
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	// OPTIONAL. JSON array containing a list of the JWE encryption algorithms
	// (enc values) supported by the OP for the ID Token to encode the Claims
	// in a JWT.
	IDTokenEncryptionEncValuesSupported []string `json:"id_token_encryption_enc_values_supported"`
	// OPTIONAL. JSON array containing a list of the JWS signing algorithms
	// (alg values) supported by the UserInfo Endpoint to encode the Claims
	// in a JWT. The value none MAY be included.
	UserInfoSigningAlgValuesSupported []string `json:"userinfo_signing_alg_values_supported"`
	// OPTIONAL. JSON array containing a list of the JWE encryption algorithms
	// (alg values) supported by the UserInfo Endpoint to encode the Claims
	// in a JWT.
	UserInfoEncryptionAlgValuesSupported []string `json:"userinfo_encryption_alg_values_supported"`
	// OPTIONAL. JSON array containing a list of the JWE encryption algorithms
	// (enc values) supported by the UserInfo Endpoint to encode the Claims
	// in a JWT.
	UserInfoEncryptionEncValuesSupported []string `json:"userinfo_encryption_enc_values_supported"`
	// OPTIONAL. JSON array containing a list of the JWS signing algorithms
	// (alg values) supported by the OP for Request Objects, which are described
	// in Section 6.1 of OpenID Connect Core 1.0. These algorithms are used both
	// when the Request Object is passed by value (using the request parameter)
	// and when it is passed by reference (using the request_uri parameter).
	// Servers SHOULD support none and RS256.
	RequestObjectSigningAlgValuesSupported []string `json:"request_object_signing_alg_values_supported"`
	// OPTIONAL. JSON array containing a list of the JWE encryption algorithms
	// (alg values) supported by the OP for Request Objects. These algorithms are
	// used both when the Request Object is passed by value and when it is
	// passed by reference.
	RequestObjectEncryptionAlgValuesSupported []string `json:"request_object_encryption_alg_values_supported"`
	// OPTIONAL. JSON array containing a list of the JWE encryption algorithms
	// (enc values) supported by the OP for Request Objects. These algorithms are
	// used both when the Request Object is passed by value and when it is
	// passed by reference.
	RequestObjectEncryptionEncValuesSupported []string `json:"request_object_encryption_enc_values_supported"`
	// OPTIONAL. JSON array containing a list of Client Authentication methods
	// supported by this Token Endpoint. The options are client_secret_post,
	// client_secret_basic, client_secret_jwt, and private_key_jwt, as described
	// in Section 9 of OpenID Connect Core 1.0. Other authentication methods MAY
	// be defined by extensions. If omitted, the default is client_secret_basic
	// -- the HTTP Basic Authentication Scheme specified in Section 2.3.1 of
	// OAuth 2.0.
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	// OPTIONAL. JSON array containing a list of the JWS signing algorithms
	// (alg values) supported by the Token Endpoint for the signature on the JWT
	// used to authenticate the Client at the Token Endpoint for the
	// private_key_jwt and client_secret_jwt authentication methods. Servers
	// SHOULD support RS256. The value none MUST NOT be used.
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported"`
	// OPTIONAL. JSON array containing a list of the display parameter values that
	// the OpenID Provider supports. These values are described in Section 3.1.2.1
	// of OpenID Connect Core 1.0.
	DisplayValuesSupported []string `json:"display_values_supported"`
	// OPTIONAL. JSON array containing a list of the Claim Types that the OpenID
	// Provider supports. These Claim Types are described in Section 5.6 of OpenID
	// Connect Core 1.0. Values defined by this specification are normal,
	// aggregated, and distributed. If omitted, the implementation supports only
	// normal Claims.
	ClaimTypesSupported []string `json:"claim_types_supported"`
	// RECOMMENDED. JSON array containing a list of the Claim Names of the Claims
	// that the OpenID Provider MAY be able to supply values for. Note that for
	// privacy or other reasons, this might not be an exhaustive list.
	ClaimsSupported []string `json:"claims_supported"`
	// OPTIONAL. URL of a page containing human-readable information that
	// developers might want or need to know when using the OpenID Provider.
	// In particular, if the OpenID Provider does not support Dynamic Client
	// Registration, then information on how to register Clients needs to be
	// provided in this documentation.
	ServiceDocumentation string `json:"service_documentation"`
	// OPTIONAL. Languages and scripts supported for values in Claims being
	// returned, represented as a JSON array of BCP47 language tag values. Not all
	// languages and scripts are necessarily supported for all Claim values.
	ClaimsLocalesSupported []string `json:"claims_locales_supported"`
	// OPTIONAL. Languages and scripts supported for the user interface,
	// represented as a JSON array of BCP47 language tag values.
	UILocalesSupported []string `json:"ui_locales_supported"`
	// OPTIONAL. Boolean value specifying whether the OP supports use of the
	// claims parameter, with true indicating support. If omitted, the default
	// value is false.
	ClaimsParameterSupported bool `json:"claims_parameter_supported"`
	// OPTIONAL. Boolean value specifying whether the OP supports use of the
	// request parameter, with true indicating support. If omitted, the
	// default value is false.
	RequestParameterSupported bool `json:"request_parameter_supported"`
	// OPTIONAL. Boolean value specifying whether the OP supports use of the
	// request_uri parameter, with true indicating support. If omitted, the
	// default value is true.
	RequestURIParameterSupported *bool `json:"request_uri_parameter_supported"`
	// OPTIONAL. Boolean value specifying whether the OP requires any request_uri
	// values used to be pre-registered using the request_uris registration
	// parameter. Pre-registration is REQUIRED when the value is true. If omitted,
	// the default value is false.
	RequireRequestURIRegistration bool `json:"require_request_uri_registration"`
	// OPTIONAL. URL that the OpenID Provider provides to the person registering
	// the Client to read about the OP's requirements on how the Relying Party can
	// use the data provided by the OP. The registration process SHOULD display
	// this URL to the person registering the Client if it is given.
	OpPolicyURI string `json:"op_policy_uri"`
	// OPTIONAL. URL that the OpenID Provider provides to the person registering
	// the Client to read about OpenID Provider's terms of service. The
	// registration process SHOULD display this URL to the person registering
	// the Client if it is given.
	OpTOSURI string `json:"op_tos_uri"`
}

// Valid checks that all the fields marked as required are present.
func (config Configuration) Valid() error {
	if config.Issuer == "" {
		return errors.New("authproxy/oidc: issuer is empty")
	}
	if config.AuthorizationEndpoint == "" {
		return errors.New("authproxy/oidc: authorization_endpoint is empty")
	}
	if config.JWKSURI == "" {
		return errors.New("authproxy/oidc: jwks_uri is empty")
	}
	if len(config.ResponseTypesSupported) == 0 {
		return errors.New("authproxy/oidc: response_types_supported is empty")
	}
	if len(config.SubjectTypesSupported) == 0 {
		return errors.New("authproxy/oidc: subject_types_supported is empty")
	}
	if len(config.IDTokenSigningAlgValuesSupported) == 0 {
		return errors.New("authproxy/oidc: id_token_signing_alg_values_supported is empty")
	}

	return nil
}

// FillDefaultValuesIfEmpty with the following replacements:
//  - ResponseModesSupported: ["query", "fragment"]
//  - GrantTypesSupported: ["authorization_code", "implicit"]
//  - TokenEndpointAuthMethodsSupported: ["client_secret_basic"]
//  - RequestURIParameterSupported: true
func (config *Configuration) FillDefaultValuesIfEmpty() {
	if len(config.ResponseModesSupported) == 0 {
		config.ResponseModesSupported = []string{"query", "fragment"}
	}
	if len(config.GrantTypesSupported) == 0 {
		config.GrantTypesSupported = []string{"authorization_code", "implicit"}
	}
	if len(config.TokenEndpointAuthMethodsSupported) == 0 {
		config.TokenEndpointAuthMethodsSupported = []string{"client_secret_basic"}
	}
	if config.RequestURIParameterSupported == nil {
		b := true
		config.RequestURIParameterSupported = &b
	}
}

// GetOpenIDConnectConfiguration from a well-known url. Should be Issuer + /.well-known/openid-configuration
func GetOpenIDConnectConfiguration(ctx context.Context, client *http.Client, url string) (config Configuration, err error) {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return config, errors.Wrap(err, "authproxy/oidc: couldn't create HTTP request")
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return config, errors.Wrap(err, "authproxy/oidc: HTTP request failed")
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			err = errors.Wrap(errClose, "authproxy/oidc: couldn't close HTTP response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return config, errors.Errorf("authproxy:oidc received non-200 status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&config)
	if err != nil {
		return config, errors.Wrap(err, "authproxy/oidc: couldn't decode HTTP response body to JSON")
	}

	config.FillDefaultValuesIfEmpty()

	return config, nil
}
