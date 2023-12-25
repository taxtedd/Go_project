package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"client_service/internal/app/api/requests"
	"client_service/internal/app/config"
	"client_service/internal/app/models"
	"client_service/internal/app/mongodb"
	"github.com/go-chi/chi/v5"
)

var (
	ErrorCloseReqBody = errors.New("failed to close request body")
	ErrorReadReqBody  = errors.New("failed to read request body")
)

type ClientHandler struct {
	Server   *http.Server
	Config   *config.Config
	Database *mongodb.Database
}

func NewHandler(db *mongodb.Database, cfg *config.Config) *ClientHandler {
	handler := ClientHandler{Config: cfg, Database: db}

	router := chi.NewRouter()
	router.Post("/trips", handler.createTrip)
	router.Get("/trips", handler.listTrips)
	router.Get("/trips/{trip_id}", handler.getTrip)
	router.Post("/trip/{trip_id}/cancel", handler.cancelTrip)

	handler.Server = &http.Server{
		Addr:    handler.Config.Port,
		Handler: router,
	}

	return &handler
}

func (handler *ClientHandler) getTrip(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	tripID := chi.URLParam(r, "trip_id")
	tripData, err := handler.Database.GetTripByTripId(tripID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(tripData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (handler *ClientHandler) listTrips(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userID := r.Header.Get("user_id")
	tripsData, err := handler.Database.GetTripsByUserId(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(tripsData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (handler *ClientHandler) cancelTrip(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	tripID := chi.URLParam(r, "trip_id")
	if err := handler.Database.CancelTrip(tripID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ClientHandler) createTrip(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userID := r.Header.Get("user_id")
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(ErrorReadReqBody, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(ErrorCloseReqBody, err)
		}
	}(r.Body)

	var request requests.RequestCreateTrip
	err = json.Unmarshal(bytes, &request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	offer, err := handler.getOfferDetails(request.OfferId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tripData := &mongodb.Trip{
		OfferId:  request.OfferId,
		ClientId: userID,
		From: mongodb.Geolocation{
			Lat: offer.From.Lat,
			Lng: offer.From.Lng,
		},
		To: mongodb.Geolocation{
			Lat: offer.To.Lat,
			Lng: offer.To.Lng,
		},
		Price:  mongodb.Price{Amount: offer.Price.Amount, Currency: offer.Price.Currency},
		Status: "DRIVER_SEARCH",
	}

	if err := handler.Database.CreateTrip(tripData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ClientHandler) getOfferDetails(offerID string) (*models.Offer, error) {
	resp, err := http.Get("http://127.0.0.1:8080/offers/" + offerID)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(ErrorCloseReqBody, err)
		}
	}(resp.Body)

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var offer models.Offer
	err = json.Unmarshal(bytes, &offer)
	if err != nil {
		return nil, err
	}

	return &offer, nil
}
