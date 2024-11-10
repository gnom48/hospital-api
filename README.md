# Проект: информационная система  для  клиник _Simbir.Health_

## Основное задание:
0. Database URL: http://localhost:8085 login=postgres password=postgres
1. Account URL: http://31.129.102.158:8081/swagger/index.html
2. Hospital URL: http://31.129.102.158:8082/swagger/index.html
3. Timetable URL: http://31.129.102.158:8083/swagger/index.html
4. Document URL: http://31.129.102.158:8084/swagger/index.html

## Дополнительное задание:
- Elasticsearch URL: http://localhost:8086
- Kibana URL: http://localhost:8087

## Пояснения
- Модуль models представляет собой встроенную библиотеку, в которой содержится общий для всех микросервисов код. Может дописываться.
- В проекте используется единая база данных, однако есть потенциал для создания распределенной.
- База данных заполнения некоторым количеством тестовых записей.
- Elasticsearch используется для индексирования пользователей по username для того, чтобы разгрузить db и ускорить авторизациию. Используйте поиск в [console](http://localhost:8087/app/dev_tools#/console) для вывода индексов, например:
```
GET /users*/_search
{
  "query": {
    "match_all": {}
  }
}
```

## Запуск
### Для запуска проекта:
1. перейти в корень проекта;
2. ввести команду:
```
docker-compose up --build -d
```
3. проект запущен в фоне.

![Бобер-Golang](https://cs14.pikabu.ru/post_img/2023/05/03/6/og_og_1683105879274332706.jpg)