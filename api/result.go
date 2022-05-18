package api

import (
	"errors"
	"math/rand"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

var (
	ErrNoSuchItems = errors.New("no such items")
	ErrNoMoreItems = errors.New("no more items")
)

type PickOption = func(qr *QueryResult)

func PickLoop(qr *QueryResult) {
	qr.pickLoop = true
}

func PickRand(qr *QueryResult) {
	qr.pickRand.Do(func() {
		rand.Seed(time.Now().UnixNano())
		lo.Shuffle(qr.Data)
	})
}

func (qr *QueryResult) Index() int {
	return qr.pIdx
}

func (qr *QueryResult) Pick(opts ...PickOption) (*Wallpaper, error) {
	for _, opt := range opts {
		opt(qr)
	}

	if qr.Meta.Total == 0 || len(qr.Data) == 0 {
		return nil, ErrNoSuchItems
	}

	qr.pLock.Lock()
	defer qr.pLock.Unlock()

	size := len(qr.Data)
	if qr.pIdx > size-1 {
		if err := qr.loadNext(); err != nil {
			return nil, err
		}
	}

	w := qr.Data[qr.pIdx]
	qr.pIdx++

	return w, nil
}

func (qr *QueryResult) loadNext() error {
	if qr.Meta.CurrentPage == qr.Meta.LastPage {
		if qr.pickLoop {
			qr.cond.Page = 0
		} else {
			return ErrNoMoreItems
		}
	}

	qc := qr.cond
	qc.Page++

	qr.api.log.With(zap.Int("page", qc.Page)).Debug("try loading next page")

	ret, err := qr.api.Query(qc)
	if err != nil {
		return err
	}

	qr.pIdx = 0
	qr.Data = ret.Data
	qr.Meta = ret.Meta

	return nil
}
