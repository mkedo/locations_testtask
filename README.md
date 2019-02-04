Запуск веб-сервера
--

```docker-compose up -f docker-compose.tarantool.yml``` (tarantool + vinyl)

или

```docker-compose up -f docker-compose.pgredis.yml``` (postgres + redis)

Схемы
--
PostgreSQL ```postgres/init.sql```

Tarantool ```tarantool/app.lua```

Команды
--
```testtask/cmd/pg``` веб-сервер (postgres + redis)

```testtask/cmd/tarantool``` веб-сервер (tarantool + vinyl)

```testtask/cmd/populate_tarantool``` наполнить случайными данными tarantool

```testtask/cmd/populate_pg``` наполнить случайными данными postgres

Методы
---
##### Связать объявление с адресами

POST /putItemLocations

Пример тела запроса:
```
{
 "ItemId": 5,
 "LocationIds": [1,2,5,46554]
}
```

##### Получить адреса привязанные к объявлению

GET /getItemLocations?ItemId=5

Пример ответа:
```
[
    {
    "ID":46554,
    "Location":"... location json ...",
    "Coordinates":{"X":53.4747,"Y":-14.2082}
    },
    ....
]
```

Допущение
--

Предполагается, что сервис объявлений должен получать LocationId по адресу
либо у этого сервиса, либо у другого.
Получение LocationId по адресу в данном сервисе не реализованно.
