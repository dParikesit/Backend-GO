package controllers

import (
	"github.com/dParikesit/bnmo-backend/models"
	"github.com/dParikesit/bnmo-backend/utils"
)

func TransactionGetBatch(page int64, pageSize int64, id uint) ([]models.Transaction, error) {
	var transactions []models.Transaction

	result := utils.Db.Scopes(utils.Paginate(int(page), int(pageSize))).Where("id = ?", id).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}
