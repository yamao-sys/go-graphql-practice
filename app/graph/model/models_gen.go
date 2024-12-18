// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateTodoInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Mutation struct {
}

type Query struct {
}

type SignInInput struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type SignUpInput struct {
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type UpdateTodoInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
