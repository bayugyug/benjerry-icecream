package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bayugyug/benjerry-icecream/config"
	"github.com/bayugyug/benjerry-icecream/driver"
	"github.com/bayugyug/benjerry-icecream/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

const (
	svcOptionWithHandler  = "svc-opts-handler"
	svcOptionWithAddress  = "svc-opts-address"
	svcOptionWithDbConfig = "svc-opts-db-config"
	svcOptionWithDumpFile = "svc-opts-dump-file"
)

var ApiInstance *ApiService

type ApiService struct {
	Router       *chi.Mux
	DumpFile     string
	Address      string
	Api          *ApiHandler
	DB           *sql.DB
	Context      context.Context
	DriverConfig *driver.DbConnectorConfig
	AppJwtToken  *utils.AppJwtConfig
}

//WithSvcOptHandler opts for handler
func WithSvcOptHandler(r *ApiHandler) *config.Option {
	return config.NewOption(svcOptionWithHandler, r)
}

//WithSvcOptAddress opts for port#
func WithSvcOptAddress(r string) *config.Option {
	return config.NewOption(svcOptionWithAddress, r)
}

//WithSvcOptDbConf opts for db connector
func WithSvcOptDbConf(r *driver.DbConnectorConfig) *config.Option {
	return config.NewOption(svcOptionWithDbConfig, r)
}

//WithSvcOptDumpFile opts for port#
func WithSvcOptDumpFile(r string) *config.Option {
	return config.NewOption(svcOptionWithDumpFile, r)
}

//NewApiService service new instance
func NewApiService(opts ...*config.Option) (*ApiService, error) {

	//default
	svc := &ApiService{
		Address:     ":8989",
		Api:         &ApiHandler{},
		Context:     context.Background(),
		AppJwtToken: utils.NewAppJwtConfig(),
	}

	//add options if any
	for _, o := range opts {
		//chk opt-name
		switch o.Name() {
		case svcOptionWithHandler:
			if s, oks := o.Value().(*ApiHandler); oks && s != nil {
				svc.Api = s
			}
		case svcOptionWithAddress:
			if s, oks := o.Value().(string); oks && s != "" {
				svc.Address = s
			}
		case svcOptionWithDumpFile:
			if s, oks := o.Value().(string); oks && s != "" {
				svc.DumpFile = s
			}
		case svcOptionWithDbConfig:
			if s, oks := o.Value().(*driver.DbConnectorConfig); oks && s != nil {
				svc.DriverConfig = s
			}
		}
	} //iterate all opts

	//set the actual router
	svc.Router = svc.MapRoute()

	//get db
	dbh, err := driver.NewDbConnector(svc.DriverConfig)
	if err != nil {
		return svc, err
	}

	//save
	svc.DB = dbh
	return svc, nil
}

//MapRoute route map all endpoints
func (svc *ApiService) PrepareData() int {
	log.Println(svc.DumpFile)
	if svc.DumpFile != "" {
		tot, err := svc.Api.Preload(svc.DumpFile)
		if err != nil {
			log.Println("Prepare data:", err)
			return -1
		}
		return tot
	}
	return 0
}

//Run run the http server based on settings
func (svc *ApiService) Run() {

	//gracious timing
	srv := &http.Server{
		Addr:         svc.Address,
		Handler:      svc.Router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//async run
	go func() {
		log.Println("Listening on port", svc.Address)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
			os.Exit(0)
		}

	}()

	//watcher
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	log.Println("Shutting down service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped!")
}

//MapRoute route map all endpoints
func (svc *ApiService) MapRoute() *chi.Mux {

	// Multiplexer
	router := chi.NewRouter()

	// Basic settings
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)

	// Basic gracious timing
	router.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	router.Use(cors.Handler)

	router.Get("/", svc.Api.IndexPage)

	/*
		@users
		POST     /v1/api/user
		PUT      /v1/api/user
		GET      /v1/api/user/{id}
		DELETE   /v1/api/user/{id}
		POST     /v1/api/otp
		POST     /v1/api/login



		@Icecream
		POST     /v1/api/icecream
		PUT      /v1/api/icecream
		GET      /v1/api/icecream/{id}
		DELETE   /v1/api/icecream/{id}

		@ingredients
		POST     /v1/api/ingredient
		PUT      /v1/api/ingredient
		GET      /v1/api/ingredient/{id}
		DELETE   /v1/api/ingredient/{id}


		@sourcing_values
		POST     /v1/api/sourcing
		PUT      /v1/api/sourcing
		GET      /v1/api/sourcing/{id}
		DELETE   /v1/api/sourcing/{id}


	*/

	// Protected routes
	router.Route("/v1", func(r chi.Router) {
		r.Use(svc.SetContextKeyVal("api.version", "v1"))
		r.Mount("/api/user",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(svc.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Put("/", api.UpdateUser)
				sr.Get("/{id}", api.GetUser)
				sr.Delete("/{id}", api.DeleteUser)
				return sr
			}(svc.Api))
		r.Mount("/api/icecream",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(svc.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Post("/", api.CreateIcecream)
				sr.Put("/{id}", api.UpdateIcecream)
				sr.Get("/{id}", api.GetIcecream)
				sr.Delete("/{id}", api.DeleteIcecream)
				return sr
			}(svc.Api))
		r.Mount("/api/ingredient",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(svc.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Post("/{id}", api.CreateIngredient)
				sr.Get("/{id}", api.GetIngredient)
				sr.Delete("/{id}", api.DeleteIngredient)
				return sr
			}(svc.Api))
		r.Mount("/api/sourcing",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(svc.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Post("/{id}", api.CreateSourcing)
				sr.Get("/{id}", api.GetSourcing)
				sr.Delete("/{id}", api.DeleteSourcing)
				return sr
			}(svc.Api))
		//not-yet implemented ;-)
		r.Mount("/api/logout",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Use(jwtauth.Verifier(svc.AppJwtToken.TokenAuth), svc.BearerChecker)
				sr.Post("/", api.Logout)
				return sr
			}(svc.Api))

	})

	//public-routes
	router.Group(func(r chi.Router) {
		r.Post("/v1/api/login", svc.Api.Login)
		r.Post("/v1/api/user", svc.Api.CreateUser)
		r.Post("/v1/api/otp", svc.Api.Otp)
	})
	return router
}

//SetContextKeyVal version context
func (svc *ApiService) SetContextKeyVal(k, v string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), k, v))
			next.ServeHTTP(w, r)
		})
	}
}

//BearerChecker check token
func (svc *ApiService) BearerChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			switch err {
			default:
				log.Println("ERROR:", err)
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			case jwtauth.ErrExpired:
				log.Println("ERROR: Expired")
				http.Error(w, "Expired", http.StatusUnauthorized)
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			case jwtauth.ErrUnauthorized:
				log.Println("ERROR: ErrUnauthorized")
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
		}

		if token == nil || !token.Valid {
			svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})

}
