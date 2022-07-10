# Курсовой проект URL Shortener

## Описание работы сервиса

- При заходе на стартовую страницу (эндпоинт `/`) видим форму с одним полем для ввода длинной ссылки и кнопку для ее отправки.
- Указываем ссылку и нажимаем на кнопку.
- Без перезагрузки страницы запрос улетает на сервер, обрабатывается и возвращает JSON c результатом.
- Фронт обрабатывает JSON и в специальном div формирует HTML контент: исходная ссылка, короткая ссылка и админская ссылка.
- Если переходим по короткой ссылке и она существует в БД, то сервер перенаправит на длинную.
- При каждом переходе по короткой ссылке в БД фиксируется время, IP и инкрементируется счетчик посещений.
- Если переходим по админской ссылке и она существует в БД, то сервер сгенерирует статистику и вернет JSON, который на лету обработает фронт и сгенерирует HTML
- Если в короткой ссылке сделать ошибку (убрать или добавить символы), то сервер перенаправляет на специальную страницу с ошибкой (эндпоинт `/err`). На неверные админские ссылки возвращается JSON с ошибкой, перенаправление на эндпоинт `/err` ложится на логику фронта. 

В качестве фронтенд был использован bootstrap и дополнительно написан простой JS код. Авторизацию не добавлял умышленно, хотя умею - сама логика сервиса не предполагает ее наличие, как мне кажется.

Постарался реализовать возможность работы по REST API, чтобы удобно было тестировать и использовать в связке с другими клиентами: бот, мобильное приложение, Postman и тд. Описание с примерами будет в конце.

Текущая версия проекта развернута у меня на сайте (платформа NetAngels) - [gbt.alextonkonogov.ru](https://gbt.alextonkonogov.ru/)

![image](img.png)

![image](img_1.png)

Также есть возможность запуститься локально через Docker Compose.

![image](img_2.png)

Тесты планирую добавить завтра.
Код в ветке final_v: https://github.com/alextonkonogov/gb-go-url-shortener/tree/final_v

---
#### СОЗДАНИЕ
Отправляем JSON методом POST на эндпоинт `/s/create`

**Запрос:**
```shell
curl --location --request POST 'http://localhost:8000/s/create' \
--header 'Content-Type: application/json' \
--data-raw '{
    "long":"https://gbcdn.mrgcdn.ru/uploads/asset/3001858/attachment/c3640e219eb26045352728efea6d443e.pdf"
}'
```

**Ответ:**

Если все ок:
```shell
{
    "id": 4,
    "created": "2022-07-10T03:16:29.101Z",
    "long": "https://gbcdn.mrgcdn.ru/uploads/asset/3001858/attachment/c3640e219eb26045352728efea6d443e.pdf",
    "short": "w9fio2NI7hGLZ6XX",
    "admin": "g53tUw73xehbhC1I"
}
```

Если пустой JSON
```shell
{
"status": "Invalid request.",
"error": "render: unable to automatically decode the request content type"
}
```

Если URL невалидный:
```shell
{
"status": "Invalid request.",
"error": "invalid URL"
}
```

---
#### ЧТЕНИЕ (короткой ссылки)
GET запрос на эндпоинт `/s/{short}` с коротким кодом в составе пути

**Запрос:**
```shell
curl --location --request GET 'http://localhost:8000/s/w9fio2NI7hGLZ6XX'
```

**Ответ:**
Редирект либо на длинную ссылку, либо на страницу с ошибкой (эндпоинт `/err`)

---
#### ЧТЕНИЕ (админской ссылки)
POST запрос на эндпоинт `/a` с JSON содержимым. В ответ получим JSON со статистикой

**Запрос:**
```shell
curl --location --request POST 'http://localhost:8000/a' \
--header 'Content-Type: application/json' \
--data-raw '{
    "admin":"g53tUw73xehbhC1I"
}'
```

**Ответ:**

Если все ок:
```shell
{
    "ip": "172.28.0.1:56220",
    "viewed": "2022-07-10T03:20:15.403Z",
    "count": 3,
    "long": "https://gbcdn.mrgcdn.ru/uploads/asset/3001858/attachment/c3640e219eb26045352728efea6d443e.pdf",
    "short": "w9fio2NI7hGLZ6XX",
    "admin": ""
}
```
Если отправим пустой JSON:
```shell
{
"status": "Invalid request.",
"error": "render: unable to automatically decode the request content type"
}
```
Если укажем неверный код админа:
```shell
{
"status": "Resource not found.",
"error": "error when reading: read statistics error: sql: no rows in result set"
}
```
