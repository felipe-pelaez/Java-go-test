package middleware

import (
	"MicroserviceTemplate/pkg/web"
	"context"
	"crypto/tls"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

// ? ==================== Structs ==================== ?

// Claims is the structure that maps the content of the JWT
type Claims struct {
	RealmAccess    clientRoles `json:"realm_access,omitempty"`
	ResourceAccess client      `json:"resource_access,omitempty"`
	JTI            string      `json:"jti,omitempty"`
}

// * =========== *

// Client is the structure containing the client's characteristics
type client struct {
	Gateway clientRoles `json:"Gateway,omitempty"`
}

// * =========== *

// ClientRoles is the structure containing the client's roles inside an array of strings
type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

// ? ==================== Functions ==================== ?

// AuthorizationFailed returns an authorization error in case the token is not valid for gin
func authorizationFailed(message string, c *gin.Context) {

	data := web.ErrorResponse{
		Code:    "authorization_failed",
		Status:  http.StatusUnauthorized,
		Message: message,
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"response": data})

}

// ? ==================== Middlewares ==================== ?

// IsAuthorizedJWT is the middleware that is in charge of validating the JWT token and verifying that the user has the necessary permissions to access the route
func IsAuthorizedJWT(excludePaths ...string) gin.HandlerFunc {

	Realm := viper.GetString("keycloak.realm")
	RealmConfigURL := viper.GetString("keycloak.url") + "/realms/" + Realm

	return func(c *gin.Context) {

		for _, path := range excludePaths {

			mainRoute := ""

			if strings.Contains(path, "/**") {
				mainRoute = strings.Split(path, "/**")[0]
			} else if strings.Contains(path, "/*any") {
				mainRoute = strings.Split(path, "/*any")[0]
			} else {
				mainRoute = path
			}

			if strings.Contains(c.Request.URL.Path, mainRoute) {
				c.Next()
				return
			}

		}

		// The header token is obtained by means of the Authorization key of type Bearer.
		rawAccessToken := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)

		// The transport is responsible for validating the token using the certificate authority's certificate (TLS for HTTPS).
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		// An HTTP client is created with the previously created transport
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}

		// An OpenId Connector context is created with the previously created HTTP client
		ctx := oidc.ClientContext(context.Background(), client)

		RealmUrl := strings.ReplaceAll(RealmConfigURL, "\\", "")

		// An OpenId Connector provider is created with the authorization provider (Keycloak) and the previously created context
		provider, err := oidc.NewProvider(ctx, RealmUrl)
		if err != nil {
			authorizationFailed("authorization failed while getting the provider: "+err.Error(), c) // An authorization error is returned in case the supplier is invalid.
			return
		}

		// An OpenId Connector configuration object is created based on the ClientID
		oidcConfig := &oidc.Config{
			ClientID: "account",
		}

		// The integrity of the token is validated using the authorization provider (Keycloak) and the OpenId Connector configuration created earlier.

		verifier := provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(ctx, rawAccessToken)

		if err != nil {
			authorizationFailed("authorization failed while verifying the token: "+err.Error(), c) // An authorization error is returned in case the token is invalid.
			return
		}

		// If the token is valid, a Claims object is created to map the contents of the token.
		var IDTokenClaims Claims
		if err := idToken.Claims(&IDTokenClaims); err != nil {
			authorizationFailed("claims : "+err.Error(), c) // An authorization error is returned in case the token is invalid.
			return
		}

		// We obtain the roles that are associated to the client in our case the roles are located in {"realm_access": {"roles": ["EDITOR", "USER"]}}}
		userAccessRoles := IDTokenClaims.RealmAccess.Roles
		for _, b := range userAccessRoles {
			// if the token contains the indicated role, you are allowed access.
			if b != "" {
				c.Next()
				return
			}
		}

		authorizationFailed("user not allowed to access this api", c) // An authorization error is returned in case the token is invalid.
	}
}
