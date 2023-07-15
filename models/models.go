package models

// used to map the data from go to database

type Users struct {
	ID      int64  `json:"id"` //json : "id" indicates that the ID field should be represented as "id" in json output
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Roll_no int64  `json:"rollNo"`
	Age     int64  `json:"age"`
	Class   string `json:"class"`
}
