package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liimadiego/schoolmatch/internal/database"
	"github.com/liimadiego/schoolmatch/internal/models"
)

type CreateSchoolRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

type UpdateSchoolRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type SchoolResponse struct {
	ID            uint            `json:"id"`
	Name          string          `json:"name"`
	Address       string          `json:"address"`
	Type          string          `json:"type"`
	UserID        uint            `json:"user_id"`
	AverageRating float32         `json:"average_rating"`
	Reviews       []models.Review `json:"reviews,omitempty"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

func CreateSchool(c *gin.Context) {
	var req CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	school := models.School{
		Name:    req.Name,
		Address: req.Address,
		Type:    req.Type,
		UserID:  userID.(uint),
	}

	if err := database.DB.Create(&school).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create school"})
		return
	}

	schoolResponse := SchoolResponse{
		ID:            school.ID,
		Name:          school.Name,
		Address:       school.Address,
		Type:          school.Type,
		UserID:        school.UserID,
		AverageRating: 0,
		CreatedAt:     school.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     school.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "school created successfully",
		"school":  schoolResponse,
	})
}

func GetSchools(c *gin.Context) {
	type SchoolWithRating struct {
		models.School
		AverageRating float32 `json:"average_rating"`
	}

	var schoolsWithRatings []SchoolWithRating

	query := `
		SELECT s.*, COALESCE(AVG(r.rating), 0) as average_rating
		FROM schools s
		LEFT JOIN reviews r ON s.id = r.school_id
		WHERE s.deleted_at IS NULL
		GROUP BY s.id
	`

	if err := database.DB.Raw(query).Scan(&schoolsWithRatings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch schools with ratings"})
		return
	}

	var schoolsResponse []SchoolResponse
	for _, s := range schoolsWithRatings {
		schoolResponse := SchoolResponse{
			ID:            s.ID,
			Name:          s.Name,
			Address:       s.Address,
			Type:          s.Type,
			UserID:        s.UserID,
			AverageRating: s.AverageRating,
			CreatedAt:     s.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     s.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		schoolsResponse = append(schoolsResponse, schoolResponse)
	}

	c.JSON(http.StatusOK, gin.H{"schools": schoolsResponse})
}

func GetSchool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school ID"})
		return
	}

	type SchoolWithRating struct {
		models.School
		AverageRating float32 `json:"average_rating"`
	}

	var schoolWithRating SchoolWithRating

	query := `
		SELECT s.*, COALESCE(AVG(r.rating), 0) as average_rating
		FROM schools s
		LEFT JOIN reviews r ON s.id = r.school_id
		WHERE s.id = ? AND s.deleted_at IS NULL
		GROUP BY s.id
	`

	if err := database.DB.Raw(query, id).Scan(&schoolWithRating).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}

	var reviews []models.Review
	if err := database.DB.Where("school_id = ?", id).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reviews"})
		return
	}

	schoolResponse := SchoolResponse{
		ID:            schoolWithRating.ID,
		Name:          schoolWithRating.Name,
		Address:       schoolWithRating.Address,
		Type:          schoolWithRating.Type,
		UserID:        schoolWithRating.UserID,
		AverageRating: schoolWithRating.AverageRating,
		Reviews:       reviews,
		CreatedAt:     schoolWithRating.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     schoolWithRating.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"school": schoolResponse})
}

func UpdateSchool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var school models.School
	if err := database.DB.First(&school, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}

	if school.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update schools you created"})
		return
	}

	if req.Name != "" {
		school.Name = req.Name
	}
	if req.Address != "" {
		school.Address = req.Address
	}
	if req.Type != "" {
		school.Type = req.Type
	}

	if err := database.DB.Save(&school).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update school"})
		return
	}

	schoolResponse := SchoolResponse{
		ID:            school.ID,
		Name:          school.Name,
		Address:       school.Address,
		Type:          school.Type,
		UserID:        school.UserID,
		AverageRating: 0,
		CreatedAt:     school.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     school.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "school updated successfully",
		"school":  schoolResponse,
	})
}

func DeleteSchool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var school models.School
	if err := database.DB.First(&school, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}

	if school.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete schools you created"})
		return
	}

	if err := database.DB.Where("school_id = ?", id).Delete(&models.Review{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete related reviews"})
		return
	}

	if err := database.DB.Delete(&school).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete school"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "school and related reviews deleted successfully"})
}
