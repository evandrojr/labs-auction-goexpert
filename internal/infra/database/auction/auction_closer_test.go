package auction_test

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/worker/auction_closer"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"testing"
	"time"
)

func TestAuctionCloserWorker(t *testing.T) {
	ctx := context.Background()

	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://admin:admin@localhost:27017/auctions?authSource=admin"
	}
	os.Setenv("MONGODB_URL", mongoURL)
	os.Setenv("MONGODB_DB", "auctions")
	client, dbClient, err := mongodb.NewMongoDBConnection(ctx)
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	auctionRepo := auction.NewAuctionRepository(dbClient)

	// Clear existing auctions for a clean test environment
	_, err = auctionRepo.Collection.DeleteMany(ctx, bson.M{})
	assert.NoError(t, err)

	// Create an auction that should expire soon
	expiresIn := 1 * time.Second
	auctionEntity, err := auction_entity.CreateAuction(
		"Test Product", "Electronics", "Description", auction_entity.New, expiresIn)
	assert.Nil(t, err)

	internalErr := auctionRepo.CreateAuction(ctx, auctionEntity)
	assert.Nil(t, internalErr)

	// Start the worker with a short interval
	workerInterval := 500 * time.Millisecond
	auctionCloser := auction_closer.NewAuctionCloserWorker(auctionRepo, workerInterval)
	go auctionCloser.Start()

	// Wait for the auction to expire and the worker to run
	time.Sleep(expiresIn + (5 * workerInterval))

	// Verify the auction status
	foundAuction, internalErr := auctionRepo.FindAuctionById(ctx, auctionEntity.Id)
	assert.Nil(t, internalErr)
	assert.NotNil(t, foundAuction)
	assert.Equal(t, auction_entity.Completed, foundAuction.Status)
}
