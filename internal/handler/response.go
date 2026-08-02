package handler

import (
	"fmt"
	"net/http"
	"server/internal/db"
	"server/internal/models"

	"github.com/gin-gonic/gin"
)

func responseGetHandler(c *gin.Context) {
	fmt.Println("response get handler")
	rq, ctx := db.PrepareQueries()
	responses, err := rq.ResponseRepo.GetAllResponses(ctx)
	if err != nil {
		fmt.Println("get response handler")
	}
	c.JSON(http.StatusOK, gin.H{"responses": responses})
	fmt.Println("get responses successful")
}

func responseCreateHandler(c *gin.Context) {
	fmt.Println("create response handler")
	response, rq, ctx, err := generalPurpose[models.Response](c)
	if err != nil {
		fmt.Println(err)
	}

	err2 := rq.ResponseRepo.Create(ctx, &response)
	if err2 != nil {
		fmt.Printf("Create response failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func responseUpdateHandler(c *gin.Context) {
	fmt.Println("update response handler")
	response, rq, ctx, err := generalPurpose[models.Response](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.ResponseRepo.Update(ctx, &response)
	if err2 != nil {
		fmt.Printf("Update response failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func responseDeleteHandler(c *gin.Context) {
	fmt.Println("delete response handler")
	responseId, rq, ctx, err := generalPurpose[int64](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.ResponseRepo.Delete(ctx, responseId)
	if err2 != nil {
		fmt.Printf("Delete response failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"response": responseId})
}

func responseGetWithIdHandler(c *gin.Context) {
	fmt.Println("get response with id handler")
	responseId, rq, ctx, err := generalPurpose[int64](c)
	if err != nil {
		fmt.Println(err)
	}
	response, err2 := rq.ResponseRepo.GetByID(ctx, responseId)
	if err2 != nil {
		fmt.Printf("Get response failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}
