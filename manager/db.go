package manager

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDB(host, port, user, dbname, passwd string) *gorm.DB {
	s := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host, port, user, dbname, passwd)
	log.Print("db info %s", s)
	db, err := gorm.Open("postgres", s)

	if err != nil {
		log.Print("error init db", err)
		return nil
	}

	if !db.HasTable(&Plugin{}) {
		db.CreateTable(&Plugin{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Plugin{})
	} else {
		log.Print("plugins table already exist")
	}

	if !db.HasTable(&Strategy{}) {
		db.CreateTable(&Strategy{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Strategy{})
	} else {
		log.Print("strategies table already exist")
	}
	/*
		p := Plugin{Name: "pipeline", Kind: "nil", Status: "enable", Description: "nil", SpecJsonStr: "nil", Manual: "nil"}

		db.NewRecord(p)
		db.Create(&p)
	*/

	return db
	//defer db.Close()
}
