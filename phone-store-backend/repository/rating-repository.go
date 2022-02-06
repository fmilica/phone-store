package repository

import (
	"database/sql"
	"fmt"
	"math"
	"phone-store-backend/model"

	_ "github.com/lib/pq"
)

type RatingRepository interface {
	Save(rating *model.Rating) (*model.Rating, error)
	FindAll() ([]model.Rating, error)
}

type ratingRepo struct{}

func NewRatingRepository() RatingRepository {
	return &ratingRepo{}
}

func (*ratingRepo) Save(rating *model.Rating) (*model.Rating, error) {

	fmt.Println("*** Adding rating ***")

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database when return
	defer db.Close()

	// insert to db
	insertStmt := `insert into "Rating"("id", "displayId", "parentId", "mark") values($1, $2, $3, $4)`
	_, e := db.Exec(insertStmt, rating.Id, rating.DisplayId, rating.ParentId, rating.Mark)
	CheckError(e)

	//update display? average rate
	ratings := getRatesByDisplay(rating.DisplayId)

	var sum float64
	sum = 0.0
	for _, rating := range ratings {
		sum += float64(rating.Mark)
	}

	newAverage := int(math.Ceil(sum / float64(len(ratings))))

	// update
	insertStmtDisplay := `update "Display" set "averagerate" = $1 where "id" = $2`
	_, e1 := db.Exec(insertStmtDisplay, newAverage, rating.DisplayId)
	CheckError(e1)

	return rating, nil
}

func (*ratingRepo) FindAll() ([]model.Rating, error) {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckErrorRating(err)

	// close database
	defer db.Close()

	rows, err := db.Query(`SELECT "id", "displayId", "parentId", "mark" FROM "Rating"`)
	CheckErrorRating(err)

	defer rows.Close()

	var ratings []model.Rating

	for rows.Next() {
		var id string
		var displayId string
		var parentId string
		var mark int

		err = rows.Scan(&id, &displayId, &parentId, &mark)
		CheckErrorRating(err)

		ratings = append(ratings, model.Rating{Id: id, DisplayId: displayId, ParentId: parentId,
			Mark: mark})
	}

	return ratings, nil
}

func CheckErrorRating(err error) {
	if err != nil {
		panic(err)
	}
}
