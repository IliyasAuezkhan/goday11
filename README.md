# 📝 Go JWT ToDo REST API

![Go](https://shields.io)
![Gin](https://shields.io)
![PostgreSQL](https://shields.io)
![Docker](https://shields.io)
![Railway](https://shields.io)

RESTful API для управления задачами с JWT-аутентификацией, PostgreSQL, GORM, Docker и деплоем на Railway.

---

## 🛠️ Стек технологий (Tech Stack)

- **Язык разработки:** Go (Golang)
- **Веб-фреймворк:** Gin Gonic
- **База данных:** PostgreSQL
- **ORM-библиотека:** GORM
- **Авторизация:** JWT (github.com/golang-jwt/jwt/v5) & bcrypt
- **Контейнеризация:** Docker
- **Облачный хостинг:** Railway PaaS

---

## ✨ Возможности (Features)

- **JWT Authentication** — надежная защита приватных эндпоинтов с помощью токенов.
- **Password Hashing (bcrypt)** — безопасное хэширование пользовательских паролей.
- **PostgreSQL** — надежное хранение данных пользователей и задач в реляционной базе.
- **GORM AutoMigrate** — автоматическое создание и обновление таблиц при запуске сервера.
- **RESTful API** — чистая реализация стандартных HTTP-методов (GET, POST, PUT, DELETE).
- **Docker** — многоэтапная сборка (Multi-stage build) для минимизации размера образа.
- **Railway Deployment** — облачный хостинг приложения и базы данных.

---

## 🚀 Живое Демо (Live Deployment)

> 🌐 **Production API URL:**  
> `https://goday11-production.up.railway.app`

*Сервер развернут на платформе Railway и использует облачную базу данных PostgreSQL.*

---

## 📡 Спецификация API (Эндпоинты)

Вся бизнес-логика разделена на публичный сектор и зону, защищенную middleware авторизации (`AuthMiddleware`).

| Метод | Маршрут | Описание Функционала | Уровень Доступа |
| :---: | :--- | :--- | :---: |
| 🟢 `POST` | `/register` | Регистрация нового аккаунта | **`🔐 Публичный`** |
| 🟢 `POST` | `/login` | Аутентификация и выдача JWT-токена | **`🔐 Публичный`** |
| 🔵 `GET` | `/api/todos` | Получить список задач текущего пользователя | **`🔑 JWT Bearer`** |
| 🟢 `POST` | `/api/todos` | Создать задачу (`user_id` автоматически извлекается из JWT) | **`🔑 JWT Bearer`** |
| 🟡 `PUT` | `/api/todos/:id` | Изменить текст или статус `completed` задачи | **`🔑 JWT Bearer`** |
| 🔴 `DELETE` | `/api/todos/:id` | Удалить задачу из базы данных по её ID | **`🔑 JWT Bearer`** |

### 📊 Возвращаемые HTTP статус-коды
- `200 OK` — Запрос успешно выполнен.
- `201 Created` — Ресурс (задача) успешно создан.
- `400 Bad Request` — Неверный формат JSON или синтаксическая ошибка.
- `401 Unauthorized` — Токен отсутствует, некорректен или истек.
- `404 Not Found` — Запрашиваемый маршрут или задача не найдены.
- `500 Internal Server Error` — Внутренняя ошибка сервера или сбой базы данных.

---

## 📖 Пошаговое Руководство по Тестированию (Postman)

### 1. Регистрация Пользователя
- **Метод:** `POST`
- **URL:** `https://goday11-production.up.railway.app/register`
- **Body (raw -> JSON):**
```json
{
  "username": "iliyas_test",
  "password": "supersecurepassword123"
}
```
- **Ответ (`200 OK`):**
```json
{
  "message": "Registered successfully"
}
```

### 2. Вход и Получение Токена
- **Метод:** `POST`
- **URL:** `https://goday11-production.up.railway.app/login`
- **Body (raw -> JSON):** Передайте те же данные.
- **Ответ (`200 OK`):** Скопируйте строку токена.
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Создание Задачи
- **Метод:** `POST`
- **URL:** `https://goday11-production.up.railway.app/api/todos`
- **Headers:** `Authorization: Bearer <ваш_токен>` (Либо вкладка *Authorization -> Bearer Token*)
- **Body (raw -> JSON):**
```json
{
  "title": "Успешно завершить День 15"
}
```
- **Ответ (`201 Created`):**
```json
{
  "id": 1,
  "title": "Успешно завершить День 15",
  "completed": false,
  "user_id": 1
}
```

> ⚠️ **Важно:** При вызове методов `PUT` или `DELETE` заменяйте параметр `:id` в URL-адресе на реальное число (например, `/api/todos/1`), иначе сервер вернет ошибку `404 Not Found`.

---

## 🏗️ Архитектура Проекта

Проект организован по слоям для разделения ответственности между обработчиками, моделями данных и middleware.

```text
goday11/
├── Dockerfile          # Скрипт многоэтапной сборки контейнера
├── go.mod              # Модуль зависимостей Go
├── go.sum              # Контрольные суммы пакетов
└── main.go             # Монолитный файл: модели, middleware, эндпоинты и старт сервера

```

---

## 🐳 Запуск через Docker

Проект можно собрать и запустить локально в изолированном Docker-контейнере.

1. **Сборка Docker-образа:**
   ```bash
   docker build -t go-jwt-todo .
   ```
2. **Запуск контейнера в фоновом режиме:**
   ```bash
   docker run -d -p 8080:8080 --env-file .env go-jwt-todo
   ```

---

## ⚙️ Автоматический деплой (CD)

В проекте настроен непрерывный деплой (CD) через прямую интеграцию **GitHub** и **Railway**:
- Платформа автоматически пересобирает и перезапускает контейнер после каждого `git push` в ветку `main`.
- Сборщик Railway считывает `Dockerfile`, оптимизирует слои и обновляет версию запущенного API.

---

## 💻 Локальный Запуск (Без Docker)

1. **Клонируйте репозиторий и перейдите в папку:**
   ```bash
   git clone https://github.com
   cd goday11
   ```
2. **Создайте файл `.env`** в корне проекта:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=todo_db
   JWT_SECRET=your_super_secret_key
   APP_PORT=8080
   ```
3. **Запустите сервер:**
   ```bash
   go mod tidy
   go run main.go
   ```
