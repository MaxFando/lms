package providers

type ServiceProvider struct {
	provider *RepositoryProvider
}

func NewServiceProvider(provider *RepositoryProvider) *ServiceProvider {
	return &ServiceProvider{
		provider: provider,
	}
}

func (r *ServiceProvider) RegisterDependencies() {
}
