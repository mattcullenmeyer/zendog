package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattcullenmeyer/zendog/pkg/models"
)

type AddRepEventPayload struct {
	StartUtc int      `json:"start_utc" binding:"required"`
	EndUtc   int      `json:"end_utc" binding:"required"`
	Day      string   `json:"day" binding:"required"`
	Start    string   `json:"start" binding:"required"`
	Duration int      `json:"duration" binding:"required"`
	Goal     int      `json:"goal" binding:"required"`
	Success  *bool    `json:"success" binding:"required"`
	Comment  *string  `json:"comment" binding:"required"`
	Behavior []string `json:"behavior" binding:"required"`
	User     string   `json:"user" binding:"required"`
}

func AddRepEvent(c *gin.Context) {
	var payload AddRepEventPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createRepEventArgs := models.CreateRepEventParams{
		StartUtc: payload.StartUtc,
		Day:      payload.Day,
		Start:    payload.Start,
		Duration: payload.Duration,
		Goal:     payload.Goal,
		Success:  *payload.Success,
		Comment:  *payload.Comment,
		Behavior: payload.Behavior,
		User:     payload.User,
	}

	if err := models.CreateRepEvent(createRepEventArgs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updateLastRepArgs := models.UpdateLastRepParams{
		EndUtc: payload.EndUtc,
	}

	if err := models.UpdateLastRep(updateLastRepArgs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fetchDaysEventsArgs := models.FetchTodaysEventsParams{
		Day: payload.Day,
	}

	response, err := models.FetchDaysEvents(fetchDaysEventsArgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repCount := 0
	successCount := 0
	durationSum := 0
	for _, event := range response {
		if event.Type == "rep" {
			repCount++
			durationSum += event.Duration

			if event.Success {
				successCount++
			}
		}
	}

	averageDuration := int(float64(durationSum) / float64(repCount))

	updateDaysStatsArgs := models.UpdateDaysStatsParams{
		Day:             payload.Day,
		RepCount:        repCount,
		SuccessCount:    successCount,
		AverageDuration: averageDuration,
	}

	if err := models.UpdateDaysStats(updateDaysStatsArgs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}
