# gb-go-backend-1
## lesson-8
### Task1

_Попробуйте “докеризировать” Go-приложения, написанные вами ранее. Обратите внимание на то, что в контейнеры можно положить не только сервисы, но и cli-приложения. Для лучшего закрепления навыков можно попрактиковаться и с теми, и с другими видами приложений._

Впервые познакомился с Docker в 2020 и даже делал тестовое задание - простое приложение, которое читает конфиг фтп, подключается туда и делает подсчет ascii символов из файлов: https://github.com/alextonkonogov/ascii_counter 
Также уже пробовал докеризировать Go-приложения на предыдущих курсах. Например: https://github.com/alextonkonogov/gb-go-observability/tree/homework2/homework2

---
### Task2

_Для приложений, которые работали с подключаемыми ресурсами (например, с базами данных) попробуйте использовать инструмент docker compose._

Моя первая версия курсового проекта уже работает через docker compose: https://github.com/alextonkonogov/gb-go-url-shortener

---
### Task3

_Если вы уже сталкивались с докеризацией Go-приложений, проведите аудит. Соответствуют ли ваши докерфайлы лучшим практикам? Все ли полезные параметры заданы в файлах конфигурации docker compose?_

Очень полезная лекция была, я узнал много нового. Например, что нужно собирать бинарник непосредственно в контейнере и что сами контейнеры еще делятся на слои. Также не знал, что переменные окружения можно в некоторых случаях там же задавать.
Обязательно буду использовать новые знания, когда займусь переработкой курсового проекта.

---
### Task4

_Запустите и проверьте (curl запросами) что система из рассмотренного на занятии примера - работает (по аналогии с предыдущими ДЗ)_

Запускаем docker compose:

```shell
a.tonkonogov@admins-MacBook-Pro lesson-8 % docker compose up
[+] Running 3/2
 ⠿ Container postgres  Created                                                                                                                                                                                                    0.1s
 ⠿ Container registry  Created                                                                                                                                                                                                    0.1s
 ⠿ Container reguser   Created                                                                                                                                                                                                    0.1s
Attaching to postgres, registry, reguser
Error response from daemon: Failed to inspect container c8bf90859d0ca6ae2254bf146723b04cef9204b36c58933e87b6664311c7658d: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
a.tonkonogov@admins-MacBook-Pro lesson-8 % docker compose up
[+] Running 4/4
 ⠿ Network lesson-8_regusernet  Created                                                                                                                                                                                           0.0s
 ⠿ Container registry           Created                                                                                                                                                                                           0.1s
 ⠿ Container postgres           Created                                                                                                                                                                                           0.1s
 ⠿ Container reguser            Created                                                                                                                                                                                           0.1s
Attaching to postgres, registry, reguser
registry  | time="2022-06-19T17:27:57.784723137Z" level=warning msg="No HTTP secret provided - generated random secret. This may cause problems with uploads if multiple registries are behind a load-balancer. To provide a shared secret, fill in http.secret in the configuration file or set the REGISTRY_HTTP_SECRET environment variable." go.version=go1.16.15 instance.id=54e82276-90ca-45d3-a7ac-7abd14388fb2 service=registry version="v2.8.1+unknown" 
registry  | time="2022-06-19T17:27:57.785099804Z" level=info msg="redis not configured" go.version=go1.16.15 instance.id=54e82276-90ca-45d3-a7ac-7abd14388fb2 service=registry version="v2.8.1+unknown" 
registry  | time="2022-06-19T17:27:57.784869929Z" level=info msg="Starting upload purge in 32m0s" go.version=go1.16.15 instance.id=54e82276-90ca-45d3-a7ac-7abd14388fb2 service=registry version="v2.8.1+unknown" 
registry  | time="2022-06-19T17:27:57.791448304Z" level=info msg="using inmemory blob descriptor cache" go.version=go1.16.15 instance.id=54e82276-90ca-45d3-a7ac-7abd14388fb2 service=registry version="v2.8.1+unknown" 
registry  | time="2022-06-19T17:27:57.791678512Z" level=info msg="listening on [::]:5000" go.version=go1.16.15 instance.id=54e82276-90ca-45d3-a7ac-7abd14388fb2 service=registry version="v2.8.1+unknown" 
postgres  | The files belonging to this database system will be owned by user "postgres".
postgres  | This user must also own the server process.
postgres  | 
postgres  | The database cluster will be initialized with locale "en_US.utf8".
postgres  | The default database encoding has accordingly been set to "UTF8".
postgres  | The default text search configuration will be set to "english".
postgres  | 
postgres  | Data page checksums are disabled.
postgres  | 
postgres  | fixing permissions on existing directory /var/lib/postgresql/data ... ok
postgres  | creating subdirectories ... ok
postgres  | selecting dynamic shared memory implementation ... posix
postgres  | selecting default max_connections ... 100
postgres  | selecting default shared_buffers ... 128MB
postgres  | selecting default time zone ... Etc/UTC
postgres  | creating configuration files ... ok
postgres  | running bootstrap script ... ok
postgres  | performing post-bootstrap initialization ... ok
postgres  | syncing data to disk ... ok
postgres  | 
postgres  | 
postgres  | Success. You can now start the database server using:
postgres  | 
postgres  |     pg_ctl -D /var/lib/postgresql/data -l logfile start
postgres  | 
postgres  | initdb: warning: enabling "trust" authentication for local connections
postgres  | You can change this by editing pg_hba.conf or using the option -A, or
postgres  | --auth-local and --auth-host, the next time you run initdb.
postgres  | waiting for server to start....2022-06-19 17:27:58.438 UTC [49] LOG:  starting PostgreSQL 12.11 (Debian 12.11-1.pgdg110+1) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 10.2.1-6) 10.2.1 20210110, 64-bit
postgres  | 2022-06-19 17:27:58.440 UTC [49] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
postgres  | 2022-06-19 17:27:58.463 UTC [50] LOG:  database system was shut down at 2022-06-19 17:27:58 UTC
postgres  | 2022-06-19 17:27:58.467 UTC [49] LOG:  database system is ready to accept connections
postgres  |  done
postgres  | server started
postgres  | CREATE DATABASE
postgres  | 
postgres  | 
postgres  | /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/init_users.sql
reguser   | 2022/06/19 20:27:58 failed to connect to `host=postgres user=postgres database=test`: dial error (dial tcp 172.26.0.3:5432: connect: connection refused)
postgres  | CREATE TABLE
postgres  | 
postgres  | 
postgres  | waiting for server to shut down....2022-06-19 17:27:58.755 UTC [49] LOG:  received fast shutdown request
postgres  | 2022-06-19 17:27:58.756 UTC [49] LOG:  aborting any active transactions
postgres  | 2022-06-19 17:27:58.760 UTC [49] LOG:  background worker "logical replication launcher" (PID 56) exited with exit code 1
postgres  | 2022-06-19 17:27:58.761 UTC [51] LOG:  shutting down
postgres  | 2022-06-19 17:27:58.778 UTC [49] LOG:  database system is shut down
reguser exited with code 1
postgres  |  done
postgres  | server stopped
postgres  | 
postgres  | PostgreSQL init process complete; ready for start up.
postgres  | 
postgres  | 2022-06-19 17:27:58.861 UTC [1] LOG:  starting PostgreSQL 12.11 (Debian 12.11-1.pgdg110+1) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 10.2.1-6) 10.2.1 20210110, 64-bit
postgres  | 2022-06-19 17:27:58.861 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
postgres  | 2022-06-19 17:27:58.861 UTC [1] LOG:  listening on IPv6 address "::", port 5432
postgres  | 2022-06-19 17:27:58.863 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
postgres  | 2022-06-19 17:27:58.872 UTC [86] LOG:  database system was shut down at 2022-06-19 17:27:58 UTC
postgres  | 2022-06-19 17:27:58.876 UTC [1] LOG:  database system is ready to accept connections
reguser   | 2022/06/19 20:27:59 Start

```

Тестируем приложение:

---
#### Создание

**Запрос:**
```shell
curl --location --request POST 'localhost:8000/create' \
--header 'Authorization: Basic YWRtaW46YWRtaW4=' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "George",
    "data": "Lorem Ipsum",
    "perms": 765
}'
```

**Ответ:**
```shell
{"id":"49e2f306-5861-45c3-8beb-6e60cbc6310e","name":"George","data":"Lorem Ipsum","perms":0}
```

--
#### Поиск

**Запрос:**
```shell
curl --location --request GET 'localhost:8000/search/e' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```

**Ответ:**
```shell
[
{"id":"49e2f306-5861-45c3-8beb-6e60cbc6310e","name":"George","data":"Lorem Ipsum","perms":493}
]
```

---
#### Чтение

**Запрос:**
```shell
curl --location --request GET 'localhost:8000/read/49e2f306-5861-45c3-8beb-6e60cbc6310e' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```

**Ответ:**
```shell
{"id":"49e2f306-5861-45c3-8beb-6e60cbc6310e","name":"George","data":"Lorem Ipsum","perms":0}
```
---
#### Удаление

**Запрос**:
```shell
curl --location --request DELETE 'localhost:8000/delete/49e2f306-5861-45c3-8beb-6e60cbc6310e' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```

**Ответ:**
```shell
{"id":"49e2f306-5861-45c3-8beb-6e60cbc6310e","name":"George","data":"Lorem Ipsum","perms":0}
```
---
#### Проверим БД
Проверим, что в БД у нас все записывается и удаляется:

```shell
a.tonkonogov@admins-MacBook-Pro gb-go-backend-1 % psql -h localhost -p 5432 -d test -U postgres
Password for user postgres: 
psql (14.1, server 12.11 (Debian 12.11-1.pgdg110+1))
Type "help" for help.

test=# \d
         List of relations
 Schema | Name  | Type  |  Owner   
--------+-------+-------+----------
 public | users | table | postgres
(1 row)

test=# select * from users;
                  id                  |          created_at           |          updated_at          |          deleted_at           |  name  |    data     | perms 
--------------------------------------+-------------------------------+------------------------------+-------------------------------+--------+-------------+-------
 49e2f306-5861-45c3-8beb-6e60cbc6310e | 2022-06-19 17:47:47.132506+00 | 2022-06-19 17:47:47.13252+00 | 2022-06-19 17:50:04.671997+00 | George | Lorem Ipsum |     0
(1 row)

```

---
Помним про тот же баг, который я правил в 6-ой домашке: perms не записывается при создании и потом выдается в виде 0.
Снова не стал его править, потому что это не влияет на суть задания.
