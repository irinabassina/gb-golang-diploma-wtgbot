# gb-golang-diploma-wtgbot

Дипломная работа по направлению "Программирование", специализация "Разработчик GoLang".

Warehouse-diploma-bot - корпоративный Telegram чат-бот для ведения автоматизированного учета товарных запасов на складе,
который позволяет анализировать и прогнозировать операции движения товаров при помощи ИИ.

### Необходимые для работы переменные окружения

* TOKEN_TG : токен для Telegram бота
* TOKEN_OPEN_AI : токен для Openai

Также для работы бота необходима БД Postgres. Можно запустить БД в отдельном локальном контейнере или подключиться к
существующей БД.
Параметры подключения бота к БД по-умолчанию, а также название переменных окружения для управления подключением
находятся в файле [internal/env/config.go](internal/env/config.go)

Файлы инициализации Postgres БД можно найти здесь [migrations](/migrations).

### Пример запуска

``` bash
echo TOKEN_TG=YOUR_TG_BOT_TOKEN >> .env; \
echo TOKEN_OPEN_AI=YOUR_OPEN_AI_TOKEN >> .env; \
docker-compose up -d
```
