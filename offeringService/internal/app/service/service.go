package service

import (
	"Go_project/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"log"
	"math"
	"time"
)

var jwtSecretKey = "uI/GxXcOxIMlAT+wMMrPmlK0kuIvPFm2O7lh5J5mQn8="

type OfferingService struct {
}

func NewService() *OfferingService {
	return &OfferingService{}
}

func getRadians(degree float64) float64 {
	return degree * math.Pi / 180
}

func getDistance(from models.Geolocation, to models.Geolocation) float64 {
	latFrom := getRadians(from.Lat)
	latTo := getRadians(to.Lat)
	lngFrom := getRadians(from.Lng)
	lngTo := getRadians(to.Lng)

	// D - расстояние между пунктами, измеряемое в радианах длиной дуги большого круга земного шара
	cosD := math.Sin(latFrom)*math.Sin(latTo) + math.Cos(latFrom)*math.Cos(latTo)*math.Cos(lngFrom-lngTo)
	distance := math.Acos(cosD) * Radius

	return math.Round(distance*10) / 10
}

func (offeringService *OfferingService) GetPrice(from models.Geolocation, to models.Geolocation) *models.Price {
	price := getDistance(from, to) * kmPrice
	price = price + startPrice

	return &models.Price{Amount: price, Currency: currency}
}

func (offeringService *OfferingService) EncodeJwt(offer *models.Offer) (string, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
	}
	offerJSON, err := json.Marshal(offer)
	if err != nil {
		logger.Error("Error signing JWT token", zap.Error(err))
		return "", err
	}

	payload := jwt.MapClaims{
		"sub": string(offerJSON),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		logger.Error("Error signing JWT token", zap.Error(err))
		return "", err
	}

	return t, err
}

func (offeringService *OfferingService) DecodeJwt(jwtToken string) (*models.Offer, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
	}
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		logger.Error("Error parsing JWT token", zap.Error(err))
		return nil, err
	}

	if !token.Valid {
		logger.Error("Invalid JWT token")
		return nil, fmt.Errorf("invalid jwt token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		offer := claims["sub"].(string)

		var offerJSON models.Offer
		err = json.Unmarshal([]byte(offer), &offerJSON)
		if err != nil {
			logger.Error("Error decoding offer from JSON", zap.Error(err))
			return nil, err
		}

		return &offerJSON, nil
	}

	return nil, fmt.Errorf("error in jwt token")
}
