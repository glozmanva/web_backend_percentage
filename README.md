# Лабораторная работа №3

## REST API для сервиса расчёта ежемесячных процентных выплат по банковскому вкладу

Backend реализован на Go с использованием Gin, GORM, PostgreSQL и MinIO.

## Методы API

| Метод | URL | Назначение |
|---|---|---|
| GET | `/api/deposit-months?min_days_count=28&max_days_count=30` | Список опубликованных расчётных месяцев с фильтрацией по количеству дней |
| GET | `/api/deposit-months/feed` | Лента; параметры `id` и `next=true` позволяют получить конкретную или следующую запись |
| GET | `/api/deposit-months/draft` | Получение черновика текущего пользователя |
| POST | `/api/deposit-months` | Создание черновика и загрузка изображения и видео |
| PUT | `/api/deposit-months/draft` | Заполнение полей и публикация черновика |
| DELETE | `/api/deposit-months/:id` | Логическое удаление услуги её создателем |
| POST | `/api/deposit-months/:id/like` | Установка или снятие лайка: `{"liked": 1}` или `{"liked": 0}` |
| POST | `/api/users` | Регистрация пользователя |
| POST | `/api/users/authentication` | Заглушка аутентификации для лабораторной работы №4 |
| POST | `/api/users/deauthentication` | Заглушка деавторизации для лабораторной работы №4 |

## Таблицы базы данных

### `users`

- `id` — первичный ключ;
- `full_name` — ФИО пользователя;
- `birth_date` — дата рождения;
- `email` — электронная почта;
- `password` — пароль.

### `deposit_months`

- `id` — первичный ключ;
- `name` — название расчётного месяца;
- `description` — описание;
- `status` — `draft`, `published` или `deleted`;
- `image_url` — URL изображения;
- `video_url` — URL видео;
- `month_number` — номер месяца;
- `days_count` — количество дней;
- `creator_id` — идентификатор создателя;
- `created_at` — дата создания;
- `formed_at` — дата публикации.

Удаление выполняется логически: запись остаётся в таблице, а её статус меняется на `deleted`.

### `deposit_month_likes`

- `id` — первичный ключ;
- `user_id` — идентификатор пользователя;
- `deposit_month_id` — идентификатор расчётного месяца.

Пара `user_id + deposit_month_id` уникальна.

## Текущий пользователь

В лабораторной работе текущий пользователь фиксирован константой:

```text
currentUserID = 1
```

Получение пользователя выполняется через singleton-функцию `GetCurrentUserID()`.

## Запуск

```bash
docker compose up -d postgres adminer
docker start minio_storage
go run ./cmd/deposit_month
```

API:

```text
http://localhost:8080/api
```

Adminer:

```text
http://localhost:8081
```

MinIO:

```text
http://localhost:9001
```