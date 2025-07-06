package main

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/infra/api/web/controller/auction_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/bid_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/user_controller"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/infra/database/bid"
	"fullcycle-auction_go/internal/infra/database/user"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"fullcycle-auction_go/internal/worker/auction_closer"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	auctionDurationMinutes, err := strconv.Atoi(os.Getenv("AUCTION_DURATION_MINUTES"))
	if err != nil {
		log.Fatal("Error trying to load AUCTION_DURATION_MINUTES env variable")
		return
	}

	auctionCloserIntervalSeconds, err := strconv.Atoi(os.Getenv("AUCTION_CLOSER_INTERVAL_SECONDS"))
	if err != nil {
		log.Fatal("Error trying to load AUCTION_CLOSER_INTERVAL_SECONDS env variable")
		return
	}

	client, databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	defer client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	auctionRepository := auction.NewAuctionRepository(databaseConnection)
	auctionCloserWorker := auction_closer.NewAuctionCloserWorker(
	auctionRepository, time.Duration(auctionCloserIntervalSeconds)*time.Second)
	go auctionCloserWorker.Start()

	router := gin.Default()

	userController, bidController, auctionsController := initDependencies(databaseConnection, auctionRepository, time.Duration(auctionDurationMinutes)*time.Minute)

	router.GET("/auction", auctionsController.FindAuctions)
	router.GET("/auction/:auctionId", auctionsController.FindAuctionById)
	router.POST("/auction", auctionsController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionsController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database, auctionRepository *auction.AuctionRepository, auctionDuration time.Duration) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController) {

	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))
	auctionController = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository, auctionDuration))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}
