// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreatePostInput struct {
	Content string `json:"content"`
}

type CreateSessionInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserInput struct {
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Email                string `json:"email"`
	PhoneNumber          string `json:"phoneNumber"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}