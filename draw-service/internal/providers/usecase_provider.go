package providers

type UsecaseProvider struct {
	provider *ServiceProvider
}

func NewUsecaseProvider(provider *ServiceProvider) *UsecaseProvider {
	return &UsecaseProvider{
		provider: provider,
	}
}

func (r *UsecaseProvider) RegisterDependencies() {
}
