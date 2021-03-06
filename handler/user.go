package handler

import (
	"net/http"
	"time"

	"github.com/c8112002/twitter_clone_go/entities"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Users(c echo.Context) error {
	time.Sleep(time.Second * 1)

	users, err := h.userStore.FetchUsers(maxID(c), minID(c), limit(c))

	if err != nil {
		c.Logger().Error("db error: " + err.Error())
		return err
	}

	firstUser, err := h.userStore.FetchFirstUser()

	if err != nil {
		c.Logger().Error("db error: " + err.Error())
		return err
	}

	res := newEmptyUsersResponse()
	for _, u := range *users {
		ur := newUserResponse(u, u.IsFollowedBy(entities.LoginUserID))
		res.Users = append(res.Users, ur)
	}

	res.ContainsFirstUser = containsFirstUser(firstUser, users)

	return c.JSON(http.StatusOK, res)
}

// usersにfirstUserが含まれている場合trueを返す
func containsFirstUser(firstUser *entities.User, users *entities.Users) bool {
	for _, u := range *users {
		if u.ID == firstUser.ID {
			return true
		}
	}

	return false
}
