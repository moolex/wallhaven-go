package utils

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"

	"github.com/moolex/wallhaven-go/api"
)

var d = resty.New().SetDoNotParseResponse(true)

func GetThumbImage(wp *api.Wallpaper, size string) (image.Image, error) {
	url := lo.
		If(size == api.ThumbSmall, wp.Thumbs.Small).
		ElseIf(size == api.ThumbLarge, wp.Thumbs.Large).
		Else(wp.Thumbs.Original)

	resp, err := d.R().Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.RawBody().Close()

	img, _, err2 := image.Decode(resp.RawBody())
	return img, err2
}
