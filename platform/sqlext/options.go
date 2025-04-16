package sqlext

import (
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type ConnOption func(*config) error

// WithMaxConns устанавливает максимальные значения для числа открытых и простаивающих соединений в настройках config.
// Возвращает ошибку, если лимит простаивающих соединений превышает лимит открытых соединений.
func WithMaxConns(idle, open int) ConnOption {
	return func(c *config) error {
		if idle > open {
			return fmt.Errorf("ожидаемое количество простаивающих соединений не может быть больше, чем открытых (%v, %v)", idle, open)
		}

		c.maxIdleConns, c.maxOpenConns = idle, open
		return nil
	}
}

// WithConnTime задает время жизни соединения и время его бездействия.
// Возвращает ошибку, если значения отрицательные.
func WithConnTime(life, idle time.Duration) ConnOption {
	return func(c *config) error {
		if life < 0 {
			return fmt.Errorf("время жизни соединения не может быть отрицательным: %v", life)
		}
		if idle < 0 {
			return fmt.Errorf("время бездействия соединения не может быть отрицательным: %v", idle)
		}

		c.connLifeTime, c.connIdleTime = life, idle
		return nil
	}
}

// WithTracerProvider устанавливает провайдер трассировки в настройках config.
func WithTracerProvider(tracerProvider trace.TracerProvider) ConnOption {
	return func(c *config) error {
		c.tracerProvider = tracerProvider
		return nil
	}
}
