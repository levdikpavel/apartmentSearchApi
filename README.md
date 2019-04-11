API для иммитации работы с БД
=============================

В БД MySql хранятся квартиры. Параметры подключения к БД задаются в конфиге.

Доступны запросы /search и /add.

Структура сущности таблица:
- Город
- Район
- Адрес
- Жилой комплекс (ЖК)
- Корпус
- Всего этажей
- Этаж
- Комнатность
- Площадь квартиры
- Стоимость аренды


При старте сервиса проверяется подключение, проверяется наличие БД, если нужно - создается. Проверяется наличие таблицы, если нужно - создается.


