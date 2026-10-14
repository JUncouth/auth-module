# auth-module

Отдельный сервис авторизации для appliance-store на Go.

## Функционал

- регистрация пользователей
- проверка пары логин/пароль
- выдача и проверка JWT
- хеширование паролей с использованием bcrypt
- межсервисное взаимодействие через gRPC
- описание интерфейса сервиса с помощью protobuf
- хранение данных в postgreSQL

## Интеграция с appliance-store

- генерация proto client для соответствующего языка
- настройка подключения к сервису auth
- добавление middleware и эндпоинтов

## Технологии

- GO, gRPC, protobuf, PostgreSQL, JWT, bcrypt, uber/fx (dependency injection)