package api

import (
	"errors"
	"strconv"

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

		rl, _ := strconv.ParseInt(resp.Header().Get("x-ratelimit-limit"), 10, 64)
		rr, _ := strconv.ParseInt(resp.Header().Get("x-ratelimit-remaining"), 10, 64)
		qr.RateLimitQuota = int(rl)
		qr.RateLimitRemain = int(rr)

		return qr, nil
	}

	s.log.With(zap.String("status", resp.Status()), zap.Int("code", resp.StatusCode())).Info("unknown response")

	return nil, ErrUnknown
}
