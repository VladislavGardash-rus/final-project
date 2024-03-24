package http_handlers

import (
	"context"
	"github.com/gardashvs/final-project/internal/services"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type PreviewerHandler struct {
	previewerService *services.PreviewerService
}

func NewPreviewerHandler() *PreviewerHandler {
	return &PreviewerHandler{previewerService: services.NewPreviewerService()}
}

func (h *PreviewerHandler) GetPreview(ctx context.Context, r *http.Request) (interface{}, error) {
	routeParams := mux.Vars(r)
	widthParam := routeParams["width"]
	heightParam := routeParams["height"]
	url := routeParams["url"]

	width, err := strconv.Atoi(widthParam)
	if err != nil {
		return nil, err
	}

	height, err := strconv.Atoi(heightParam)
	if err != nil {
		return nil, err
	}

	return h.previewerService.MakePreview(height, width, url, r.Header)
}
