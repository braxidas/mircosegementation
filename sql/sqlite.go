package sql

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"errors"
)


var db *gorm.DB

func InitDB(){
	var err error
	if db, err = gorm.Open(sqlite.Open("/root/.cache/trivy/java-db/trivy-java.db"),&gorm.Config{}); err !=nil{
		log.Fatal("failed to connect database", err)
	}
	db.AutoMigrate(&Indices{})
}

//sha1值是否存在在java-db中 存在返回true 否则返回false
func QueryBySha1(sha1 []byte) (bool, error){
	indices := Indices{}
	if res := db.Where("sha1 = ?", sha1).First(&indices); res.Error != nil{
		if errors.Is(res.Error, gorm.ErrRecordNotFound){
			return false, res.Error
		}else{
			return false, nil
		}
	}
	return true, nil
}
