package models

// used to map the data from go to database

type Students struct {
	ID      int    `json:"id"` //json : "id" indicates that the ID field should be represented as "id" in json output
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Roll_no int    `json:"rollNo"`
	Age     int    `json:"age"`
	Class   string `json:"class"`
}
