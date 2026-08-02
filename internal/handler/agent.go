package handler

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/db"
	"server/internal/models"

	"github.com/gin-gonic/gin"
)

func agentGetHandler(c *gin.Context) {
	fmt.Println("agent with id get handler")
	rq, ctx := db.PrepareQueries()
	agents, err := rq.AgentRepo.GetAllAgents(ctx)
	if err != nil {
		log.Fatalf("Get Agents failed: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"agents": agents})
	fmt.Println("agents get successfully")
}

func agentRegisterHandler(c *gin.Context) {
	fmt.Println("agent register handler")
	agent, rq, ctx, err := generalPurpose[models.Agent](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.AgentRepo.Create(ctx, &agent)
	if err2 != nil {
		log.Fatalf("Create Agent failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"agent": agent})
	fmt.Println("agent register successfully")
}

func agentUpdateHandler(c *gin.Context) {
	fmt.Println("agent update handler")
	agent, rq, ctx, err := generalPurpose[models.Agent](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.AgentRepo.Update(ctx, &agent)
	if err2 != nil {
		log.Fatalf("Update Agent failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"agent": agent})
	fmt.Println("agent update successfully")
}

func agentDeleteHandler(c *gin.Context) {
	fmt.Println("agent delete handler")
	agentId, rq, ctx, err := generalPurpose[int64](c)
	if err != nil {
		fmt.Println(err)
	}
	err2 := rq.AgentRepo.Delete(ctx, agentId)
	if err2 != nil {
		log.Fatalf("Delete Agent failed: %v", err2)
	}
	c.JSON(http.StatusOK, gin.H{"agent": agentId})
	fmt.Println("agent deleted successfully")
}

func agentGetWithIdHandler(c *gin.Context) {
	fmt.Println("agent get with id handler")
	agentId, rq, ctx, err := generalPurpose[int64](c)
	if err != nil {
		fmt.Println(err)
	}

	agent, err := rq.AgentRepo.GetByID(ctx, agentId)
	if err != nil {
		log.Fatalf("Get Agent failed: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"agent": agent})
	fmt.Println("agent get successfully")
}
