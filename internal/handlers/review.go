package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liimadiego/schoolmatch/internal/database"
	"github.com/liimadiego/schoolmatch/internal/models"
)

type CreateReviewRequest struct {
	Rating   float32 `json:"rating" binding:"required,min=1,max=5"`
	Comment  string  `json:"comment"`
	SchoolID uint    `json:"school_id" binding:"required"`
}

type UpdateReviewRequest struct {
	Rating  float32 `json:"rating" binding:"min=1,max=5"`
	Comment string  `json:"comment"`
}

func CreateReview(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var school models.School
	if err := database.DB.First(&school, req.SchoolID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}

	var existingReview models.Review
	result := database.DB.Where("user_id = ? AND school_id = ?", userID, req.SchoolID).First(&existingReview)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you already have a review for this school, please update it instead"})
		return
	}

	review := models.Review{
		Rating:   req.Rating,
		Comment:  req.Comment,
		UserID:   userID.(uint),
		SchoolID: req.SchoolID,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "review created successfully",
		"review":  review,
	})
}

func GetReviews(c *gin.Context) {
	schoolID, err := strconv.ParseUint(c.Param("school_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school ID"})
		return
	}

	var reviews []models.Review
	if err := database.DB.Where("school_id = ?", schoolID).Preload("User").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func GetReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	var review models.Review
	if err := database.DB.Preload("User").First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": review})
}

func UpdateReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var review models.Review
	if err := database.DB.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	if review.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own reviews"})
		return
	}

	if req.Rating >= 1 && req.Rating <= 5 {
		review.Rating = req.Rating
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}

	if err := database.DB.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "review updated successfully",
		"review":  review,
	})
}

func DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var review models.Review
	if err := database.DB.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	if review.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own reviews"})
		return
	}

	if err := database.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}
