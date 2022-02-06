package repository

import (
	"database/sql"
	"fmt"
	"phone-store-backend/model"
	"time"

	"github.com/google/uuid"
)

func getDisplayFromDTO(displayDTO *model.DisplayDTO) *model.Display {

	// Create phone
	var phone model.Phone
	phone.Id = uuid.New().String()
	phone.Brand = displayDTO.Brand
	phone.Model = displayDTO.Model
	phone.Date = displayDTO.Date
	phone.Processor = displayDTO.Processor
	phone.Battery = displayDTO.Battery
	phone.RAM = displayDTO.RAM

	// Empty comment and rate lists
	comments := []model.Comment{}
	ratings := []model.Rating{}

	var display model.Display
	display.Id = uuid.New().String()
	display.Phone = phone
	display.Price = displayDTO.Price
	display.Date = time.Now()
	display.Ratings = ratings
	display.Comments = comments

	return &display
}

func getRatesByDisplay(id string) []model.Rating {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckErrorDisplay(err)

	// close database
	defer db.Close()

	query := `
	SELECT "id", "displayId", "parentId", "mark"
	FROM "Rating"
	WHERE "displayId" = $1 and "parentId" = ''
	`
	rows, err := db.Query(query, id)
	CheckErrorDisplay(err)

	defer rows.Close()

	var ratings = []model.Rating{}

	for rows.Next() {
		var id string
		var displayId string
		var parentId string
		var mark int

		err = rows.Scan(&id, &displayId, &parentId, &mark)
		CheckErrorDisplay(err)

		ratings = append(ratings, model.Rating{Id: id, DisplayId: displayId, ParentId: parentId,
			Mark: mark})
	}

	return ratings
}

/*
	Find all comments by display id
*/
func getCommentsByDisplay(id string) []model.Comment {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckErrorDisplay(err)

	// close database
	defer db.Close()

	query := `
	SELECT "id", "displayId", "parentId", "content"
	FROM "Comment"
	WHERE "displayId" = $1 and "parentId" = ''
	`
	rows, err := db.Query(query, id)
	CheckErrorDisplay(err)

	defer rows.Close()

	var comments = []model.Comment{}

	for rows.Next() {
		var id string
		var displayId string
		var parentId string
		var content string

		err = rows.Scan(&id, &displayId, &parentId, &content)
		CheckErrorDisplay(err)

		commentComments := getCommentsByComment(id)
		ratings := getRatingsByComment(id)

		comments = append(comments, model.Comment{Id: id, DisplayId: displayId, ParentId: parentId,
			Content: content, Comments: commentComments, Ratings: ratings})
	}

	return comments
}

func getCommentsByComment(id string) []model.Comment {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckErrorDisplay(err)

	// close database
	defer db.Close()

	query := `
	SELECT "id", "displayId", "parentId", "content"
	FROM "Comment"
	WHERE "parentId" = $1
	`
	rows, err := db.Query(query, id)
	CheckErrorDisplay(err)

	defer rows.Close()

	var comments = []model.Comment{}

	for rows.Next() {
		var id string
		var displayId string
		var parentId string
		var content string

		err = rows.Scan(&id, &displayId, &parentId, &content)
		CheckErrorDisplay(err)

		commentComments := getCommentsByComment(id)
		ratings := getRatingsByComment(id)

		comments = append(comments, model.Comment{Id: id, DisplayId: displayId, ParentId: parentId,
			Content: content, Comments: commentComments, Ratings: ratings})
	}

	return comments
}

func getRatingsByComment(id string) []model.Rating {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckErrorDisplay(err)

	// close database
	defer db.Close()

	query := `
	SELECT "id", "displayId", "parentId", "mark"
	FROM "Rating"
	WHERE "parentId" = $1
	`
	rows, err := db.Query(query, id)
	CheckErrorDisplay(err)

	defer rows.Close()

	var ratings = []model.Rating{}

	for rows.Next() {
		var id string
		var displayId string
		var parentId string
		var mark int

		err = rows.Scan(&id, &displayId, &parentId, &mark)
		CheckErrorDisplay(err)

		ratings = append(ratings, model.Rating{Id: id, DisplayId: displayId, ParentId: parentId,
			Mark: mark})
	}

	return ratings
}

func CheckErrorDisplay(err error) {
	if err != nil {
		panic(err)
	}
}

func searchPreprocess(search *model.DisplaySearchDTO) {

	if search.PriceTo <= 0 {
		search.PriceTo = 99999999
	}

	// if search.RAM <= 0 {
	// 	search.RAM = 9999
	// }
}

func getQuery(search *model.DisplaySearchDTO) string {

	query := `
	SELECT "D"."id", "D"."price", "D"."date", "D"."averagerate", "P"."id",
	"P"."brand", "P"."model", "P"."date", "P"."processor", "P"."battery", "P"."ram"
	FROM "Display" "D", "Phone" "P"
	WHERE "D"."phoneId" = "P"."id"
	AND "D"."price" BETWEEN $1 AND $2
	AND "P"."date" BETWEEN $3 AND $4
	` + queryFilterBrand(search.Brand) + queryFilterProcessor(search.Processor) +
		queryFilterBattery(search.Battery) + queryFilterRAM(search.RAM) +
		querySortPrice(search.Sort) + querySortDate(search.Sort) + querySortAverageRate(search.Sort)

	return query
}

func queryFilterBrand(brand string) string {

	if brand != "" {
		return `AND "P"."brand" ilike '%` + brand + `%'`
	}

	return ``
}

func queryFilterProcessor(processor string) string {

	if processor != "" {
		return `AND "P"."processor" ilike '%` + processor + `%'`
	}

	return ``
}

func queryFilterBattery(battery string) string {

	if battery != "" {
		return `AND "P"."battery" ilike '%` + battery + `%'`
	}

	return ``
}

func queryFilterRAM(ram int) string {

	if ram != 0 {
		return fmt.Sprintf("%s%d", `AND "P"."ram" = `, ram)
	}

	return ``
}

func querySortPrice(sort string) string {

	if sort == "price up" {
		return `ORDER BY "D"."price" ASC`
	} else if sort == "price down" {
		return `ORDER BY "D"."price" DESC`
	}

	return ``
}

func querySortDate(sort string) string {

	if sort == "oldest" {
		return `ORDER BY "P"."date" ASC`
	} else if sort == "latest" {
		return `ORDER BY "P"."date" DESC`
	}

	return ``
}

func querySortAverageRate(sort string) string {

	if sort == "average rate" {
		return `ORDER BY "D"."averagerate" DESC`
	}
	return ``
}

/*
	Remove offers whose publish date
	is not in dateFrom-dateTo range.
*/
// func filterByDate(offers []model.Display, dateFrom time.Time, dateTo time.Time) []model.Display {

// 	for idx, offer := range offers {
// 		if !(offer.Vehicle.Date.After(dateFrom) && offer.Vehicle.Date.Before(dateTo)) {
// 			return append(offers[0:idx], offers[idx+1:]...)
// 		}
// 	}

// 	return offers
// }

/*
	Sort list by publish date.
	If asc is true sort offer by newest,
	otherwise by oldest.
*/
// func sortByDate(offers []model.Offer, sortStr string) []model.Offer {

// 	if sortStr == "newest" { // sort asc
// 		sort.Slice(offers, func(i, j int) bool {
// 			return offers[i].Date.After(offers[j].Date)
// 		})
// 	} else if sortStr == "oldest" { //sort desc
// 		sort.Slice(offers, func(i, j int) bool {
// 			return offers[i].Date.Before(offers[j].Date)
// 		})
// 	}

// 	return offers

// }
