package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Katotodan/go-web/db"
	"github.com/gorilla/mux"
)

type CreateTableRequest struct {
	TableName string `json:"tableName"`
}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Table    string `json:"table"`
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

func InsertUser(w http.ResponseWriter, r *http.Request) {

	var userData User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to get the user data", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (username, password, created_at) VALUES (?, ?, NOW())
	
	`, userData.Table)

	result, err := db.Database.Exec(query, userData.Username, userData.Password)
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
