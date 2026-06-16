package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Katotodan/go-web/db"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type CreateTableRequest struct {
	TableName string `json:"tableName"`
}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Table    string `json:"table"`
}
type ListOfUser struct {
	Id        string
	UserName  string
	CreatedAt time.Time
}

func CreateTable(w http.ResponseWriter, r *http.Request) {

	var body CreateTableRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if body.TableName == "" {
		http.Error(w, "Table name is empty or invalid", http.StatusBadRequest)
		return
	}
	// Create table name inside the db
	query := fmt.Sprintf(`
		CREATE TABLE %s (
			id INT AUTO_INCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME,
			PRIMARY KEY (id)
		)
	`, body.TableName)
	_, err := db.Database.Exec(query)

	if err != nil {
		http.Error(w, "Failed to save table name", http.StatusInternalServerError)
		return
	}
	// Success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Table created successfully",
		"table":   body.TableName,
	})

}

func DropTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["tableName"]

	query := fmt.Sprintf(`
	    DROP TABLE %s
	`, table)
	_, err := db.Database.Exec(query)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to delete table %s", table)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Table deleted successfully",
		"table":   table,
	})
}

func GetAllTable(w http.ResponseWriter, r *http.Request) {
	query := `SHOW TABLES`
	rows, err := db.Database.Query(query)

	if err != nil {
		http.Error(w, "Failed to get the list of tables", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tables []string

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			http.Error(w, "Error reading tables", http.StatusInternalServerError)
			return
		}

		tables = append(tables, table)

	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tables": tables,
	})

}
func hashPasswordFunc(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func InsertUser(w http.ResponseWriter, r *http.Request) {

	var userData User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to get the user data", http.StatusBadRequest)
		return
	}

	hassPassword, err := hashPasswordFunc(userData.Password)

	if err != nil {
		http.Error(w, "Something went wrong, failed to hash the password", http.StatusInternalServerError)
		return
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (username, password, created_at) VALUES (?, ?, NOW())
	
	`, userData.Table)

	result, err := db.Database.Exec(query, userData.Username, hassPassword)
	if err != nil {
		http.Error(w, "Insertion failed", http.StatusInternalServerError)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get inserted ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "Success",
		"insertedId": userID,
	})

}

func GetUserData(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)
	userId := user["id"]
	table := user["table"]

	var (
		id        int
		username  string
		createdAt time.Time
	)

	query := fmt.Sprintf(`
	    SELECT id, username, created_at FROM %s WHERE id = ?
	`, table)

	err := db.Database.QueryRow(query, userId).Scan(&id, &username, &createdAt)

	if err != nil {
		http.Error(w, "Failed to get user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "Success",
		"id":        id,
		"username":  username,
		"createdAt": createdAt,
	})

}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var users []ListOfUser
	vars := mux.Vars(r)
	table := vars["table"]
	query := fmt.Sprintf(`
	    SELECT id, username, created_at FROM %s
	`, table)

	rows, err := db.Database.Query(query)
	if err != nil {
		http.Error(w, "Failed to get all users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user ListOfUser
		err := rows.Scan(&user.Id, &user.UserName, &user.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to iterate over list of users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)

	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to iterate over list of users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "Success",
		"users":  users,
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	table := vars["table"]

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)

	_, err := db.Database.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "Success",
		"id":      id,
		"message": "User deleted successfully",
	})

}
