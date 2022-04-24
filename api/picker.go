package api

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrNoMoreItems = errors.New("no more items")
)

func (qr *QueryResult) Pick() (*Wallpaper, error) {
	if qr.Meta.Total == 0 || len(qr.Data) == 0 {
		return nil, ErrNoMoreItems
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

	qr.api.log.With(zap.String("id", w.Id)).Debug("picked wallpaper")
	return w, nil
}

func (qr *QueryResult) loadNext() error {
	if qr.Meta.CurrentPage == qr.Meta.LastPage {
		return ErrNoMoreItems
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
