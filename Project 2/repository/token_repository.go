package repository

import (
	"sync"

	"Project2/model"

	"github.com/pkg/errors"
)

type TokenRepository struct {
	sync.RWMutex
	tokens []*model.Token
}

func ProvideTokenRepository() *TokenRepository {
	return &TokenRepository{tokens: make([]*model.Token, 0)}
}

func (tr *TokenRepository) Create(id int64) {
	t := &model.Token{Id: id}
	tr.Lock()
	tr.tokens = append(tr.tokens, t)
	tr.Unlock()
}

func (tr *TokenRepository) Write(tk *model.Token) error {
	var token *model.Token
	for _, t := range tr.tokens {
		if t.Id == tk.Id {
			token = t
			break
		}
	}

	if token == nil {
		return errors.New("token not found")
	}

	if (tk.Low > tk.Mid) || (tk.Mid >= tk.High) {
		return errors.New("low, mid and high must satisfy low <= mid < high relation")
	}

	tr.Lock()
	token.Low = tk.Low
	token.High = tk.High
	token.Mid = tk.Mid
	token.Name = tk.Name
	token.PartialVal = tk.PartialVal
	token.FinalVal = tk.FinalVal
	tr.Unlock()

	return nil
}

func (tr *TokenRepository) Read(id int64) (*model.Token, error) {
	var token *model.Token
	for _, t := range tr.tokens {
		if t.Id == id {
			token = t
			break
		}
	}

	if token == nil {
		return nil, errors.New("token not found")
	}

	return token, nil
}

func (tr *TokenRepository) Drop(id int64) error {
	var token *model.Token
	var idx int

	for i, t := range tr.tokens {
		if t.Id == id {
			token = t
			idx = i
			break
		}
	}

	if token == nil {
		return errors.New("token not found")
	}

	tr.Lock()
	lastIdx := len(tr.tokens) - 1
	if lastIdx > 0 {
		tr.tokens[idx] = tr.tokens[lastIdx]
		tr.tokens = tr.tokens[0 : lastIdx-1]
	} else {
		tr.tokens = make([]*model.Token, 0)
	}
	tr.Unlock()

	return nil
}

func (tr *TokenRepository) GetAllTokens() []int64 {
	allTokens := make([]int64, 0)
	for _, t := range tr.tokens {
		allTokens = append(allTokens, t.Id)
	}

	return allTokens
}
