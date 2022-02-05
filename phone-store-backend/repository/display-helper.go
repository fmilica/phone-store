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
	SELECT "id", "displayId", "mark"
	FROM "Rating"
	WHERE "displayId" = $1 and "parentId" is null
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

		err = rows.Scan(&id, &displayId, &mark)
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
	SELECT "id", "displayId", "content"
	FROM "Comment"
	WHERE "displayId" = $1 and "parentId" is null
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

		err = rows.Scan(&id, &displayId, &content)
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
