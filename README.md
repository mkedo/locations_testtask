Запуск сервиса
--

```docker-compose up -f docker-compose.tarantool.yml``` (tarantool + vinyl)

или

```docker-compose up -f docker-compose.pgredis.yml``` (postgres + redis)

Основные файлы и директории
--
```testtask/testtask``` исходный код сервиса

```postgres/init.sql``` схема PostgreSQL

```tarantool/app.lua``` схема Tarantool

Команды
--
```testtask/cmd/pg``` веб-сервер (postgres + redis)

```testtask/cmd/tarantool``` веб-сервер (tarantool + vinyl)

```testtask/cmd/populate_tarantool``` наполнить случайными данными tarantool

```testtask/cmd/populate_pg``` наполнить случайными данными postgres

Методы
---

|url|Метод|Описание|
|---|------|-------|
|/item/{ItemId}/locations|POST| Связать объявление с адресами |

Пример тела запроса:
```
{
 "LocationIds": [1,2,5,46554]
}
```

|url|Метод|Описание|
|---|------|-------|
|/item/{ItemId}/locations|GET| Получить адреса привязанные к объявлению |

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
