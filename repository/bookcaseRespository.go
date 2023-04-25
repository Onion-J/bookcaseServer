package repository

import (
	"BookcaseServer/common"
	"BookcaseServer/model"
	"BookcaseServer/validator"
	"gorm.io/gorm"
)

type BookcaseRepository struct {
	DB *gorm.DB
}

func NewBookcaseRepository() BookcaseRepository {
	return BookcaseRepository{
		DB: common.GetDB(),
	}
}

// Create 创建书柜（支持批量创建）
func (b BookcaseRepository) Create(bookcase []validator.CreateBookcaseRequest) ([]model.Cabinet, error) {
	cabinetList := make([]model.Cabinet, len(bookcase))
	for i := 0; i < len(bookcase); i++ {
		cabinetList[i] = model.Cabinet{
			Area:           bookcase[i].Area,
			SequenceNumber: bookcase[i].SequenceNumber,
			Occupied:       false,
		}

		row := b.DB.Where("area = ? and sequence_Number = ?", cabinetList[i].Area, cabinetList[i].SequenceNumber).First(&model.Cabinet{}).RowsAffected
		// 不存在则创建，存在则跳过创建
		if row == 0 {
			if err := b.DB.Create(&cabinetList[i]).Error; err != nil {
				return nil, err
			}
		}
	}
	return cabinetList, nil

}

func (b BookcaseRepository) SelectByAreaAndSequenceNumber(area string, sequenceNumber int) (*model.Cabinet, error) {
	var cabinet model.Cabinet
	if err := b.DB.Where("area = ? and sequence_Number = ?", area, sequenceNumber).First(&cabinet).Error; err != nil {
		return nil, err
	}
	return &cabinet, nil
}
