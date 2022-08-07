package db

import (
	"github.com/sukenda/golang-krakend/product-service/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgre struct {
	GormBD *gorm.DB
}

func Init(url string) Postgre {
	open, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	open.AutoMigrate(&model.Product{})
	open.AutoMigrate(&model.StockDecreaseLog{})

	return Postgre{open}
}
