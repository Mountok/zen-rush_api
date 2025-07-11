# ZenRush Backend

Бэкенд для ZenRush на Go + PostgreSQL.

## Быстрый старт (Docker Compose)

1. Клонируйте репозиторий и перейдите в папку backend:
   ```sh
   cd backend
   ```
2. Запустите сервисы:
   ```sh
   docker-compose up --build
   ```
3. Бэкенд будет доступен по адресу: http://localhost:8080
   PostgreSQL — на порту 5432 (user/pass/db: zenrush)

**Переменные окружения для backend:**
- `DB_HOST` — адрес базы (по умолчанию: db)
- `DB_PORT` — порт базы (5432)
- `DB_USER` — пользователь базы (zenrush)
- `DB_PASSWORD` — пароль базы (zenrush)
- `DB_NAME` — имя базы (zenrush)
- `JWT_SECRET` — секрет для подписи JWT (замените на свой в проде)

> ⚡️ Миграции выполняются автоматически при запуске backend — ничего руками делать не нужно.

---

# Документация API

## Аутентификация

### Регистрация
`POST /api/auth/register`

**Что передать:**
```
{
  "username": "имя_пользователя",
  "password": "пароль"
}
```
**Ответ:**
- 201 Created — всё ок
- 400 — пользователь уже есть или ошибка валидации

### Вход (логин)
`POST /api/auth/login`

**Что передать:**
```
{
  "username": "имя_пользователя",
  "password": "пароль"
}
```
**Ответ:**
```
{
  "token": "ваш_JWT_токен"
}
```

**Важно:**
Для всех защищённых эндпоинтов нужен заголовок:
```
Authorization: Bearer <ваш_JWT_токен>
```

---

## Рекомендации (Activities)

### Получить список активностей
`GET /api/activities`

**Можно фильтровать по параметрам:**
- `min_budget` — минимальный бюджет
- `max_budget` — максимальный бюджет
- `time` — время (часы)
- `mood` — настроение (например: Весело)g
- `weather` — погода (sunny/cloudy/rainy)
- `people_count` — количество людей (1, 2, 3, 4, 5+)

**Пример:**
```
curl -H "Authorization: Bearer <JWT>" "http://localhost:8080/api/activities?min_budget=0&max_budget=1000&mood=Весело&weather=sunny"
```

### Получить одну активность
`GET /api/activities/:id`

### Создать активность (только moderator/admin)
`POST /api/activities`

**Что передать:**
```
{
  "name": "Название",
  "description": "Описание",
  "budget": 0,
  "time": 2,
  "weather": "sunny",
  "people_count": 2,
  "moods": ["Нейтрально", "Хорошо"]
}
```

### Обновить активность (только moderator/admin)
`PUT /api/activities/:id`

### Удалить активность (только moderator/admin)
`DELETE /api/activities/:id`

---

## Избранное (Favorites)

### Получить избранное
`GET /api/favorites`

### Добавить в избранное
`POST /api/favorites/:activity_id`

### Удалить из избранного
`DELETE /api/favorites/:activity_id`

---

## История просмотров (History)

### Получить последние 10 просмотренных
`GET /api/history`

### Добавить просмотр
`POST /api/history/:activity_id`

---

## Как обычно работает процесс
1. Регистрируешься: `/api/auth/register`
2. Логинишься: `/api/auth/login` (получаешь JWT)
3. Используешь JWT для всех остальных запросов
4. Получаешь рекомендации, добавляешь в избранное, смотришь историю — всё просто!

---

## Примеры CURL

**Регистрация:**
```
curl -X POST http://localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username": "test", "password": "123456"}'
```

**Логин:**
```
curl -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username": "test", "password": "123456"}'
```

**Получить активности:**
```
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/api/activities
``` 