package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (ar *AuctionRepository) FindAuctionById(
	ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find auction by id = %s", id), err)
		return nil, internal_error.NewInternalServerError("Error trying to find auction by id")
	}

	return &auction_entity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
		ExpiresAt:   time.Unix(auctionEntityMongo.ExpiresAt, 0),
	}, nil
}

func (repo *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category string,
	productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["productName"] = primitive.Regex{Pattern: productName, Options: "i"}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding auctions", err)
		return nil, internal_error.NewInternalServerError("Error finding auctions")
	}
	defer cursor.Close(ctx)

	var auctionsMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding auctions", err)
		return nil, internal_error.NewInternalServerError("Error decoding auctions")
	}

	var auctionsEntity []auction_entity.Auction
	for _, auction := range auctionsMongo {
		auctionsEntity = append(auctionsEntity, auction_entity.Auction{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Status:      auction.Status,
			Description: auction.Description,
			Condition:   auction.Condition,
			Timestamp:   time.Unix(auction.Timestamp, 0),
			ExpiresAt:   time.Unix(auction.ExpiresAt, 0),
		})
	}

	return auctionsEntity, nil
}

func (ar *AuctionRepository) FindExpiredAuctions(
	ctx context.Context) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{
		"status":     auction_entity.Active,
		"expires_at": bson.M{"$lte": time.Now().Unix()},
	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error trying to find expired auctions", err)
		return nil, internal_error.NewInternalServerError("Error trying to find expired auctions")
	}
	defer cursor.Close(ctx)

	var auctionsMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding expired auctions", err)
		return nil, internal_error.NewInternalServerError("Error decoding expired auctions")
	}

	var auctionsEntity []auction_entity.Auction
	for _, auction := range auctionsMongo {
		auctionsEntity = append(auctionsEntity, auction_entity.Auction{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   auction.Condition,
			Status:      auction.Status,
			Timestamp:   time.Unix(auction.Timestamp, 0),
			ExpiresAt:   time.Unix(auction.ExpiresAt, 0),
		})
	}

	return auctionsEntity, nil
}