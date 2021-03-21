package main

import (
	"fmt"
	"os"
	votingpb "github.com/guilherme-brandao/crypto-vote-platform/proto"
	"log"
	"io"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/guilherme-brandao/crypto-vote-platform/models"
	"google.golang.org/grpc"
)

var client votingpb.VotingServiceClient

var cryptos []models.Cryptocurrency

func main() {
	conn, err := grpc.Dial(os.Getenv("SERVER_URL"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client = votingpb.NewVotingServiceClient(conn)

	g := gin.Default()
	g.POST("/crypto", createCrypto)
	g.POST("/upvote/:id", upvoteCrypto)
	g.POST("/downvote/:id", downvoteCrypto)
	g.GET("/crypto/:id", getCrypto)
	g.GET("/cryptos", listCryptos)
	g.DELETE("/delete/:id", deleteCrypto)

	log.Fatal(g.Run(":8080"))

}

func createCrypto(ctx *gin.Context) {

	crypto := votingpb.Crypto{}
	err := ctx.ShouldBindJSON(&crypto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(&crypto)

	res, err := client.CreateCrypto(ctx, &crypto)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(res),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func upvoteCrypto(ctx *gin.Context) {
	uid := ctx.Param("id")

	obj := votingpb.UpvoteCryptoReq{Id: uid}

	res, err := client.UpvoteCrypto(ctx, &obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func downvoteCrypto(ctx *gin.Context) {
	uid := ctx.Param("id")

	obj := votingpb.DownvoteCryptoReq{Id: uid}

	res, err := client.DownvoteCrypto(ctx, &obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func getCrypto(ctx *gin.Context) {
	id := ctx.Param("id")

	obj := votingpb.GetCryptoReq{Id: id}

	res, err := client.GetCrypto(ctx, &obj)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func listCryptos(ctx *gin.Context) {

	obj := votingpb.ListCryptosReq{}

	stream, err := client.ListCryptos(ctx, &obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for {

		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(res.GetCrypto())
		ctx.JSON(http.StatusOK, res.GetCrypto())
	}

}

func deleteCrypto(ctx *gin.Context) {

	uid := ctx.Param("id")

	obj := votingpb.DeleteCryptoReq{Id: uid}

	res, err := client.DeleteCrypto(ctx, &obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
