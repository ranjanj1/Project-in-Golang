package mapper

import (
	"Project2/model"
	tokenv1 "Project2/proto"
)

type Mapper interface {
	Token(*tokenv1.Token) *model.Token
	TokenPb(*model.Token) *tokenv1.Token
}

type mapper struct{}

func ProvideMapper() Mapper {
	return &mapper{}
}

func (m *mapper) Token(t *tokenv1.Token) *model.Token {
	return &model.Token{
		Id:   t.Id,
		Name: t.Name,
		Low:  t.Low,
		Mid:  t.Mid,
		High: t.High,
	}
}

func (m *mapper) TokenPb(t *model.Token) *tokenv1.Token {
	return &tokenv1.Token{
		Id:           t.Id,
		Name:         t.Name,
		Low:          t.Low,
		Mid:          t.Mid,
		High:         t.High,
		PartialValue: t.PartialVal,
		FinalValue:   t.FinalVal,
	}
}
