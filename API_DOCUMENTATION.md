# ZenRush API Documentation

## Base URL
```
http://localhost:8080/api
```

## Аутентификация
Все защищённые эндпоинты требуют JWT токен в заголовке:
```
Authorization: Bearer <JWT_TOKEN>
```

---

## 1. Аутентификация

### Регистрация пользователя
**POST** `/auth/register`

**Тело запроса:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Ответы:**
- `201 Created` - пользователь успешно создан
- `400 Bad Request` - пользователь уже существует или ошибка валидации

**Пример:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "123456"}'
```

### Вход в систему
**POST** `/auth/login`

**Тело запроса:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Ответ:**
```json
{
  "token": "JWT_TOKEN_STRING"
}
```

**Пример:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

---

## 2. Рекомендации (Activities)

### Получить список всех активностей
**GET** `/activities`

**Query параметры:**
- `min_budget` (int) - минимальный бюджет
- `max_budget` (int) - максимальный бюджет
- `time` (int) - время в часах
- `mood` (string) - настроение (например: "Весело")
- `weather` (string) - погода ("sunny", "cloudy", "rainy", "any")

**Ответ:**
```json
[
  {
    "id": 1,
    "name": "Прогулка в парке",
    "description": "Приятная прогулка на свежем воздухе",
    "budget": 0,
    "time": 2,
    "weather": "sunny",
    "moods": ["Нейтрально", "Хорошо", "Весело"],
    "created_at": "2025-07-10T21:00:00Z"
  }
]
```

**Примеры:**
```bash
# Все активности
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/api/activities

# С фильтрами
curl -H "Authorization: Bearer <JWT>" \
  "http://localhost:8080/api/activities?min_budget=0&max_budget=500&mood=Весело&weather=sunny"
```

### Получить одну активность
**GET** `/activities/{id}`

**Ответ:**
```json
{
  "id": 1,
  "name": "Прогулка в парке",
  "description": "Приятная прогулка на свежем воздухе",
  "budget": 0,
  "time": 2,
  "weather": "sunny",
  "moods": ["Нейтрально", "Хорошо", "Весело"],
  "created_at": "2025-07-10T21:00:00Z"
}
```

### Создать активность (только admin/moderator)
**POST** `/activities`

**Тело запроса:**
```json
{
  "name": "string",
  "description": "string",
  "budget": 0,
  "time": 2,
  "weather": "sunny",
  "moods": ["string"]
}
```

**Ответы:**
- `201 Created` - активность создана
- `403 Forbidden` - недостаточно прав
- `400 Bad Request` - ошибка валидации

### Обновить активность (только admin/moderator)
**PUT** `/activities/{id}`

**Тело запроса:** (аналогично созданию)

**Ответы:**
- `200 OK` - активность обновлена
- `403 Forbidden` - недостаточно прав
- `404 Not Found` - активность не найдена

### Удалить активность (только admin/moderator)
**DELETE** `/activities/{id}`

**Ответы:**
- `204 No Content` - активность удалена
- `403 Forbidden` - недостаточно прав
- `404 Not Found` - активность не найдена

---

## 3. Избранное (Favorites)

### Получить избранное пользователя
**GET** `/favorites`

**Ответ:**
```json
[
  {
    "id": 1,
    "name": "Прогулка в парке",
    "description": "Приятная прогулка на свежем воздухе",
    "budget": 0,
    "time": 2,
    "weather": "sunny",
    "moods": ["Нейтрально", "Хорошо", "Весело"],
    "created_at": "2025-07-10T21:00:00Z"
  }
]
```

### Добавить в избранное
**POST** `/favorites/{activity_id}`

**Ответы:**
- `201 Created` - добавлено в избранное
- `400 Bad Request` - уже в избранном или неверный ID

### Удалить из избранного
**DELETE** `/favorites/{activity_id}`

**Ответы:**
- `204 No Content` - удалено из избранного
- `400 Bad Request` - неверный ID

---

## 4. История просмотров (History)

### Получить последние 10 просмотренных активностей
**GET** `/history`

**Ответ:**
```json
[
  {
    "id": 1,
    "name": "Прогулка в парке",
    "description": "Приятная прогулка на свежем воздухе",
    "budget": 0,
    "time": 2,
    "weather": "sunny",
    "moods": ["Нейтрально", "Хорошо", "Весело"],
    "created_at": "2025-07-10T21:00:00Z"
  }
]
```

### Добавить просмотр активности
**POST** `/history/{activity_id}`

**Ответы:**
- `201 Created` - просмотр добавлен
- `400 Bad Request` - неверный ID

---

## 5. Модели данных

### User
```json
{
  "id": 1,
  "username": "string",
  "role": "user|moderator|admin",
  "created_at": "2025-07-10T21:00:00Z"
}
```

### Activity
```json
{
  "id": 1,
  "name": "string",
  "description": "string",
  "budget": 0,
  "time": 2,
  "weather": "sunny|cloudy|rainy|any",
  "moods": ["string"],
  "created_at": "2025-07-10T21:00:00Z"
}
```

---

## 6. Коды ошибок

### HTTP Status Codes
- `200 OK` - успешный запрос
- `201 Created` - ресурс создан
- `204 No Content` - успешное удаление
- `400 Bad Request` - ошибка валидации
- `401 Unauthorized` - не авторизован
- `403 Forbidden` - недостаточно прав
- `404 Not Found` - ресурс не найден
- `500 Internal Server Error` - ошибка сервера

### Формат ошибок
```json
{
  "error": "описание ошибки"
}
```

---

## 7. Примеры использования

### Типичный flow приложения:

1. **Регистрация/Логин**
```bash
# Регистрация
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "user1", "password": "123456"}'

# Логин
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user1", "password": "123456"}'
```

2. **Получение рекомендаций**
```bash
# Все активности
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/api/activities

# С фильтрами
curl -H "Authorization: Bearer <JWT>" \
  "http://localhost:8080/api/activities?min_budget=0&max_budget=1000&weather=sunny"
```

3. **Работа с избранным**
```bash
# Добавить в избранное
curl -X POST -H "Authorization: Bearer <JWT>" \
  http://localhost:8080/api/favorites/1

# Получить избранное
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/api/favorites
```

4. **История просмотров**
```bash
# Добавить просмотр
curl -X POST -H "Authorization: Bearer <JWT>" \
  http://localhost:8080/api/history/1

# Получить историю
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/api/history
```

---

## 8. Тестовые данные

При первом запуске автоматически создаются:

### Пользователь-админ
- **Username:** `admin`
- **Password:** `admin123`
- **Role:** `admin`

### Примеры активностей
1. **Прогулка в парке** - бюджет: 0, время: 2ч, погода: sunny
2. **Чтение книги** - бюджет: 0, время: 3ч, погода: cloudy  
3. **Кофе с другом** - бюджет: 300, время: 1ч, погода: any

---

## 9. Заметки для разработчиков

- Все временные метки в формате ISO 8601
- Массивы настроений (moods) поддерживают любые строковые значения
- Погода может быть: "sunny", "cloudy", "rainy", "any"
- Роли пользователей: "user", "moderator", "admin"
- Только moderator/admin могут создавать/редактировать/удалять активности
- JWT токен действителен 24 часа 