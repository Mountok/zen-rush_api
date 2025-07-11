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
			people_count INT DEFAULT 1,
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
			// Бесплатные активности
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Прогулка в парке', 'Приятная прогулка на свежем воздухе', 0, 2, 'sunny', 1, 
			 ARRAY['Нейтрально', 'Хорошо', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Чтение книги', 'Уютно устроиться с интересной книгой', 0, 3, 'cloudy', 1, 
			 ARRAY['Спокойно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Медитация', 'Расслабляющая медитация для души', 0, 1, 'any', 1, 
			 ARRAY['Спокойно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Йога дома', 'Утренняя практика для бодрости', 0, 1, 'any', 1, 
			 ARRAY['Спокойно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Рисование', 'Творческий процесс с красками', 0, 2, 'any', 1, 
			 ARRAY['Вдохновенно', 'Спокойно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Прослушивание музыки', 'Любимые треки для настроения', 0, 1, 'any', 1, 
			 ARRAY['Весело', 'Спокойно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Фотографирование', 'Съёмка интересных моментов', 0, 2, 'sunny', 1, 
			 ARRAY['Вдохновенно', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Вечерняя прогулка', 'Романтичная прогулка под звёздами', 0, 1, 'any', 2, 
			 ARRAY['Романтично', 'Спокойно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Пикник на природе', 'Отдых на свежем воздухе', 0, 4, 'sunny', 4, 
			 ARRAY['Весело', 'Дружелюбно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Написание дневника', 'Запись мыслей и планов', 0, 1, 'any', 1, 
			 ARRAY['Спокойно', 'Вдохновенно'], NOW())`,

			// Недорогие активности (до 500₽)
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Кофе с другом', 'Встретиться и поболтать за чашкой кофе', 300, 1, 'any', 2, 
			 ARRAY['Весело', 'Дружелюбно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Посещение музея', 'Культурное просвещение', 400, 3, 'any', 2, 
			 ARRAY['Вдохновенно', 'Интересно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Кино в кинотеатре', 'Новый фильм на большом экране', 500, 3, 'any', 2, 
			 ARRAY['Весело', 'Интересно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Боулинг', 'Активная игра с друзьями', 400, 2, 'any', 4, 
			 ARRAY['Весело', 'Активно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Лазертаг', 'Захватывающая командная игра', 450, 2, 'any', 6, 
			 ARRAY['Активно', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Квест-комната', 'Интеллектуальное развлечение', 500, 2, 'any', 4, 
			 ARRAY['Интересно', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Мастер-класс по рисованию', 'Творческое развитие', 400, 2, 'any', 8, 
			 ARRAY['Вдохновенно', 'Интересно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Скалодром', 'Активный спорт для всех', 350, 2, 'any', 2, 
			 ARRAY['Активно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Бильярд', 'Классическая игра для компании', 300, 2, 'any', 4, 
			 ARRAY['Весело', 'Дружелюбно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Настольные игры', 'Интеллектуальное развлечение', 200, 3, 'any', 4, 
			 ARRAY['Весело', 'Интересно'], NOW())`,

			// Средние активности (500-1500₽)
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Ресторан', 'Ужин в хорошем ресторане', 1200, 2, 'any', 2, 
			 ARRAY['Романтично', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('СПА-салон', 'Расслабляющие процедуры', 1500, 3, 'any', 1, 
			 ARRAY['Спокойно', 'Романтично'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Концерт', 'Живая музыка и эмоции', 1000, 4, 'any', 4, 
			 ARRAY['Весело', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Театр', 'Классическое искусство', 800, 4, 'any', 2, 
			 ARRAY['Вдохновенно', 'Интересно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Картинг', 'Скорость и адреналин', 800, 2, 'any', 2, 
			 ARRAY['Активно', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Пейнтбол', 'Командная игра на природе', 600, 3, 'sunny', 8, 
			 ARRAY['Активно', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Верёвочный парк', 'Активный отдых на высоте', 700, 3, 'sunny', 4, 
			 ARRAY['Активно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Массаж', 'Расслабляющий массаж', 1000, 2, 'any', 1, 
			 ARRAY['Спокойно', 'Романтично'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Кулинарный мастер-класс', 'Обучение готовке', 800, 3, 'any', 6, 
			 ARRAY['Интересно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Экскурсия по городу', 'Познавательная прогулка', 600, 4, 'sunny', 8, 
			 ARRAY['Интересно', 'Вдохновенно'], NOW())`,

			// Дорогие активности (1500₽+)
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Прыжок с парашютом', 'Экстремальные эмоции', 5000, 4, 'sunny', 1, 
			 ARRAY['Активно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Полёт на воздушном шаре', 'Романтичное приключение', 8000, 3, 'sunny', 2, 
			 ARRAY['Романтично', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Дайвинг', 'Исследование подводного мира', 3000, 5, 'sunny', 2, 
			 ARRAY['Активно', 'Интересно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Сёрфинг', 'Покорение волн', 2500, 4, 'sunny', 1, 
			 ARRAY['Активно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Горные лыжи', 'Зимний спорт', 4000, 6, 'cloudy', 2, 
			 ARRAY['Активно', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Сноуборд', 'Экстремальный зимний спорт', 3500, 5, 'cloudy', 1, 
			 ARRAY['Активно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Вертолётная экскурсия', 'Вид на город с высоты', 6000, 2, 'sunny', 4, 
			 ARRAY['Вдохновенно', 'Романтично'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Баня с друзьями', 'Традиционный отдых', 2000, 4, 'any', 6, 
			 ARRAY['Весело', 'Дружелюбно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Рыбалка', 'Спокойный отдых на природе', 1500, 6, 'sunny', 2, 
			 ARRAY['Спокойно', 'Интересно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Охота', 'Активный отдых в лесу', 3000, 8, 'sunny', 4, 
			 ARRAY['Активно', 'Интересно'], NOW())`,

			// Домашние активности
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Готовка нового блюда', 'Кулинарные эксперименты', 500, 2, 'any', 2, 
			 ARRAY['Интересно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Просмотр сериала', 'Уютный вечер дома', 0, 3, 'any', 2, 
			 ARRAY['Спокойно', 'Весело'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Уборка и организация', 'Приведение дома в порядок', 0, 2, 'any', 1, 
			 ARRAY['Спокойно', 'Вдохновенно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Игра на музыкальном инструменте', 'Творческое самовыражение', 0, 1, 'any', 1, 
			 ARRAY['Вдохновенно', 'Спокойно'], NOW())`,
			`INSERT INTO activities (name, description, budget, time, weather, people_count, moods, created_at) 
			 VALUES ('Вязание или рукоделие', 'Создание чего-то своими руками', 200, 2, 'any', 1, 
			 ARRAY['Спокойно', 'Вдохновенно'], NOW())`,
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
