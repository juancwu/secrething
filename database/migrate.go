package database

func Migrate() {
	db.AutoMigrate(&User{})
}
