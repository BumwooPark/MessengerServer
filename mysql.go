package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type database struct {
	db *sql.DB
}

var db database

func init() {

	zip, err := sql.Open("mysql","root:1q2w3e4r@tcp(210.100.238.118:3306)/ZIP")
	if err != nil {
		log.Fatalln(err)
	}

	db = database{db: zip}
}

func (database *database) searchIDAndRoom(user string, room string) (userID int, roomID int ,err error) {

	row := database.db.QueryRow("SELECT chat_room.id, chat_member.id FROM chat_room JOIN chat_member ON chat_room.id = chat_member.ROOM_id AND chat_member.id=? AND chat_room.id=?",user,room)

	err = row.Scan(&userID,&roomID)
	if err != nil {
		return
	}
	return
}