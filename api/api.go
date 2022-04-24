package api

import (
	"errors"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	s.log = s.log.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return l.Core()
	}))
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

		s.log.With(zap.String("url", resp.Request.URL), zap.String("status", resp.Status())).Debug("query done")
		return qr, nil
	}

	s.log.With(zap.String("status", resp.Status()), zap.Int("code", resp.StatusCode())).Info("unknown response")

	return nil, ErrUnknown
}
