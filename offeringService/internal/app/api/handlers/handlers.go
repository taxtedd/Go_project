package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type BankHandlers interface {
	GetBalanceHandler(w http.ResponseWriter, r *http.Request)
	TopUpBalanceHandler(w http.ResponseWriter, r *http.Request)
	TransferMoneyHandler(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	Bank *bank.Bank
}

func NewHandler() *BankHandler {
	return &BankHandler{Bank: bank}
}

func (bankHandler *BankHandler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.New("invalid request method").Error(), http.StatusBadRequest)
		return
	}

	inputJSON, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}(r.Body)

	var input models.AccountBalanceRequest
	err = json.Unmarshal(inputJSON, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := bankHandler.Bank.FindAccount(input.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	outputJson, err := json.Marshal(models.AccountBalanceResponse{Balance: account.GetBalance(), Status: http.StatusOK})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(outputJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (bankHandler *BankHandler) TopUpBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.New("invalid request method").Error(), http.StatusBadRequest)
		return
	}

	inputJSON, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}(r.Body)

	var input models.TopUpBalanceRequest
	err = json.Unmarshal(inputJSON, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := bankHandler.Bank.FindAccount(input.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	bankHandler.Bank.TopUpAccountBalance(account, input.Amount)
	fmt.Println(bankHandler.Bank.Accountes)
	outputJson, err := json.Marshal(models.TopUpBalanceResponse{Status: http.StatusOK})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(outputJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (bankHandler *BankHandler) MoneyTransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusBadRequest)
		return
	}

	inputJSON, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}(r.Body)

	var input models.MoneyTransferRequest
	err = json.Unmarshal(inputJSON, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.FromId == input.ToId {
		http.Error(w, "you can not transfer money to yourself", http.StatusBadRequest)
		return
	}

	fromAccount, err := bankHandler.Bank.FindAccount(input.FromId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	toAccount, err := bankHandler.Bank.FindAccount(input.ToId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = bankHandler.Bank.WithdrawMoneyFromAccount(fromAccount, input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bankHandler.Bank.TopUpAccountBalance(toAccount, input.Amount)

	outputJson, err := json.Marshal(models.MoneyTransferResponse{Status: http.StatusOK})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(outputJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
