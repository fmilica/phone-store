package main

import (
	"database/sql"
	"fmt"
	"phone-store-backend/model"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "ntp"
)

/*
	Method to create all tables.
*/
func createTablesDB() {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database when return
	defer db.Close()

	// insert to db
	insertStmt := `
	CREATE TABLE "Phone"(
		"id" TEXT NOT NULL,
		"brand" TEXT NOT NULL,
		"model" TEXT NOT NULL,
		"date" TEXT NOT NULL,
		"processor" TEXT NOT NULL,
		"battery" TEXT NOT NULL,
		"ram" INTEGER NOT NULL
	);
	ALTER TABLE
		"Phone" ADD PRIMARY KEY("id");
	CREATE TABLE "Display"(
		"id" TEXT NOT NULL,
		"phoneId" TEXT NOT NULL,
		"price" INTEGER NOT NULL,
		"date" TEXT NOT NULL
	);
	ALTER TABLE
		"Display" ADD PRIMARY KEY("id");
	CREATE TABLE "Comment"(
		"id" TEXT NOT NULL,
		"displayId" TEXT NOT NULL,
		"parentId" TEXT,
		"content" TEXT NOT NULL
	);
	ALTER TABLE
		"Comment" ADD PRIMARY KEY("id");
	CREATE TABLE "Rating"(
		"id" TEXT NOT NULL,
		"displayId" TEXT NOT NULL,
		"parentId" TEXT,
		"mark" INTEGER NOT NULL
	);
	ALTER TABLE
		"Rating" ADD PRIMARY KEY("id");
	ALTER TABLE
		"Display" ADD CONSTRAINT "display_phoneid_foreign" FOREIGN KEY("phoneId") REFERENCES "Phone"("id");
	ALTER TABLE
		"Comment" ADD CONSTRAINT "comment_displayid_foreign" FOREIGN KEY("displayId") REFERENCES "Display"("id");
	ALTER TABLE
		"Comment" ADD CONSTRAINT "comment_parentid_foreign" FOREIGN KEY("parentId") REFERENCES "Comment"("id");
	ALTER TABLE
		"Rating" ADD CONSTRAINT "rating_displayid_foreign" FOREIGN KEY("displayId") REFERENCES "Display"("id");
	ALTER TABLE
		"Rating" ADD CONSTRAINT "rating_parentid_foreign" FOREIGN KEY("parentId") REFERENCES "Comment"("id");
	`
	_, e := db.Exec(insertStmt)
	CheckError(e)

	fmt.Println("---------- DB tables are created ----------")
}

/*
	Create one entity for all tables.
*/
func createAllInit() {

	var display model.Display
	display = getDisplay()

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database when return
	defer db.Close()

	var insertStmt string

	// insert phone into db
	date := display.Phone.Date.String()
	insertStmt = `insert into "Phone"("id", "brand", "model", "date", "processor", "battery", "ram") values($1, $2, $3, $4, $5, $6, $7)`
	_, e1 := db.Exec(insertStmt, display.Phone.Id, display.Phone.Brand, display.Phone.Model,
		date[0:10], display.Phone.Processor, display.Phone.Battery, display.Phone.RAM)
	CheckError(e1)

	// insert phone display into db
	publishDate := display.Date.String()
	insertStmt = `insert into "Display"("id", "phoneId", "price", "date") values($1, $2, $3, $4)`
	_, e2 := db.Exec(insertStmt, display.Id, display.Phone.Id, display.Price,
		publishDate[0:10])
	CheckError(e2)

	// insert rate into db
	insertStmt = `insert into "Rating"("id", "displayId", "mark") values($1, $2, $3)`
	_, e3 := db.Exec(insertStmt, display.Ratings[0].Id, display.Id, display.Ratings[0].Mark)
	CheckError(e3)

	// insert comment into db
	insertStmt = `insert into "Comment"("id", "displayId", "content") values($1, $2, $3)`
	_, e4 := db.Exec(insertStmt, display.Comments[0].Id, display.Id, display.Comments[0].Content)
	CheckError(e4)

	commentOfDisplay := display.Comments[0]

	// insert comment of comment into db
	insertStmt = `insert into "Comment"("id", "displayId", "parentId", "content") values($1, $2, $3, $4)`
	_, e6 := db.Exec(insertStmt, commentOfDisplay.Comments[0].Id, display.Id, commentOfDisplay.Comments[0].ParentId, commentOfDisplay.Comments[0].Content)
	CheckError(e6)

	// insert comment rate into db
	insertStmt = `insert into "Rating"("id", "displayId", "parentId", "mark") values($1, $2, $3, $4)`
	_, e5 := db.Exec(insertStmt, commentOfDisplay.Ratings[0].Id, display.Id, commentOfDisplay.Ratings[0].ParentId, commentOfDisplay.Ratings[0].Mark)
	CheckError(e5)

	fmt.Println("---------- Entities added to db ----------")
}

/*
	Delete all entities from all tables
*/
func deleteAll() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database when return
	defer db.Close()

	var insertStmt string

	// delete all from rating
	insertStmt = `DELETE FROM "Rating"`
	_, e1 := db.Exec(insertStmt)
	CheckError(e1)

	// delete all from comment
	insertStmt = `DELETE FROM "Comment"`
	_, e2 := db.Exec(insertStmt)
	CheckError(e2)

	// delete all from display
	insertStmt = `DELETE FROM "Display"`
	_, e3 := db.Exec(insertStmt)
	CheckError(e3)

	// delete all from phone
	insertStmt = `DELETE FROM "Phone"`
	_, e4 := db.Exec(insertStmt)
	CheckError(e4)

	// insertStmt = `DROP TABLE "Rating"`
	// _, e5 := db.Exec(insertStmt)
	// CheckError(e5)

	// insertStmt = `DROP TABLE "Comment"`
	// _, e6 := db.Exec(insertStmt)
	// CheckError(e6)

	// insertStmt = `DROP TABLE "Display"`
	// _, e7 := db.Exec(insertStmt)
	// CheckError(e7)

	// insertStmt = `DROP TABLE "Phone"`
	// _, e8 := db.Exec(insertStmt)
	// CheckError(e8)

}

/*
	manually create and return display
*/
func getDisplay() model.Display {
	const layout = "2006-01-02"
	d, _ := time.Parse(layout, "2018-05-05")
	dPhone, _ := time.Parse(layout, "2018-04-05")

	var display model.Display
	display.Id = uuid.New().String()
	display.Date = d
	display.Price = 260

	var phone model.Phone

	phone.Id = uuid.New().String()
	phone.Brand = "Xiaomi"
	phone.Model = "Mi 9 Lite"
	phone.Processor = "Snapdragon 832"
	phone.Battery = "Lithium 5000"
	phone.RAM = 8

	phone.Date = dPhone

	// add phone to display
	display.Phone = phone

	var comment1 model.Comment
	comment1.Id = uuid.New().String()
	comment1.DisplayId = display.Id
	comment1.Content = "this phone is amazing"

	var comment2 model.Comment
	comment2.Id = uuid.New().String()
	comment2.DisplayId = display.Id
	comment2.ParentId = comment1.Id
	comment2.Content = "this comment is right"

	//add comment to comment
	comment1.Comments = append(comment1.Comments, comment2)

	var rate1 model.Rating
	rate1.Id = uuid.New().String()
	rate1.DisplayId = display.Id
	rate1.Mark = 4

	// add rate to display
	display.Ratings = append(display.Ratings, rate1)

	var rate2 model.Rating
	rate2.Id = uuid.New().String()
	rate2.DisplayId = display.Id
	rate2.ParentId = comment1.Id
	rate2.Mark = 4

	// add rate to comment
	comment1.Ratings = append(comment1.Ratings, rate2)

	// add comment to display
	display.Comments = append(display.Comments, comment1)

	return display
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
