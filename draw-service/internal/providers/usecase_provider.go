package providers

import (
	"github.com/MaxFando/lms/draw-service/internal/core/draw/usecase"
)

type UsecaseProvider struct {
	provider *ServiceProvider

	drawUseCase *usecase.DrawUseCase
}

func NewUsecaseProvider(provider *ServiceProvider) *UsecaseProvider {
	return &UsecaseProvider{
		provider: provider,
	}
}

func (r *UsecaseProvider) RegisterDependencies() {
	r.drawUseCase = usecase.NewDrawUseCase(r.provider.drawService)
}
