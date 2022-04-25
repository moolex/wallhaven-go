package api

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	ErrUnknown = errors.New("unknown error")
	ErrServer  = errors.New("server error")
)

func New(key string) *API {
	return &API{
		cli: resty.New().
			SetBaseURL("https://wallhaven.cc/api/v1").
			SetQueryParam("apikey", key),
		log: zap.NewNop(),
	}
}

type API struct {
	cli *resty.Client
	log *zap.Logger
}

func (s *API) SetDebug() {
	s.cli.SetDebug(true)
}

func (s *API) SetLogger(l *zap.Logger) {
	s.log = l
}

func (s *API) Query(qc *QueryCond) (*QueryResult, error) {
	params, err := qc.ToMap()
	if err != nil {
		return nil, err
	}

	result := &QueryResult{}

	resp, err := s.cli.R().
		SetQueryParams(params).
		SetResult(result).
		Get("/search")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, ErrServer
	}

	if qr, is := resp.Result().(*QueryResult); is {
		qr.api = s
		qr.cond = qc

		s.log.With(
			zap.String("status", resp.Status()),
			zap.Int("total", qr.Meta.Total),
			zap.String("page", fmt.Sprintf("%d/%d", qr.Meta.CurrentPage, qr.Meta.LastPage)),
			zap.String("ratelimit", fmt.Sprintf(
				"%s/%s",
				resp.Header().Get("x-ratelimit-remaining"),
				resp.Header().Get("x-ratelimit-limit"),
			)),
		).Debug("api query done")
		return qr, nil
	}

	s.log.With(zap.String("status", resp.Status()), zap.Int("code", resp.StatusCode())).Info("unknown response")

	return nil, ErrUnknown
}
