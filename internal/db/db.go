package db

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	// Получаем переменные окружения с значениями по умолчанию
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "db"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "zenrush"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "zenrush"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "zenrush"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	log.Printf("Подключаюсь к БД: host=%s, port=%s, db=%s, user=%s", dbHost, dbPort, dbName, dbUser)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := autoMigrate(); err != nil {
		return err
	}
	if err := seedInitialData(); err != nil {
		return err
	}
	return nil
}

func autoMigrate() error {
	log.Println("Начинаю создание таблиц...")

	// Создаём таблицы вручную через SQL
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(64) UNIQUE NOT NULL,
			password_hash VARCHAR(128) NOT NULL,
			role VARCHAR(16) DEFAULT 'user',
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS activities (
			id SERIAL PRIMARY KEY,
			name VARCHAR(128) NOT NULL,
			description TEXT,
			budget INT,
			time INT,
			weather VARCHAR(16),
			moods VARCHAR(64)[],
			created_at TIMESTAMP DEFAULT NOW(),
			deleted_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS favorites (
			user_id INT REFERENCES users(id),
			activity_id INT REFERENCES activities(id),
			PRIMARY KEY (user_id, activity_id)
		)`,
		`CREATE TABLE IF NOT EXISTS history (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			activity_id INT REFERENCES activities(id),
			viewed_at TIMESTAMP DEFAULT NOW()
		)`,
	}

	for i, query := range queries {
		if err := DB.Exec(query).Error; err != nil {
			log.Printf("Ошибка создания таблицы %d: %v", i+1, err)
			return err
		}
	}

	log.Println("Все таблицы успешно созданы")
	return nil
}

func seedInitialData() error {
	log.Println("Начинаю создание начальных данных...")

	// Проверяем подключение к БД
	var result int
	if err := DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		log.Printf("Ошибка проверки подключения к БД: %v", err)
		return err
	}
	log.Println("Подключение к БД работает")

	// Админ - используем простой запрос
	var adminExists bool
	if err := DB.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = 'admin')").Scan(&adminExists).Error; err != nil {
		log.Printf("Ошибка проверки существования админа: %v", err)
		return err
	}

	if !adminExists {
		log.Println("Создаю пользователя admin...")
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Ошибка хеширования пароля: %v", err)
			return err
		}

		insertQuery := `INSERT INTO users (username, password_hash, role, created_at) 
						VALUES ('admin', ?, 'admin', NOW())`
		if err := DB.Exec(insertQuery, string(hash)).Error; err != nil {
			log.Printf("Ошибка создания админа: %v", err)
			return err
		}
		log.Println("Админ создан успешно")
	} else {
		log.Println("Админ уже существует")
	}

	// Проверяем количество активностей
	var activityCount int
	if err := DB.Raw("SELECT COUNT(*) FROM activities").Scan(&activityCount).Error; err != nil {
		log.Printf("Ошибка подсчёта активностей: %v", err)
		return err
	}

	if activityCount == 0 {
		log.Println("Создаю примеры активностей...")

		// Используем кастомный SQL для правильной работы с массивами
		queries := []string{
			`INSERT INTO activities (name, description, budget, time, weather, moods, created_at)
			 VALUES ('Прогулка в парке', 'Приятная прогулка на свежем воздухе', 0, 2, 'sunny',
			 ARRAY['Нейтрально', 'Хорошо', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, moods, created_at)
			 VALUES ('Чтение книги', 'Уютно устроиться с интересной книгой', 0, 3, 'cloudy',
			 ARRAY['Спокойно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, moods, created_at)
			 VALUES ('Кофе с другом', 'Встретиться и поболтать за чашкой кофе', 300, 1, 'any',
			 ARRAY['Весело', 'Дружелюбно'], NOW())`,
		}

		for i, query := range queries {
			if err := DB.Exec(query).Error; err != nil {
				log.Printf("Ошибка создания активности %d: %v", i+1, err)
				return err
			}
		}
		log.Printf("Создано %d активностей", len(queries))
	} else {
		log.Printf("Активности уже существуют (%d штук)", activityCount)
	}

	log.Println("Начальные данные созданы успешно")
	return nil
}
