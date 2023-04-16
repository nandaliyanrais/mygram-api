package middlewares

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"mygram-api/models/domain"
	"mygram-api/models/response"
	"mygram-api/comments/service"
)

func Authorization(commentService service.CommentService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			comment domain.Comment
			err     error
		)

		commentID, _ := strconv.ParseUint(ctx.Param("commentId"), 10, 32)
		userData := ctx.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))

		if comment, err = commentService.GetOne(uint(commentID)); err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, response.ErrorResponse{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Errors: gin.H{
					"message": "Comment not found",
				},
			})

			return
		}

		if comment.UserID != userID {
			ctx.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse{
				Code:   http.StatusForbidden,
				Status: "Forbidden",
				Errors: gin.H{
					"message": "You don't have permission",
				},
			})

			return
		}

		ctx.Next()
	}
}