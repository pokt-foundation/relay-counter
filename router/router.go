package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pokt-foundation/relay-counter/types"
	jsonresponse "github.com/pokt-foundation/utils-go/json-response"
	"golang.org/x/sync/errgroup"

	"github.com/sirupsen/logrus"
)

type Driver interface {
	WriteRelayCount(ctx context.Context, count types.RelayCount) error
	ReadRelayCounts(ctx context.Context, from, to time.Time) ([]types.RelayCount, error)
}

type Router struct {
	router  *mux.Router
	driver  Driver
	apiKeys map[string]bool
	port    string
	log     *logrus.Logger
}

func (rt *Router) logError(err error) {
	fields := logrus.Fields{
		"err": err.Error(),
	}

	rt.log.WithFields(fields).Error(err)
}

func respondWithResultOK(w http.ResponseWriter) {
	jsonresponse.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "ok"})
}

// NewRouter returns router instance
func NewRouter(driver Driver, apiKeys map[string]bool, port string, logger *logrus.Logger) (*Router, error) {
	rt := &Router{
		driver:  driver,
		router:  mux.NewRouter(),
		apiKeys: apiKeys,
		port:    port,
		log:     logger,
	}

	rt.router.HandleFunc("/", rt.HealthCheck).Methods(http.MethodGet)

	rt.router.HandleFunc("/v0/count", rt.CreateCount).Methods(http.MethodPost)

	rt.router.Use(rt.AuthorizationHandler)

	return rt, nil
}

func (rt *Router) RunServer(ctx context.Context) {
	httpServer := &http.Server{
		Addr:    ":" + rt.port,
		Handler: rt.router,
	}

	rt.log.Printf("Relay Counter is running in port: %s", rt.port)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		rt.log.Info("HTTP router context finished")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			rt.logError(fmt.Errorf("Error closing http server: %s", err))
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		rt.log.Infof("exit reason: %s", err.Error())
	}
}

func (rt *Router) AuthorizationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is the path of the health check endpoint
		if r.URL.Path == "/" {
			h.ServeHTTP(w, r)

			return
		}

		if !rt.apiKeys[r.Header.Get("Authorization")] {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unauthorized"))
			if err != nil {
				panic(err)
			}

			return
		}

		h.ServeHTTP(w, r)
	})
}

func (rt *Router) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Transaction HTTP DB is up and running!"))
	if err != nil {
		panic(err)
	}
}

func (rt *Router) CreateCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	var count types.RelayCount
	err := decoder.Decode(&count)
	if err != nil {
		rt.logError(fmt.Errorf("WriteCount in JSON decoding failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	err = rt.driver.WriteRelayCount(ctx, count)
	if err != nil {
		rt.logError(fmt.Errorf("WriteCount in WriteSession failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithResultOK(w)
}
