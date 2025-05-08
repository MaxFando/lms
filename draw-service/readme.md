Запросы для проверки API

### 1. Проверка наличие методов 
```
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 list draw_service.v1.DrawService
```

### 2. Добавить тираж 
```
grpcurl -plaintext -d '{
  "lottery_type": "daily",
  "start_time": "2025-05-01T10:00:00Z",
  "end_time": "2025-05-11T10:30:00Z"
}' localhost:50051 draw_service.v1.DrawService/CreateDraw
```

### 3. Посмотреть тиражи 
```
grpcurl -plaintext -d '{}' localhost:50051 draw_service.v1.DrawService/GetDrawsList
```

### 4. Получить законченные тиражи 
```
grpcurl -plaintext -d '{}' localhost:50051 draw_service.v1.DrawService/GetCompletedDrawsList
```

### 5. Получение результата тиража 
```
grpcurl -plaintext -d '{
  "id": 123
}' localhost:50051 draw_service.v1.DrawService/GetDrawResult
```