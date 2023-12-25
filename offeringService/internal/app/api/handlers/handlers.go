package handlers

import (
	"Go_project/internal/app/api/requests"
	"Go_project/internal/app/service"
	"Go_project/pkg/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type OfferingHandler struct {
	Service *service.OfferingService
	Logger  *zap.Logger
	Server  *http.Server
}

func NewHandler(logger *zap.Logger) *OfferingHandler {
	offerService := service.NewService()
	handler := OfferingHandler{Service: offerService, Logger: logger}

	router := chi.NewRouter()
	router.Post("/offers", handler.CreateOffer)
	router.Get("/offers/{offerID}", handler.ParseOffer)

	handler.Server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &handler
}

func (handler *OfferingHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		handler.Logger.Error("Error reading request body", zap.Error(err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			handler.Logger.Error("Error closing request body", zap.Error(err))
		}
	}(r.Body)

	var offerRequest requests.CreateOfferRequest
	err = json.Unmarshal(bytes, &offerRequest)
	if err != nil {
		handler.Logger.Error("Error unmarshalling JSON", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	offer := models.Offer{From: offerRequest.From, To: offerRequest.To, ClientId: offerRequest.ClientId}
	offer.Price = *handler.Service.GetPrice(offer.From, offer.To)

	encodedJwt, err := handler.Service.EncodeJwt(&offer)
	if err != nil {
		handler.Logger.Error("Error encoding JWT", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := requests.CreateOfferResponse{OfferId: encodedJwt}
	res, err := json.Marshal(response)
	if err != nil {
		handler.Logger.Error("Error marshalling JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		handler.Logger.Error("Error writing response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *OfferingHandler) ParseOffer(w http.ResponseWriter, r *http.Request) {
	offerID := chi.URLParam(r, "offerID")

	decodedOffer, err := handler.Service.DecodeJwt(offerID)
	if err != nil {
		handler.Logger.Error("Error decoding JWT", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := requests.ParseOfferResponse{From: decodedOffer.From, To: decodedOffer.To, Price: decodedOffer.Price, ClientId: decodedOffer.ClientId}
	res, err := json.Marshal(response)
	if err != nil {
		handler.Logger.Error("Error marshalling JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		handler.Logger.Error("Error writing response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
