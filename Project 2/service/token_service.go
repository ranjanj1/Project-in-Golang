package service

import (
	"encoding/json"
	"fmt"
	"log"

	"Project2/model"
	"Project2/repository"
)

type TokenService struct {
	tokenRepo *repository.TokenRepository
	hashSvc   HashService
}

func ProvideTokenService(tokenRepo *repository.TokenRepository, hashSvc HashService) *TokenService {
	return &TokenService{tokenRepo: tokenRepo, hashSvc: hashSvc}
}

func (ts TokenService) Create(id int64) {
	ts.tokenRepo.Create(id)
	log.Println(fmt.Sprintf("token created successfully with id: %d", id))
	ts.dump(nil)
}

func (ts TokenService) Write(token *model.Token) error {
	token.PartialVal = ts.hashSvc.FindArgMin(token.Name, token.Low, token.Mid-1)

	err := ts.tokenRepo.Write(token)
	if err != nil {
		return err
	}
	ts.dump(token)
	return nil
}

func (ts TokenService) Read(id int64) (*model.Token, error) {
	tk, err := ts.tokenRepo.Read(id)
	if err != nil {
		return nil, err
	}

	partialVal := ts.hashSvc.FindArgMin(tk.Name, tk.Mid, tk.High-1)
	if partialVal > tk.PartialVal {
		tk.FinalVal = tk.PartialVal
	} else {
		tk.FinalVal = partialVal
	}
	ts.dump(tk)
	return tk, nil
}

func (ts TokenService) Drop(id int64) error {
	err := ts.tokenRepo.Drop(id)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("token dropped successfully with id: %d", id))
	ts.dump(nil)
	return nil
}

func (ts *TokenService) dump(token *model.Token) {
	if token != nil {
		b, _ := json.Marshal(token)
		log.Printf("current token info: %s\n", string(b))
	}
	log.Printf("tokens in store: %v\n", ts.tokenRepo.GetAllTokens())
}
