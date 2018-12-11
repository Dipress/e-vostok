package send

import (
	"github.com/pkg/errors"
)

type GetBody interface {
	Get(url string) (string, error)
}

type Sender interface {
	Send(body string, to []string) error
}

type Service struct {
	GetBody
	Sender
}

func NewService(g GetBody, s Sender) *Service {
	srv := Service{
		GetBody: g,
		Sender:  s,
	}

	return &srv
}

func (s *Service) Deliver(url string, to []string) error {
	body, err := s.GetBody.Get(url)
	if err != nil {
		return errors.Wrap(err, "get body")
	}

	if err := s.Sender.Send(body, to); err != nil {
		return errors.Wrap(err, "send mail")
	}

	return nil
}
