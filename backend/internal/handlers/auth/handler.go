package handlers

import (
	"net/http"

	"github.com/diagnosis/deploywatchv2/internal/config"
	refreshtokenstore "github.com/diagnosis/deploywatchv2/internal/store/refreshtoken"
	userstore "github.com/diagnosis/deploywatchv2/internal/store/user"
	"github.com/diagnosis/go-toolkit/secure"
)

type AuthHandler struct {
	cfg               *config.Config
	jwt               *secure.JWTSigner
	userStore         userstore.UserStore
	refreshTokenStore refreshtokenstore.RefreshTokenStore
}

func NewAuthHandler(
	cfg *config.Config,
	jwt *secure.JWTSigner,
	userStore userstore.UserStore,
	refreshTokenStore refreshtokenstore.RefreshTokenStore,
) *AuthHandler {
	return &AuthHandler{
		cfg:               cfg,
		jwt:               jwt,
		userStore:         userStore,
		refreshTokenStore: refreshTokenStore,
	}
}

func (h *AuthHandler) AuthFunc(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}
	claims, err := h.jwt.VerifyAccess(cookie.Value)
	if err != nil {
		return "", err
	}
	return claims.Sub, nil
}
