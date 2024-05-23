package database

import (
	"fmt"
	"log"
	"os"
	"vcs_server/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SMSDatabse interface {
	CreateServer(server *entity.Server) error
	ViewServer(id int) (entity.Server, error)
	ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error)
	UpdateServer(server *entity.Server) error
	DeleteServer(id int) error
	CheckServerExistence(ip string) bool
	CheckServerName(name string) bool
	GetServerByIp(ip string) entity.Server
	FindUserByUsername(username string) entity.User
	AddUser(user *entity.User) error
}

type smsDatabase struct {
	connection *gorm.DB
}

var DB SMSDatabse

func init() {

	// Get environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Data source name
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh", host, user, password, dbname, port)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Println("Database connected")
	db.AutoMigrate(&entity.Server{})
	db.AutoMigrate(&entity.User{})
	DB = &smsDatabase{connection: db}
}

func (db *smsDatabase) CreateServer(server *entity.Server) error {
	result := db.connection.Create(server)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *smsDatabase) ViewServer(id int) (entity.Server, error) {
	var server entity.Server

	if err := db.connection.Where("id = ?", id).First(&server).Error; err != nil {
		return server, err
	}
	fmt.Println("Execute ViewServer() function from sms-database.go file")
	return server, nil
}

func (db *smsDatabase) ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error) {
	var servers []entity.Server

	sql := "SELECT * FROM servers"

	if filter != "" {
		sql = fmt.Sprintf("%s WHERE status='%s'", sql, filter)
	}

	if !(from == 0 && to == 0 && perpage == 0) {
		if from == 0 && to != 0 {
			from = to
		}
		if from != 0 && to == 0 {
			to = from
		}
		sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perpage*(to-from+1), (from-1)*perpage)
	}

	if order == "" {
		order = "asc"
	}
	if sortby != "" {
		sql = fmt.Sprintf("%s ORDER BY %s %s", sql, sortby, order)
	}

	if err := db.connection.Raw(sql).Scan(&servers).Error; err != nil {
		return nil, err
	}

	fmt.Println("Execute ViewServers() function from sms-database.go file")
	return servers, nil
}

func (db *smsDatabase) UpdateServer(server *entity.Server) error {
	id := server.Id
	err := db.connection.Where("id = ?", id).Updates(server).Error
	if err != nil {
		return err
	}
	fmt.Println("Execute UpdateServer() function from sms-database.go file")
	return nil
}
func (db *smsDatabase) DeleteServer(id int) error {
	err := db.connection.Table(entity.Server{}.TableName()).Where("id = ?", id).Delete(nil).Error
	if err != nil {
		return err
	}
	fmt.Println("Execute DeleteServer() function from sms-database.go file")
	return nil
}

func (db *smsDatabase) CheckServerExistence(ip string) bool {
	var server entity.Server
	if err := db.connection.Where("ipv4 = ?", ip).First(&server).Error; err != nil {
		return false
	}
	return true
}

func (db *smsDatabase) CheckServerName(name string) bool {
	var server entity.Server
	if err := db.connection.Where("name = ?", name).First(&server).Error; err != nil {
		return false
	}
	return true
}

func (db *smsDatabase) GetServerByIp(ip string) entity.Server {
	var server entity.Server
	db.connection.Where("ipv4 = ?", ip).First(&server)
	return server
}

func (db *smsDatabase) FindUserByUsername(username string) entity.User {
	var user entity.User
	db.connection.Where("username = ?", username).First(&user)
	return user
}

func (db *smsDatabase) AddUser(user *entity.User) error {
	result := db.connection.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
