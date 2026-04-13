# Проверить здоровье сервиса
Invoke-WebRequest -Uri http://localhost:8080/api/v1/health

# Получить информацию
Invoke-WebRequest -Uri http://localhost:8080/api/v1/info

# Получить все задачи
Invoke-WebRequest -Uri http://localhost:8080/api/v1/tasks

# Получить задачу по id
Invoke-WebRequest -Uri http://localhost:8080/api/v1/tasks/1

# Создать новую задачу
Invoke-RestMethod -Uri http://localhost:8080/api/v1/tasks -Method Post -Body '{"title":"New Task", "done":false}' -ContentType "application/json"

# Обновить задачу по id
Invoke-RestMethod -Uri http://localhost:8080/api/v1/tasks/1 -Method Put -Body '{"title":"Updated Title", "done":true}' -ContentType "application/json"

# Удалить задачу по id
Invoke-RestMethod -Uri http://localhost:8080/api/v1/tasks/1 -Method Delete