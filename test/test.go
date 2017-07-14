package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ezbuy/ezorm/db"

	. "./model"
)

func init() {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://127.0.0.1"
	MgoSetup(conf)
}

func main() {
	for i := 0; i < 20000000000; i++ {
		p := BlogMgr.NewBlog()
		p.Title = "I like ezorm"
		p.Slug = fmt.Sprintf("ezorm_%d", time.Now().Nanosecond())

		_, err := p.Save()
		if err != nil {
			log.Fatal(err)
		}

		id := p.Id()

		_, err = BlogMgr.FindByID(id)
		if err != nil {
			log.Fatal(err.Error())

		}
	}

	fmt.Printf("get blog ok")

}
