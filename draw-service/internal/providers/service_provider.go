package providers

import (
	"github.com/MaxFando/lms/draw-service/internal/core/draw/service"
)

type ServiceProvider struct {
	provider *RepositoryProvider

	drawService *service.DrawService
}

func NewServiceProvider(provider *RepositoryProvider) *ServiceProvider {
	return &ServiceProvider{
		provider: provider,
	}
}

func (r *ServiceProvider) RegisterDependencies() {
	r.drawService = service.NewDrawService(r.provider.drawRepository)
}
