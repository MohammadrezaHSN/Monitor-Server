package handler

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/db"
	"server/internal/models"

	"github.com/gin-gonic/gin"
)

func requestGetHandler(c *gin.Context) {
	fmt.Println("get request handler")
	rq, ctx := db.PrepareQueries()
	requests, err := rq.RequestRepo.GetAllRequests(ctx)
	if err != nil {
		log.Fatalf("Get All Requests failed: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"requests": requests})
	fmt.Println("get requests successful")
}

func requestCreateHandler(c *gin.Context) {
	fmt.Println("create request handler")
	request, rq, ctx, err := generalPurpose[models.Request](c)
	if err != nil {
		fmt.Println(err)
	}
	agent, err3 := rq.AgentRepo.GetByID(ctx, request.AgentID)
	if err3 != nil {
		fmt.Printf("Get Agent ByID failed: %v", err3)
	}
	err4 := rq.RequestRepo.Create(ctx, &request, agent.Interval)
	if err4 != nil {
		fmt.Printf("Create Request failed: %v", err4)
	}
	c.JSON(http.StatusOK, gin.H{"request": request})
}

func requestUpdateHandler(c *gin.Context) {
	fmt.Println("update request handler")
	request, rq, ctx, err := generalPurpose[models.Request](c)
	if err != nil {
		fmt.Println(err)
	}
	agent, err3 := rq.AgentRepo.GetByID(ctx, request.AgentID)
	if err3 != nil {
		fmt.Printf("Get Agent ByID failed: %v", err3)
	}
	err4 := rq.RequestRepo.Update(ctx, &request, agent.Interval)
	if err4 != nil {
		fmt.Printf("Update Request failed: %v", err4)
	}
	c.JSON(http.StatusOK, gin.H{"request": request})
}

func requestDeleteHandler(c *gin.Context) {
	fmt.Println("delete request handler")
	requestId, rq, ctx, err := generalPurpose[int64](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.RequestRepo.Delete(ctx, requestId)
	if err2 != nil {
		fmt.Printf("Delete Request failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"request deleted successfully": requestId})
}

func requestGetWithIdHandler(c *gin.Context) {
	fmt.Println("get request with id handler")
	requestId, rq, ctx, err := generalPurpose[int64](c)
	if err != nil {
		fmt.Println(err)
	}
	request, err2 := rq.RequestRepo.GetByID(ctx, requestId)
	if err2 != nil {
		fmt.Printf("Get Agent ByID failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"request": request})
}
