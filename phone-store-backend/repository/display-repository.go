package repository

import (
	"database/sql"
	"fmt"
	"phone-store-backend/model"
	"time"

	_ "github.com/lib/pq"
)

type DisplayRepository interface {
	Save(display *model.DisplayDTO) (*model.Display, error)
	Search(search *model.DisplaySearchDTO) ([]model.Display, error)
	FindAll() ([]model.Display, error)
}

type displayRepo struct{}

func NewDisplayRepository() DisplayRepository {
	return &displayRepo{}
}

func (*displayRepo) Save(displayDTO *model.DisplayDTO) (*model.Display, error) {

	fmt.Println("*** Adding phone display ***")

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database when return
	defer db.Close()

	display := getDisplayFromDTO(displayDTO)

	// insert phone to db
	date := display.Phone.Date.String()
	insertStmtPhone := `insert into "Phone"("id", "brand", "model", "date", "processor", "battery", "ram") values($1, $2, $3, $4, $5, $6, $7)`
	_, ePhone := db.Exec(insertStmtPhone, display.Phone.Id, display.Phone.Brand, display.Phone.Model,
		date[0:10], display.Phone.Processor, display.Phone.Battery, display.Phone.RAM)
	CheckError(ePhone)

	// insert display to db
	publishDate := display.Date.String()
	insertStmtDisplay := `insert into "Display"("id", "phoneId", "price", "date") values($1, $2, $3, $4)`
	_, eDisplay := db.Exec(insertStmtDisplay, display.Id, display.Phone.Id, display.Price, publishDate[0:10])
	CheckError(eDisplay)

	fmt.Println("Printamo display da vidimo sta smo napravili")
	fmt.Println(display)

	return display, nil
}

func (*displayRepo) FindAll() ([]model.Display, error) {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckErrorDisplay(err)

	// close database
	defer db.Close()

	query := `
		SELECT "D"."id", "D"."price", "D"."date", "D"."averagerate", "P"."id",
		"P"."brand", "P"."model", "P"."date", "P"."processor", "P"."battery", "P"."ram"
		FROM "Display" "D", "Phone" "P"
		WHERE "D"."phoneId" = "P"."id"`

	rows, err := db.Query(query)
	CheckErrorDisplay(err)

	defer rows.Close()

	displays := []model.Display{}

	// layout for parse string to date
	const layout = "2006-01-02"

	for rows.Next() {
		var id string
		var price int
		var date string
		var averageRate int
		var phoneId string
		var brand string
		var phoneModel string
		var phoneDate string
		var processor string
		var battery string
		var ram int

		err = rows.Scan(&id, &price, &date, &averageRate,
			&phoneId, &brand, &phoneModel, &phoneDate, &processor, &battery, &ram)
		CheckErrorDisplay(err)

		// Create phone
		var phone model.Phone
		phone.Id = phoneId
		phone.Brand = brand
		phone.Model = phoneModel
		phone.Processor = processor
		d1, _ := time.Parse(layout, phoneDate[0:10])
		phone.Date = d1
		phone.Battery = battery
		phone.RAM = ram

		// create rates
		ratings := getRatesByDisplay(id)

		// create comments
		comments := getCommentsByDisplay(id)

		d2, _ := time.Parse(layout, date[0:10])
		displays = append(displays, model.Display{Id: id, Price: price, Date: d2, AverageRating: averageRate,
			Phone: phone, Ratings: ratings, Comments: comments})
	}

	return displays, nil
}

func (*displayRepo) Search(search *model.DisplaySearchDTO) ([]model.Display, error) {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckErrorDisplay(err)

	// close database
	defer db.Close()

	searchPreprocess(search)

	query := getQuery(search)
	fmt.Println(query)

	//									$1				$2				$3				$4
	rows, err := db.Query(query, search.PriceFrom, search.PriceTo, search.DateFrom, search.DateTo)
	CheckErrorDisplay(err)

	defer rows.Close()

	displays := []model.Display{}

	// layout for parse string to date
	const layout = "2006-01-02"

	for rows.Next() {
		var id string
		var price int
		var date string
		var averageRate int
		var phoneId string
		var brand string
		var phoneModel string
		var phoneDate string
		var processor string
		var battery string
		var ram int

		err = rows.Scan(&id, &price, &date, &averageRate,
			&phoneId, &brand, &phoneModel, &phoneDate, &processor, &battery, &ram)
		CheckErrorDisplay(err)

		var phone model.Phone
		phone.Id = phoneId
		phone.Model = phoneModel
		phone.Brand = brand
		d1, _ := time.Parse(layout, phoneDate[0:10])
		phone.Date = d1
		phone.Processor = processor
		phone.Battery = battery
		phone.RAM = ram

		// create rates
		ratings := getRatesByDisplay(id)

		// create comments
		comments := getCommentsByDisplay(id)

		d2, _ := time.Parse(layout, date[0:10])
		displays = append(displays, model.Display{Id: id, Price: price, Date: d2, AverageRating: averageRate,
			Phone: phone, Ratings: ratings, Comments: comments})
	}

	return displays, nil
}
