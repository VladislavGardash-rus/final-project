package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gardashvs/final-project/cfg"
	"github.com/gardashvs/final-project/internal/services/lru_cach_service"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"net/http"
	"time"
)

type PreviewerService struct {
	cache lru_cach_service.Cache
}

func NewPreviewerService() *PreviewerService {
	return &PreviewerService{cache: lru_cach_service.NewCache(cfg.Config().CacheCapacity)}
}

func (s *PreviewerService) MakePreview(height, width int, url string, headers http.Header) ([]byte, error) {
	image, found := s.findImageInCache(url)
	if found {
		image, err := s.resizeImage(image, height, width)
		if err != nil {
			return nil, err
		}

		return image, nil
	}

	image, err := s.getImageByApi(url, headers)
	if err != nil {
		return nil, err
	}
	s.saveImageInCache(url, image)

	image, err = s.resizeImage(image, height, width)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *PreviewerService) findImageInCache(url string) ([]byte, bool) {
	item, found := s.cache.Get(s.getImageKey(url))

	if image, ok := item.([]byte); ok {
		return image, found
	}

	return nil, found
}

func (s *PreviewerService) saveImageInCache(url string, image []byte) {
	s.cache.Set(s.getImageKey(url), image)
}

func (s *PreviewerService) getImageKey(url string) lru_cach_service.Key {
	hash := md5.Sum([]byte(url))
	return lru_cach_service.Key(hex.EncodeToString(hash[:]))
}

func (s *PreviewerService) getImageByApi(url string, headers http.Header) ([]byte, error) {
	client := new(http.Client)
	client.Timeout = 1 * time.Minute
	req, err := http.NewRequest("GET", "https://"+url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Get image by url error, resource answered: " + res.Status)
	}

	image, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return image, nil
}

func (s *PreviewerService) resizeImage(image []byte, height, width int) ([]byte, error) {
	reader := bytes.NewReader(image)
	img, err := jpeg.Decode(reader)
	if err != nil {
		return nil, err
	}

	if cfg.Config().ThumbnailMode {
		img = resize.Thumbnail(uint(width), uint(height), img, resize.NearestNeighbor)
	} else {
		img = resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
