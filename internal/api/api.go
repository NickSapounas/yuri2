package api

import (
	"fmt"
	"net/http"

	"github.com/zekroTJA/discordgo"
	"github.com/zekroTJA/yuri2/internal/api/auth"
	"github.com/zekroTJA/yuri2/internal/config"
	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/internal/static"
	"github.com/zekroTJA/yuri2/pkg/discordoauth"
)

type API struct {
	cfg *config.API

	qualifiedAddress string

	db      database.Middleware
	session *discordgo.Session

	server *http.Server
	mux    *http.ServeMux
	auth   *auth.Auth

	discordAuthFE  *discordoauth.DiscordOAuth
	discordAuthAPI *discordoauth.DiscordOAuth
}

func NewAPI(cfg *config.API, db database.Middleware, session *discordgo.Session) *API {
	api := &API{
		cfg:     cfg,
		db:      db,
		session: session,
	}

	protocol := "http"
	address := api.cfg.Address
	if api.cfg.TLS != nil && api.cfg.TLS.Enable {
		protocol = "https"
	}
	if address[0] == ':' {
		address = "localhost" + address
	}
	api.qualifiedAddress = fmt.Sprintf("%s://%s", protocol, address)

	api.mux = http.NewServeMux()

	api.server = &http.Server{
		Handler: api.mux,
		Addr:    api.cfg.Address,
	}

	api.auth = auth.NewAuth(db, static.TokenHashRounds, static.TokenLifetime)

	api.discordAuthAPI = discordoauth.NewDiscordOAuth(
		api.cfg.ClientID,
		api.cfg.ClientSecret,
		api.qualifiedAddress+"/token/authorize",
		nil, //  TODO: func(w http.ResponseWriter, r *http.ReadRequest), errResponse
		api.getTokenHandler)

	api.discordAuthFE = discordoauth.NewDiscordOAuth(
		api.cfg.ClientID,
		api.cfg.ClientSecret,
		api.qualifiedAddress+"/login/authorize",
		errPageResponse,
		nil)

	api.registerHandlers()

	return api
}

func (api *API) registerHandlers() {
	// Static file server
	api.mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./web/static"))))

	// GET /token
	api.mux.HandleFunc("/token", api.discordAuthAPI.HandlerInit)

	// GET /token/authorize
	api.mux.HandleFunc("/token/authorize", api.discordAuthAPI.HandlerCallback)

	// GET /login
	api.mux.HandleFunc("/login", api.discordAuthFE.HandlerInit)

	// GET /login/authorize
	api.mux.HandleFunc("/login/authorize", api.discordAuthFE.HandlerCallback)
}

func (api *API) StartBlocking() error {
	var err error

	if api.cfg.TLS != nil && api.cfg.TLS.Enable {
		err = api.server.ListenAndServeTLS(api.cfg.TLS.CertFile, api.cfg.TLS.KeyFile)
	} else {
		err = api.server.ListenAndServe()
	}

	return err
}

func (api *API) Close() {
	if api == nil {
		return
	}
	api.server.Close()
}
