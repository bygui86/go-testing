package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/bygui86/go-testing/rest-examples/http-server/commons"
	"github.com/bygui86/go-testing/rest-examples/http-server/database"
	"github.com/bygui86/go-testing/rest-examples/http-server/logging"
)

func (s *Server) GetProducts(writer http.ResponseWriter, request *http.Request) {
	span, ctx := retrieveSpanAndCtx(request, "get-products-handler")
	defer span.Finish()

	startTimer := time.Now()

	logging.Log.Info("Get products")

	span.SetTag("app", commons.ServiceName)

	count, _ := strconv.Atoi(request.FormValue("count"))
	start, _ := strconv.Atoi(request.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	products, err := s.db.GetProducts(start, count, ctx)
	if err != nil {
		errMsg := "Get products failed: " + err.Error()
		sendErrorResponse(writer, http.StatusInternalServerError, errMsg)

		span.SetTag("products-found", 0)
		span.SetTag("error", errMsg)
		span.LogKV("products-found", 0, "error", errMsg)
		return
	}

	span.SetTag("products-found", len(products))
	span.LogKV("products-found", len(products))

	sendJsonResponse(writer, http.StatusOK, products)

	IncreaseRestRequests("getProducts")
	ObserveRestRequestsTime("getProducts", float64(time.Now().Sub(startTimer).Milliseconds()))
}

func (s *Server) GetProduct(writer http.ResponseWriter, request *http.Request) {
	span, ctx := retrieveSpanAndCtx(request, "get-product-handler")
	defer span.Finish()

	startTimer := time.Now()

	span.SetTag("app", commons.ServiceName)

	vars := mux.Vars(request)
	id := vars["id"]

	logging.SugaredLog.Infof("Get product by ID: %s", id)

	span.SetTag("product-id", id)

	getProduct := &database.Product{ID: id}
	product := s.db.GetProduct(getProduct, ctx)
	if product == nil {
		errMsg := "Get product failed: product not found"
		sendErrorResponse(writer, http.StatusNotFound, errMsg)

		span.SetTag("product-found", false)
		span.SetTag("error", errMsg)
		span.LogKV("product-id", id, "product-found", false, "error", errMsg)
		return
	}

	span.SetTag("product-found", true)
	span.LogKV("product-id", id, "product-found", true)

	sendJsonResponse(writer, http.StatusOK, product)

	IncreaseRestRequests("getProduct")
	ObserveRestRequestsTime("getProduct", float64(time.Now().Sub(startTimer).Milliseconds()))
}

func (s *Server) CreateProduct(writer http.ResponseWriter, request *http.Request) {
	span, ctx := retrieveSpanAndCtx(request, "create-product-handler")
	defer span.Finish()

	startTimer := time.Now()

	span.SetTag("app", commons.ServiceName)

	var product *database.Product
	unmarshErr := json.NewDecoder(request.Body).Decode(&product)
	if unmarshErr != nil {
		errMsg := "Create product failed: invalid request payload"
		sendErrorResponse(writer, http.StatusBadRequest, errMsg)

		span.SetTag("product-created", false)
		span.SetTag("error", errMsg)
		span.LogKV("product-created", false, "error", errMsg)
		return
	}
	defer closeRequestBody(request.Body)

	logging.SugaredLog.Infof("Create product %s", product.String())

	createErr := s.db.CreateProduct(product, ctx)
	if createErr != nil {
		errMsg := "Create product failed: " + createErr.Error()
		sendErrorResponse(writer, http.StatusInternalServerError, errMsg)

		span.SetTag("product-created", false)
		span.SetTag("error", errMsg)
		span.LogKV("product-created", false, "error", errMsg)
		return
	}

	span.SetTag("product", product.String())
	span.SetTag("product-created", true)
	span.LogKV("product", product.String(), "product-created", true)

	sendJsonResponse(writer, http.StatusCreated, product)

	IncreaseRestRequests("createProduct")
	ObserveRestRequestsTime("createProduct", float64(time.Now().Sub(startTimer).Milliseconds()))
}

func (s *Server) UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	span, ctx := retrieveSpanAndCtx(request, "update-product-handler")
	defer span.Finish()

	startTimer := time.Now()

	span.SetTag("app", commons.ServiceName)

	vars := mux.Vars(request)
	id := vars["id"]

	var product *database.Product
	unmarshErr := json.NewDecoder(request.Body).Decode(&product)
	if unmarshErr != nil {
		errMsg := "Update product failed: invalid request payload"
		sendErrorResponse(writer, http.StatusBadRequest, errMsg)

		span.SetTag("product-updated", false)
		span.SetTag("error", errMsg)
		span.LogKV("product-updated", false, "error", errMsg)
		return
	}
	defer closeRequestBody(request.Body)

	product.ID = id
	logging.SugaredLog.Infof("Update product: %s", product.String())
	span.SetTag("product-id", id)

	updateErr := s.db.UpdateProduct(product, ctx)
	if updateErr != nil {
		errMsg := "Update product failed: " + updateErr.Error()
		sendErrorResponse(writer, http.StatusInternalServerError, errMsg)

		span.SetTag("product-updated", false)
		span.SetTag("error", errMsg)
		span.LogKV("product-updated", false, "error", errMsg)
		return
	}

	span.SetTag("product", product.String())
	span.SetTag("product-updated", true)
	span.LogKV("product", product.String(), "product-updated", true)

	sendJsonResponse(writer, http.StatusOK, product)

	IncreaseRestRequests("updateProduct")
	ObserveRestRequestsTime("updateProduct", float64(time.Now().Sub(startTimer).Milliseconds()))
}

func (s *Server) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	span, ctx := retrieveSpanAndCtx(request, "delete-product-handler")
	defer span.Finish()

	startTimer := time.Now()

	span.SetTag("app", commons.ServiceName)

	vars := mux.Vars(request)
	id := vars["id"]

	logging.SugaredLog.Infof("Delete product by ID: %s", id)
	span.SetTag("product-id", id)

	deleteErr := s.db.DeleteProduct(id, ctx)
	if deleteErr != nil {
		errMsg := "Delete product failed: " + deleteErr.Error()
		sendErrorResponse(writer, http.StatusInternalServerError, errMsg)

		span.SetTag("product-deleted", false)
		span.SetTag("error", errMsg)
		span.LogKV("product-deleted", false, "error", errMsg)
		return
	}

	span.SetTag("product-deleted", true)
	span.LogKV("product-deleted", true)

	sendJsonResponse(writer, http.StatusOK, map[string]string{"result": "success"})

	IncreaseRestRequests("deleteProduct")
	ObserveRestRequestsTime("deleteProduct", float64(time.Now().Sub(startTimer).Milliseconds()))
}
