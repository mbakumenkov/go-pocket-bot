package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/mbakumenkov/go-pocket-bot/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectUrl     string
}

func NewAuthorizationServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectUrl string) *AuthorizationServer {
	return &AuthorizationServer{
		pocketClient:    pocketClient,
		tokenRepository: tokenRepository,
		redirectUrl:     redirectUrl,
	}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}
	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIdParam := r.URL.Query().Get("chat_id")
	if chatIdParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := strconv.ParseInt(chatIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.Get(chatId, repository.RequestTokens)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResp, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.tokenRepository.Save(chatId, authResp.AccessToken, repository.AccessTokens)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("chat_id: %s\nrequest_token: %s\naccess_token:%s\n", chatIdParam, requestToken, authResp.AccessToken)

	w.Header().Add("Location", s.redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}
