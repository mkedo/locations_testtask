Архитектура
--

Реализовано 3 варианта архитектуры:

1. **tarantool + vinyl** - адреса и их связь с объявлениями хранятся в tarantool.
2. **postgres + redis (как кеш)** - адреса и их связь с объявлениями хранятся в postgres. Redis используется как кеш запросов.
3. **postgres + redis (как постоянное хранилище)** - адреса хранятся в postgres. Связи объявлений с адресами хранятся в Redis.


Запуск сервиса
--

```docker-compose up -f docker-compose.tarantool.yml``` (tarantool + vinyl)

или

```docker-compose up -f docker-compose.pgredis.yml``` (postgres + redis как кеш)

или

```docker-compose up -f docker-compose.redis_persistent.yml``` (postgres + redis как постоянное хранилище)


Основные файлы и директории
--
```testtask/testtask``` исходный код сервиса

```postgres/init.sql``` схема PostgreSQL

```tarantool/app.lua``` схема Tarantool

Команды
--
```testtask/cmd/pg``` веб-сервер (postgres + redis в качестве кеша)

```testtask/cmd/redis_persistent``` веб-сервер (postgres + redis в качестве постоянного хранилища)

```testtask/cmd/tarantool``` веб-сервер (tarantool + vinyl)

```testtask/cmd/populate_tarantool``` наполнить случайными данными tarantool

```testtask/cmd/populate_pg``` наполнить случайными данными (postgres + redis в качестве кеша)

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
Получение LocationId по адресу в данном сервисе не реализовано.
