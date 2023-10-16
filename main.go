package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("无法加载配置:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}

	store := db.NewStore(conn)
	//runGinServer(config, store)
	runGinServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("服务器创建失败:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
