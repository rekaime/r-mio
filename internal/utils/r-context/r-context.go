package rcontext

import (
    "context"
    "time"
)

func CreateTimeoutContext() (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), 10 * time.Second)
}