package tui

import (
	"image"

	"github.com/zMoooooritz/nachrichten/pkg/http"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
)

type ImageCache struct {
	images map[string]image.Image
}

func NewImageCache() *ImageCache {
	ic := ImageCache{
		images: make(map[string]image.Image),
	}
	return &ic
}

func (ic *ImageCache) LoadImage(id, url string) error {
	if _, found := ic.images[id]; found {
		return nil
	}
	img, err := http.LoadImage(url)
	if err == nil {
		ic.images[id] = img
	}
	return err
}

func (ic *ImageCache) GetImage(id, url string) image.Image {
	if img, found := ic.images[id]; found {
		return img
	}
	err := ic.LoadImage(id, url)
	if err == nil {
		return ic.images[id]
	}
	return image.Rect(0, 0, 1, 1)
}

func (ic *ImageCache) LoadThumbnails(articles []tagesschau.Article) {
	for _, article := range articles {
		ic.LoadThumbnail(article)
	}
}

func (ic *ImageCache) LoadThumbnail(article tagesschau.Article) {
	imageSpec := tagesschau.ImageSpec{Size: tagesschau.SMALL, Ratio: tagesschau.RECT}
	_ = ic.LoadImage(article.ID, tagesschau.GetImageURL(article.ImageData.ImageVariants, imageSpec))
}
