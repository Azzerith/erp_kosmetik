package cron

import (
	"context"
	"log"
	"time"

	"erp-cosmetics-backend/internal/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type TrendCron struct {
	trendService service.TrendService
	logger       *zap.Logger
	cron         *cron.Cron
}

func NewTrendCron(trendService service.TrendService, logger *zap.Logger) *TrendCron {
	return &TrendCron{
		trendService: trendService,
		logger:       logger,
		cron:         cron.New(cron.WithLocation(time.Local)),
	}
}

func (tc *TrendCron) Start() {
	// Run every 6 hours: 0 */6 * * *
	_, err := tc.cron.AddFunc("0 */6 * * *", func() {
		tc.refreshTrends()
	})
	if err != nil {
		log.Printf("Failed to add trend cron job: %v", err)
	}

	tc.cron.Start()
	tc.logger.Info("Trend cron job started")

	// Run immediately on start
	go tc.refreshTrends()
}

func (tc *TrendCron) Stop() {
	tc.cron.Stop()
	tc.logger.Info("Trend cron job stopped")
}

func (tc *TrendCron) refreshTrends() {
	ctx := context.Background()
	tc.logger.Info("Starting trend refresh scheduled job")

	if err := tc.trendService.RefreshTrends(ctx); err != nil {
		tc.logger.Error("Failed to refresh trends", zap.Error(err))
	} else {
		tc.logger.Info("Trend refresh completed successfully")
	}
}