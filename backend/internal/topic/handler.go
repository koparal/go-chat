package topic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TopicHandler struct {
	Service *TopicService
}

func NewTopicHandler(service *TopicService) *TopicHandler {
	return &TopicHandler{
		Service: service,
	}
}

// CreateTopic godoc
// @Summary Create a new topic
// @Description Create a new topic with the provided data
// @Tags topics
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header"
// @Param input body Topic true "Topic data to create"
// @Success 200 {object} map[string]string "OK"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /topics [post]
func (h *TopicHandler) CreateTopic(c *gin.Context) {
	var topic Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.CreateTopic(c.Request.Context(), &topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Topic created successfully"})
}

// GetTopics godoc
// @Summary Get all topics
// @Description Get a list of all topics
// @Tags topics
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header"
// @Success 200 {object} []Topic "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /topics [get]
func (h *TopicHandler) GetTopics(c *gin.Context) {
	topics, err := h.Service.GetTopics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topics)
}

// UpdateTopic godoc
// @Summary Update a topic
// @Description Update a topic with the provided data
// @Tags topics
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header"
// @Param id path string true "Topic ID" Format(uuid)
// @Param input body Topic true "Updated topic data"
// @Success 200 {object} map[string]string "OK"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /topics/{id} [post]
func (h *TopicHandler) UpdateTopic(c *gin.Context) {
	id := c.Param("id")

	var topic Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.UpdateTopic(c.Request.Context(), id, &topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Topic updated successfully"})
}

// DeleteTopic godoc
// @Summary Delete a topic
// @Description Delete a topic with the provided ID
// @Tags topics
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header"
// @Param id path string true "Topic ID" Format(uuid)
// @Success 200 {object} map[string]string "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /topics/delete/{id} [post]
func (h *TopicHandler) DeleteTopic(c *gin.Context) {
	id := c.Param("id")

	err := h.Service.DeleteTopic(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Topic %s deleted successfully", id)})
}
