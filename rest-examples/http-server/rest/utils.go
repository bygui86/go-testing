package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bygui86/go-testing/rest-examples/http-server/commons"
	"github.com/bygui86/go-testing/rest-examples/http-server/logging"
)

const (
	// urls
	rootProductsEndpoint = "/products"
	productsIdEndpoint   = rootProductsEndpoint + "/{id:[0-9]+}"

	contentTypeHeaderKey       = "Content-Type"
	contentTypeApplicationJson = "application/json"
)

// SERVER

func (s *Server) setupRouter() {
	logging.Log.Debug("Create new router")

	s.router = mux.NewRouter().StrictSlash(true)

	s.router.Use(requestInfoPrintingMiddleware)

	s.router.HandleFunc(rootProductsEndpoint, s.GetProducts).Methods(http.MethodGet)
	s.router.HandleFunc(productsIdEndpoint, s.GetProduct).Methods(http.MethodGet)
	s.router.HandleFunc(rootProductsEndpoint, s.CreateProduct).Methods(http.MethodPost)
	s.router.HandleFunc(productsIdEndpoint, s.UpdateProduct).Methods(http.MethodPut)
	s.router.HandleFunc(productsIdEndpoint, s.DeleteProduct).Methods(http.MethodDelete)
}

func (s *Server) setupHTTPServer() {
	logging.SugaredLog.Debugf("Create new HTTP server on port %d", s.config.RestPort)

	if s.config != nil {
		s.httpServer = &http.Server{
			Addr:    fmt.Sprintf(commons.HttpServerHostFormat, s.config.RestHost, s.config.RestPort),
			Handler: s.router,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: commons.HttpServerWriteTimeoutDefault,
			ReadTimeout:  commons.HttpServerReadTimeoutDefault,
			IdleTimeout:  commons.HttpServerIdelTimeoutDefault,
		}
		return
	}

	logging.Log.Error("HTTP server creation failed: REST server configurations not loaded")
}

// HANDLERS

func sendJsonResponse(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	writer.Header().Set(contentTypeHeaderKey, contentTypeApplicationJson)
	writer.WriteHeader(code)
	_, err := writer.Write(response)
	if err != nil {
		logging.SugaredLog.Errorf("Error sending JSON response: %s", err.Error())
	}
}

func sendErrorResponse(writer http.ResponseWriter, code int, message string) {
	sendJsonResponse(writer, code, map[string]string{"error": message})
}

func closeRequestBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logging.SugaredLog.Errorf("Closing request body failed: %s", err.Error())
	}
}