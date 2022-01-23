package main
import(
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"errors"
	"strconv"
)

var db, _ = sql.Open("sqlite3","/var/lib/peerd/foo.db")

func CreateMessagesTable(){
	statement1, err := db.Prepare(`create table messages ("service" text, "userIdentifier" text, "idType" text, "message" blob)`);
	if err == nil{
		statement1.Exec()
	}
}

// add message to messages table, (and create table if it does not exist clean up this part plz)
func AddMessage(service string, userIdentifier string, idType string, message []byte) {
	fmt.Println(message)
	statement2, err := db.Prepare("insert into messages(service, userIdentifier, idType, message) values(?,?,?,?)")
	if err != nil{
		panic(err) // should this panic??
	}
	_, err = statement2.Exec(service, userIdentifier, idType, message)
	if err != nil{
		panic(err) // should this panic??
	}
}

// return message with given id
func GetMessage(id string) ([]byte, error) {
	statement, err := db.Prepare("select message from messages where rowid=?")
	if err != nil {
		panic(err)
	}
	var message []byte
	err = statement.QueryRow(id).Scan(&message)
	if err != nil{
		return nil, errors.New("no such message")
	}
	return message, nil
}

// return all message ids
func GetMessageIds() []string {
	row, err := db.Query("select rowid from messages")
	if err != nil {
		panic(err)
	}
	messageIds := make([]string, 0)
	defer row.Close()
	for row.Next(){
		var id int
		err = row.Scan(&id)
		messageIds = append(messageIds, strconv.Itoa(id))
	}
	return messageIds
}

func CloseDatabase(){
	db.Close()
}