# Проект: информационная система  для  клиник _Simbir.Health_

# Основное задание:
0. Database URL: http://localhost:8085 login=postgres password=postgres
1. Account URL: http://localhost:8081/ui-swagger
2. Hospital URL: http://localhost:8082/ui-swagger
3. Timetable URL: http://localhost:8083/ui-swagger
4. Document URL: http://localhost:8084/ui-swagger

# Дополнительное задание:
- Elasticsearch URL: http://localhost:8086
- Kibana URL: http://localhost:8087

# Пояснения
- Модуль models представляет собой встроенную библиотеку, в которой содержится общий для всех микросервисов код. Может дописываться.
- В проекте используется единая база данных, однако есть потенциал для создания распределенной.
- База данных заполнения некоторым количеством тестовых записей.
- Elasticsearch используется для индексирования пользователей по username для того, чтобы разгрузить db и ускорить авторизациию

# Запуск
Для запуска проекта:
1. перейти в корень проекта;
2. ввести команду:
```
docker-compose up --build -d
```
3. проект запущен в фоне.