package usecase

import (
	"github.com/MaxFando/lms/draw-service/internal/core/draw/service"
)

type DrawUseCase struct {
	crudService     *service.DrawService
	producerService *service.DrawProducer
}

func NewDrawUseCase(service *service.DrawService, producerService *service.DrawProducer) *DrawUseCase {
	return &DrawUseCase{crudService: service, producerService: producerService}
}
