package requests

type SearchContacts struct {
	FirstName string `json:"first_name" validator:"max:55"`
	LastName  string `json:"last_name" validator:"max:55"`
	Email     string `json:"email" validator:"email,max:255"`
}
