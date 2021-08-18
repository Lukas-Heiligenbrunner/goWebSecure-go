package oauth

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"log"
	"net/http"
)

var srv *server.Server

// callback function when a userlogin is performed to check the user in db
var userQuery func(id string) (CustomClientInfo, error)

func InitOAuth(userquery func(id string) (CustomClientInfo, error)) {
	userQuery = userquery

	manager := manage.NewDefaultManager()
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// create new secretstore
	clientStore := NewCustomStore()
	manager.MapClientStorage(clientStore)

	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func ValidateToken(f func(rw http.ResponseWriter, req *http.Request, node int, tokenInfo *oauth2.TokenInfo), node int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokeninfo, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		f(w, r, node, &tokeninfo)
	}
}
