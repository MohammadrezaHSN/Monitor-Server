package handler

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/db"
	"server/internal/models"

	"github.com/gin-gonic/gin"
)

func fileGetHandler(c *gin.Context) {
	fmt.Println("get file handler")
	rq, ctx := db.PrepareQueries()
	files, err := rq.FileRepo.GetAllFiles(ctx)
	if err != nil {
		log.Fatalf("Get All Files failed: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
	fmt.Println("get files successful")
}

func fileStoreHandler(c *gin.Context) {
	fmt.Println("store files handler")
	file, rq, ctx, err := generalPurpose[models.File](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.FileRepo.Create(ctx, &file)
	if err2 != nil {
		log.Fatalf("Create File failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"file": file})
	fmt.Println("store file successful")
}

func fileUpdateHandler(c *gin.Context) {
	fmt.Println("file update handler")
	file, rq, ctx, err := generalPurpose[models.File](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.FileRepo.Update(ctx, &file)
	if err2 != nil {
		fmt.Printf("Update File failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"file": file})
	fmt.Println("file update successful")
}

func fileDeleteHandler(c *gin.Context) {
	fmt.Println("file delete handler")
	fileId, rq, ctx, err := generalPurpose[int64](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.FileRepo.Delete(ctx, fileId)
	if err2 != nil {
		fmt.Println(err2)
	}
	c.JSON(http.StatusOK, gin.H{"file": fileId})
	fmt.Println("file delete successful")
}

func fileGetWithIdHandler(c *gin.Context) {
	fmt.Println("get file with id handler")
	fileId, rq, ctx, err := generalPurpose[int64](c)
	if err != nil {
		fmt.Println(err)
	}
	file, err2 := rq.FileRepo.GetByID(ctx, fileId)
	if err2 != nil {
		log.Fatalf("Get File failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"file": file})
	fmt.Println("get with id file successful")
}
