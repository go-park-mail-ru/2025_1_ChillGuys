package tests

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestWithLogger(t *testing.T) {
	ctx := context.Background()
	logger := logrus.NewEntry(logrus.New())

	ctx = logctx.WithLogger(ctx, logger)

	retrieved, ok := ctx.Value(domains.LoggerKey{}).(*logrus.Entry)
	assert.True(t, ok, "Should retrieve logger from context")
	assert.Equal(t, logger, retrieved, "Retrieved logger should match stored one")
}

func TestGetLogger_WithLogger(t *testing.T) {
	ctx := context.Background()
	logger := logrus.NewEntry(logrus.New())

	ctx = logctx.WithLogger(ctx, logger)
	retrieved := logctx.GetLogger(ctx)

	assert.Equal(t, logger, retrieved, "Should retrieve stored logger")
}

func TestGetLogger_NoLogger(t *testing.T) {
	ctx := context.Background()
	retrieved := logctx.GetLogger(ctx)

	assert.NotNil(t, retrieved, "Should return non-nil logger when no logger in context")
	assert.IsType(t, &logrus.Entry{}, retrieved, "Should return logrus.Entry")
}
