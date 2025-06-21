package httpres

import "time"

func New(success bool, data any) map[string]any {
	return map[string]any{
		"success":    success,
		"accessedAt": time.Now(),
		"data":       data,
	}
}
