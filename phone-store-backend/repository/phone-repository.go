package repository

import (
	"database/sql"
	"fmt"
	"phone-store-backend/model"
	"time"

	_ "github.com/lib/pq"
)

type PhoneRepository interface {
	Save(phone *model.Phone) (*model.Phone, error)
	FindAll() ([]model.Phone, error)
	DeleteAll()
}

type phoneRepo struct{}

func NewPhoneRepository() PhoneRepository {
	return &phoneRepo{}
}

func (*phoneRepo) Save(phone *model.Phone) (*model.Phone, error) {

	fmt.Println("*** Adding phone ***")

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database when return
	defer db.Close()

	// insert to db
	insertStmt := `insert into "Phone"("id", "brand", "model", "date", "proccessor", "battery", "ram") values($1, $2, $3, $4, $5, $6, $7)`
	_, e := db.Exec(insertStmt, phone.Id, phone.Brand, phone.Model, phone.Date, phone.Processor,
		phone.Battery, phone.RAM)
	CheckError(e)

	return phone, nil
}

func (*phoneRepo) FindAll() ([]model.Phone, error) {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	rows, err := db.Query(`SELECT "id", "brand", "model", "date", "proccessor", "battery", "ram" FROM "Phone"`)
	CheckError(err)

	defer rows.Close()

	var phones []model.Phone

	for rows.Next() {
		var id string
		var brand string
		var phoneModel string
		var date string
		var processor string
		var battery string
		var ram int

		err = rows.Scan(&id, &brand, &phoneModel, &date, &processor, &battery, &ram)
		CheckError(err)

		const layout = "2006-01-02"
		d, _ := time.Parse(layout, date[0:10])
		phones = append(phones, model.Phone{Id: id, Brand: brand, Model: phoneModel,
			Date: d, Processor: processor, Battery: battery, RAM: ram})
	}

	return phones, nil
}

func (*phoneRepo) DeleteAll() {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database when return
	defer db.Close()

	// insert to db
	insertStmt := `DELETE FROM "Phone"`
	_, e := db.Exec(insertStmt)
	CheckError(e)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
