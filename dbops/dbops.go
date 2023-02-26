package dbops

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/vivekkb14/goblogbackend/common"
)

type DataBase struct {
	Db *sql.DB
}

var GlobalDatabase *DataBase

func (database *DataBase) InitialiseDatabaeServer() error {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "new_db_kbv",
	}
	// Get a database handle.
	var err error
	database.Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	pingErr := database.Db.Ping()
	if pingErr != nil {
		fmt.Println(pingErr.Error())
		return pingErr
	}
	return nil
}

func (database *DataBase) CreateProductTable() error {
	query := `CREATE TABLE IF NOT EXISTS articles(Id int primary key auto_increment, Title text, 
		Content text, Author text)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := database.Db.ExecContext(ctx, query)
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when getting rows affected", err)
		return err
	}
	fmt.Printf("Rows affected when creating table: %d\n", rows)
	return nil
}

func (database *DataBase) InserIntoTable(userInfo common.UserRequest) (int, error) {
	query := `INSERT INTO articles (Title, Content, Author) VALUES
  ('` + userInfo.Title + `', '` + userInfo.Content + `', '` + userInfo.Author + `');`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := database.Db.ExecContext(ctx, query)
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when getting rows affected", err)
		return 0, err
	}
	fmt.Printf("Rows affected when creating table: %d\n", rows)
	id := database.ExecuteQuery("")
	fmt.Println("Id is: ", id)
	return id, nil
}

func (database *DataBase) DispalyData() ([]common.Author, error) {
	var author common.Author
	var allauthors []common.Author
	query := `select * from articles;`
	res, err := database.Db.Query(query)
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return allauthors, err
	}
	for res.Next() {

		err := res.Scan(&author.Id, &author.Title, &author.Content, &author.Author)

		if err != nil {
			fmt.Println(err)
			return allauthors, err
		}

		fmt.Printf("Author info %v\n", author)
		allauthors = append(allauthors, author)
	}

	return allauthors, nil
}

func (database *DataBase) SelectAuthorInfo(id int) (common.Author, error) {
	fmt.Println("Inside selectAuthorInfo 1")
	var author common.Author
	query := `select * from articles where id=` + fmt.Sprintf("%d", id) + `;`
	res, err := database.Db.Query(query)
	respData := res.Next()
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return author, err
	}
	if respData {
		err := res.Scan(&author.Id, &author.Title, &author.Content, &author.Author)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", author)
		if author.Id == 0 {
			return author, errors.New("Couldn't not find blog of the id number")
		}
	} else {
		return author, errors.New("Couldn't not find blog of the id number")
	}
	return author, nil
}

func (database *DataBase) ExecuteQuery(queryInfo string) int {
	var id int
	query := `SELECT Id FROM articles ORDER BY id DESC LIMIT 1;`
	res, err := database.Db.Query(query)
	if err != nil {
		fmt.Printf("Error %s when creating product table", err)
		return 0
	}
	for res.Next() {
		err := res.Scan(&id)
		if err != nil {
			fmt.Println("Error getting id number\n", err)
			return 0
		}
	}
	return id
}
