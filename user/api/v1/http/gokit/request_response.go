package gokit

import "gophr.v2/user"

type getByUserIDRequest struct {
	UserID string `json:"userID"`
}

type getByUserIDResponse struct {
	User  *user.User `json:"user,omitempty"`
	Error string     `json:"error,omitempty"`
}

type getAllRequest struct {
	Cursor string `json:"cursor"`
	Number int    `json:"num"`
}

type getAllResponse struct {
	Users      []*user.User `json:"users,omitempty"`
	NextCursor string       `json:"next_cursor,omitemptY"`
	Error      string       `json:"error,omitempty"`
}

type registerRequest struct {
	*user.User
}

type registerResponse struct {
	Error string `json:"error,omitempty"`
}

type updateRequest struct {
	*user.User
}

type updateResponse struct {
	Error string `json:"error,omitempty"`
}

type deleteRequest struct {
	ID string `json:"id"`
}

type deleteResponse struct {
	Error string `json:"error,omitempty"`
}
