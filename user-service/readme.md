## Основные методы (gRPC)


1. **Register** (RegisterRequest → RegisterResponse)  
   – Создаёт нового пользователя, хеширует пароль, сохраняет в БД, генерирует `access` + `refresh` токены, публикует событие `user.created`.

2. **Login** (LoginRequest → LoginResponse)  
   – Аутентифицирует по имени и паролю, обновляет refresh-токен, публикует событие `user.logged_in`.

3. **GetUser** (GetUserRequest → UserResponse)  
   – Возвращает данные пользователя по ID.

4. **ListUsers** (Empty → ListUsersResponse)  
   – Возвращает список всех пользователей.

## Запросы для проверки API

### 1. Проверка наличие методов
```
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 list user-service.v1.UserService
```
### 2. Регистрация
```
grpcurl -plaintext -d '{
  "name": "alice",
  "password": "P@ssw0rd!"
}' localhost:50051 user-service.v1.UserService/Register
```
### 3. Логин
```
grpcurl -plaintext -d '{
  "name": "alice",
  "password": "P@ssw0rd!"
}' localhost:50051 user-service.v1.UserService/Login
```

### 4. Получение пользователя по ID
```
grpcurl -plaintext -d '{
  "id": 42
}' localhost:50051 user-service.v1.UserService/GetUser
```

### 5. Список всех пользователей
```
grpcurl -plaintext -d '{}' localhost:50051 user-service.v1.UserService/ListUsers
```


## Модули
1.  repository (internal/repository):

- Работа с PostgreSQL через github.com/jmoiron/sqlx.

- Интерфейс UserRepository и транзакционные методы BeginTx, CreateTx, UpdateRefreshTokenTx.

2. service (internal/service)
– Бизнес-логика в UserService:
- Register: проверка уникальности, bcrypt-хеширование, транзакции, генерация JWT, публикация события.
- Login: проверка имени/пароля, генерация новых токенов, обновление refresh-токена.
- GetUser, ListUsers: простые вызовы репозитория.

3. jwt (internal/jwt):
- Интерфейс JWTService (генерация Access/Refresh, валидация).
- Позволяет подменить реализацию (например, Auth-Gateway).

4. pubsub (internal/pubsub):
- Интерфейс PubSub для публикации в Redis-каналы:
- PublishUserCreated(ctx, userID)
- PublishUserLoggedIn(ctx, userID)

5. server (internal/server):пше
- Инициализация gRPC-сервера с panic-рекавери-интерсептором.
- Регистрация UserServiceServer и включение reflection.

