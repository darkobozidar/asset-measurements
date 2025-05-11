package controllers

import (
	"asset/config"
	"asset/models"

	"net/http"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetAsset(c *gin.Context) {
	var asset models.Asset

	result := config.DB.First(&asset, "id = ? AND enabled = true", c.Param("id"))
	if err := result.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func GetAssets(c *gin.Context) {
	var assets []models.Asset

	enabledParam, errEnabledParam := strconv.ParseBool(c.DefaultQuery("enabled", "true"))
	typeParam := c.Query("type")

	if errEnabledParam != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'enabled' query param value."})
	}

	query := config.DB.Where("enabled = ?", enabledParam)

	if typeParam != "" {
		query = query.Where("type = ?", typeParam)
	}

	query = query.
	    Order("name").
		Find(&assets)

	if err := query.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO add pagination (or check how to add it).
	c.JSON(http.StatusOK, assets)
}

func CreateAsset(c *gin.Context) {
	// Why custom serializer?
	// - `Id`` is auto generated.
	// - `Enabled` is `true` for every newly created Asset.
	var body struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Type        string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset := models.Asset{
		Name: body.Name,
		Description: body.Description,
		Type: body.Type,
	}

	if err := config.DB.Create(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, asset)
}

func UpdateAsset(c *gin.Context) {
	var asset models.Asset
	var body struct {
		Name        *string `json:"name" binding:"omitempty,min=1"`
		Description *string `json:"description" binding:"omitempty,min=1"`
		Type        *string `json:"type" binding:"omitempty,min=1"`
	}
	assetId := c.Param("id")

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO check if these repeated lines can be somehow simplified.
	if body.Name != nil {
		asset.Name = *body.Name
	}
	if body.Description != nil {
		asset.Description = *body.Description
	}
	if body.Type != nil {
		asset.Type = *body.Type
	}

	// Creates a single UPDATE SQL statement. TODO double check.
	resultUpdate := config.DB.Model(&models.Asset{}).
		Where("id = ? AND enabled = true", assetId).
		Updates(&asset)

	if err := resultUpdate.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resultUpdate.RowsAffected == 0 {
		msg := fmt.Sprintf("Asset %s not found.", assetId)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
    }

	// TODO transaction between update and read.
	resultRead := config.DB.First(&asset, "id = ? AND enabled = true", assetId)
    if err := resultRead.Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
    }

	c.JSON(http.StatusOK, asset)
}

func DeleteAsset(c *gin.Context) {
	assetId := c.Param("id")
	result := config.DB.Model(&models.Asset{}).
		Where("id = ? AND enabled = true", assetId).
		Update("enabled", false)

	if err := result.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected == 0 {
		msg := fmt.Sprintf("Asset %s not found.", assetId)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	msg := fmt.Sprintf("Asset %s deleted successfully.", assetId)
	c.JSON(http.StatusOK, msg)
}
