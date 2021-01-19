package pkg

import (
	"context"
	"database/sql"
	"log"

	// sqlite
	"github.com/georgysavva/scany/sqlscan"
	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ItemID    int64  `db:"itemId"      json:"itemId"`
	ListID    string `db:"listId"      json:"listId"`
	Text      string `db:"contentText" json:"text"`
	IsChecked bool   `db:"isChecked"   json:"isChecked"`
}

type Store interface {
	Bootstrap()
	GetItem(itemId int64) *Item
	GetItems(listId string) []*Item
	AddItem(item *Item) *Item
	UpdateItem(item *Item) *Item
}

type SqliteStore struct {
	DB *sql.DB
}

// // insert
// insert, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
// checkErr(err)

// res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
// checkErr(err)

// id, err := res.LastInsertId()
// checkErr(err)

func NewStore(dbfile string) Store {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		panic(err)
	}
	return &SqliteStore{
		DB: db,
	}
}

func (store *SqliteStore) Bootstrap() {
	stmt, err := store.DB.Prepare(`
		CREATE TABLE IF NOT EXISTS ListItems (
			itemId      INTEGER      PRIMARY KEY AUTOINCREMENT,
			listId      VARCHAR(128) NOT NULL,
			isChecked   TINYINT 	 DEFAULT 0,
			contentText TEXT
		) ;
		CREATE INDEX idx_list_item ON ListItems (listId, itemId);
		`)
	defer stmt.Close()
	if _, err = stmt.Exec(); err != nil {
		log.Println("created a new database")
	}
}

func (store *SqliteStore) GetItem(itemId int64) *Item {
	var items []*Item
	ctx := context.Background()
	err := sqlscan.Select(ctx, store.DB, &items, `SELECT itemId, listId, isChecked, contentText FROM ListItems WHERE itemId = ?`, itemId)
	if err != nil || len(items) != 1 {
		log.Println("item not found:", itemId)
		return nil
	}
	return items[0]
}

func (store *SqliteStore) GetItems(listId string) []*Item {
	var items []*Item
	ctx := context.Background()
	err := sqlscan.Select(ctx, store.DB, &items, `SELECT itemId, listId, isChecked, contentText FROM ListItems WHERE listId = ?`, listId)
	if err != nil {
		log.Println(err)
	}
	return items
}

func (store *SqliteStore) AddItem(item *Item) *Item {
	stmt, err := store.DB.Prepare(`
		INSERT INTO 
			ListItems (listId, isChecked, contentText)
			   VALUES (?, ?, ?)`)
	defer stmt.Close()
	var res sql.Result
	if res, err = stmt.Exec(item.ListID, item.IsChecked, item.Text); err != nil {
		log.Println(err)
	}
	id, err := res.LastInsertId()
	return store.GetItem(id)
}

func (store *SqliteStore) UpdateItem(item *Item) *Item {
	stmt, err := store.DB.Prepare(`
		UPDATE ListItems SET isChecked = ? WHERE itemId = ?`)
	defer stmt.Close()
	if _, err = stmt.Exec(item.IsChecked, item.ItemID); err != nil {
		log.Println(err)
	}
	return store.GetItem(item.ItemID)
}
