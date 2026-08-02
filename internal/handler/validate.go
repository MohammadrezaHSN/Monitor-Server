package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"server/internal/db"
	"server/internal/models"

	"github.com/gin-gonic/gin"
)

type ValidModels interface {
	models.Request | models.Response | models.Agent | models.File | int64
}

func generalPurpose[T ValidModels](c *gin.Context) (T, db.RepoQueries, context.Context, error) {
	var inputType T
	rq, ctx := db.PrepareQueries()
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return inputType, rq, nil, err
	}
	err = json.Unmarshal(rawData, &inputType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return inputType, rq, nil, err
	}
	return inputType, rq, ctx, nil

}
