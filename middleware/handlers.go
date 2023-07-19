package middleware

//middleware package serves as the bridge between APIs and the database, handling all crud operation
import (
	"encoding/json"
	"errors"
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
	ID      int64  `json:"id,omitempty"`
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

	var student models.Users
	// decoding a json request -> process of extracting the data sent in the body of an HTTP req

	err := json.NewDecoder(r.Body).Decode(&student) //Body bata leo student ko data

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertID := insertStudent(student)

	// format a response object
	res := Response{
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
		log.Fatalf("Unable to convert the string to int . %v", err)

	}

	//call getStudent func with user id ot retrieve a single user
	student, err := getStudent(int64(id))
	if err != nil {
		log.Fatalf("unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(student)
}

// "id" parameter is used when retrieving a single user in the 'GetUser' fn, to fetch a specific user based on the provided id
// fetch := action of retrieving or getting the desired data from database
func GetAllStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	students, err := getAllStudent()

	if err != nil {
		log.Fatalf("unable to get all the student. %v", err)
	}

	json.NewEncoder(w).Encode(students)

}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//get the user id from req params, key "id"

	var student models.Users
	err := json.NewDecoder(r.Body).Decode(&student) //request lai decode

	if err != nil {
		log.Fatalf("Unable to decode request body. %v", err)
	}
	fmt.Println(student.ID)
	updatedRows := updateStudent(student.ID, student)

	// message
	msg := fmt.Sprintf("Student updated successfully. Total rows affected %v", updatedRows)

	res := Response{
		ID:      student.ID,
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//get the userId from the req params "id"

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	fmt.Println(id)

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	deletedRows := deleteStudent(int64(id))
	msg := fmt.Sprintf("User updated successfully. Total rows affected: %v", deletedRows)

	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

// insert one user in db
func insertStudent(student models.Users) int64 {
	db := Database_connection()
	db.AutoMigrate(&models.Users{})
	result := db.Create(&student)
	if result.Error != nil {
		panic(fmt.Sprintf("Failed to execute the query: %v", result.Error))
	}
	fmt.Printf("Inserted a single record %v \n", student.ID)
	return student.ID
}

func getStudent(id int64) (models.Users, error) {
	db := Database_connection()

	var student models.Users
	//finding student by id
	result := db.First(&student, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("no rows were returned")
		return student, nil
	} else if result.Error != nil {
		log.Fatalf("unable to query the user. %v", result.Error)
		return student, result.Error
	}

	return student, nil

}

func getAllStudent() ([]models.Users, error) {
	db := Database_connection()
	var students []models.Users
	//retrieve all students from db
	result := db.Find(&students)
	if result.Error != nil {
		log.Fatalf("unable to find students. %v", result.Error)
	}
	return students, nil
}

func updateStudent(id int64, student models.Users) int64 {
	db := Database_connection()
	result := db.Model(&models.Users{}).Where("id = ?", id).Updates(student)
	if result.Error != nil {
		log.Fatalf("Unable to update student: %v", result.Error)
	}
	rowsAffected := result.RowsAffected
	log.Printf("total rows affected: %d", rowsAffected)
	return rowsAffected
}

func deleteStudent(id int64) int64 {
	db := Database_connection()
	result := db.Delete(&models.Users{}, id)
	if result.Error != nil {
		log.Fatalf("the record is not deleted")
	}
	rowsAffected := result.RowsAffected
	return rowsAffected
}
