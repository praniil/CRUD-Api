package middleware

//middleware package serves as the bridge between APIs and the database, handling all crud operation
import (
	"encoding/json"
	"fmt"
	"go-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Response := the data or information that is returned from server when an API request is sent
type Response struct {
	ID      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

var Database *gorm.DB

func Database_connection() *gorm.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error in loading .env file")
	}

	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database_name := os.Getenv("DB_NAME")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		log.Fatal()
	}

	psql_info := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", host, port, username, database_name, password)

	Database, err := gorm.Open(postgres.Open(psql_info), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}

	return Database
}

// Api Endpoint Handlers
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var student models.Students
	// decoding a json request -> process of extracting the data sent in the body of an HTTP req

	err := json.NewDecoder(r.Body).Decode(&student) //Body bata leo student ko data

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertID := insertStudent(student)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "User created Successfully",
	}

	json.NewEncoder(w).Encode(res) //writes the response by encoding

}

//get user

func GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//get the student id from request params, key is "id"
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"]) //convert the id type from string to int
	if err != nil {
		log.Fatalf("Unable to converty the string to int . %v", err)

	}

	//call getStudent func with user id ot retrieve a single user
	student, err := getStudent(int64(id))
	if err != nil {
		log.Fatalf("unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(student)
}

//"id" parameter is used when retrieving a single user in the 'GetUser' fn, to fetch a specific user based on the provided id
//fetch := action of retrieving or getting the desired data from database
func GetAllStudent(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	students, err := getAllStudent()

	if err!= nil{
		log.Fatalf("unable to get all the student. %v", err)
	}

	json.NewEncoder(w).Encode(students)

}
