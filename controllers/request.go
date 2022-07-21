package controllers

import (
	"github.com/dParikesit/bnmo-backend/models"
	"github.com/dParikesit/bnmo-backend/utils"
)

func RequestGetBatch(page int64, pageSize int64, id uint) ([]models.Request, error) {
	var requests []models.Request

	result := utils.Db.Scopes(utils.Paginate(int(page), int(pageSize))).Where("id = ?", id).Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}

	return requests, nil
}
