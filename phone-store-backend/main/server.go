package main

import (
	"fmt"
	"net/http"
	"phone-store-backend/controller"
	"phone-store-backend/repository"
	"phone-store-backend/service"

	"phone-store-backend/router"
)

var (
	phoneRepository repository.PhoneRepository = repository.NewPhoneRepository()
	phoneService    service.PhoneService       = service.NewPhoneService(phoneRepository)
	phoneController controller.PhoneController = controller.NewPhoneController(phoneService)

	displayRepository repository.DisplayRepository = repository.NewDisplayRepository()
	displayService    service.DisplayService       = service.NewDisplayService(displayRepository)
	displayController controller.DisplayController = controller.NewDisplayController(displayService)

	commentRepository repository.CommentRepository = repository.NewCommentRepository()
	commentService    service.CommentService       = service.NewCommentService(commentRepository)
	commentController controller.CommentController = controller.NewCommentController(commentService)

	ratingRepository repository.RatingRepository = repository.NewRatingRepository()
	ratingService    service.RatingService       = service.NewRatingService(ratingRepository)
	ratingController controller.RatingController = controller.NewRatingController(ratingService)

	httpRouter router.Router = router.NewMuxRouter()
)

func runServer() {
	const port string = ":8000"

	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Up and runing")
	})

	// Phones
	httpRouter.GET("/phones", phoneController.GetAll)
	httpRouter.POST("/phones", phoneController.Save)
	httpRouter.DELETE("/phones/delete-all", phoneController.DeleteAll)

	// Displays
	httpRouter.POST("/displays/search", displayController.Search)
	httpRouter.GET("/displays/getAll", displayController.GetAll)
	httpRouter.POST("/displays/add-displays", displayController.Save)

	// Comments
	httpRouter.POST("/comments", commentController.Save)

	// Ratings
	httpRouter.POST("/ratings", ratingController.Save)

	httpRouter.SERVE(port)
}
