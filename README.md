## Запуск проекта

### Установка и запуск

1. Клонируйте репозиторий:

   ```sh
   git clone https://github.com/Blxssy/AvitoTest.git
   cd AvitoTest
   ```
2. Запустите сервис:

   ```sh
   docker-compose up -d --build
   ```

Сервер будет запущен на `http://localhost:8080`.

## API Эндпоинты

### 1. Авторизация

#### `POST /api/auth`

Авторизует пользователя и возвращает токен.

**Запрос:**

```json
{
  "username": "user1",
  "password": "securepass"
}
```

**Ответ:**

```json
{
  "token": "access_token"
}
```

### 2. Перевод монет

#### `POST /api/sendCoin`

Отправляет монеты другому пользователю.

**Заголовки:**

- `Authorization: Bearer <token>`

**Запрос:**

```json
{
  "toUser": "user2",
  "amount": 50
}
```

**Ответ:**

- `200 OK` (успешный перевод)

### 3. Получение информации о пользователе

#### `GET /api/info`

Возвращает баланс, инвентарь и историю транзакций пользователя.

**Заголовки:**

- `Authorization: Bearer <token>`

**Ответ:**

```json
{
  "coins": 1030,
  "inventory": [
    { "type": "powerbank", "quantity": 1 }
  ],
  "coinHistory": {
    "received": [
      { "fromUser": "user2", "amount": 50 }
    ],
    "sent": [
      { "toUser": "user2", "amount": 20 }
    ]
  }
}
```

### 4. Покупка предметов

#### `GET /api/buy/:item`

Позволяет пользователю купить предмет за монеты.

**Заголовки:**

- `Authorization: Bearer <token>`

**Пример запроса:**

```sh
GET /api/buy/powerbank
```

**Ответ:**

- `200 OK` (покупка успешна)
- `500 Internal Server Error` (ошибка покупки)

