# RegistrationCodingTest

For coding test purpose

## Requirement
1. PostgreSQL Database
2. Go Language
3. Go Package
    - github.com/gorilla/mux
    - github.com/jinzhu/gorm

## Set Up Steps
1. Download or clone the project
2. Create the database
3. Change database connection configuration in config/db.go
4. Run the project (go run main.go) at least once to create necessary tables

## Note
1. Registration page's link is localhost:8080/register
2. main_test.go is for unit test purpose
3. db.go conts that need to set:
    - DB_USER is the dbms user (default is "postgres")
    - DB_PASSWORD is the dbms user password
    - DB_NAME the database's name