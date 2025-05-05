//go:generate mockgen -source=$GOFILE -destination=./mock_${GOPACKAGE}_test.go -package=${GOPACKAGE}
package service

import (
	"context"
	"github.com/MaxFando/lms/draw-service/internal/core/draw/entity"
)

type Producer interface {
	Produce(ctx context.Context, draw *entity.Draw) error
}

type DrawProducer struct {
}
