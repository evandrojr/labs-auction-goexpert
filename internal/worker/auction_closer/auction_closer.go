package auction_closer

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"time"
)

type AuctionCloserWorker struct {
	AuctionRepository *auction.AuctionRepository
	Interval          time.Duration
}

func NewAuctionCloserWorker(
	auctionRepository *auction.AuctionRepository,
	interval time.Duration) *AuctionCloserWorker {
	return &AuctionCloserWorker{
		AuctionRepository: auctionRepository,
		Interval:          interval,
	}
}

func (acw *AuctionCloserWorker) Start() {
	logger.Info("AuctionCloserWorker started")
	ticker := time.NewTicker(acw.Interval)
	defer ticker.Stop()

	for range ticker.C {
		logger.Info("Checking for expired auctions...")
		ctx := context.Background()
		auctions, err := acw.AuctionRepository.FindExpiredAuctions(ctx)
		if err != nil {
			logger.Error("Error finding expired auctions", err)
			continue
		}

		if len(auctions) == 0 {
			logger.Info("No expired auctions found.")
			continue
		}

		logger.Info("Found expired auctions. Closing them...")
		for _, auction := range auctions {
			err := acw.AuctionRepository.UpdateAuctionStatus(
				ctx, auction.Id, auction_entity.Completed)
			if err != nil {
				logger.Error("Error closing auction", err)
			}
		}
		logger.Info("Finished closing expired auctions.")
	}
}
