# Микросервис для работы с балансом пользователя

Проект создан при выполнении [тестового задания от Avito](https://github.com/avito-tech/internship_backend_2022)

## Запуск

### Используемый порт
По умолчанию сервис использует порт 9096. Номер порта для сервиса указывается в 2-х файлах:
- [./docker-compose.yaml](https://github.com/adepte-myao/avito_internship/blob/master/docker-compose.yaml)
- [./config/config.yaml](https://github.com/adepte-myao/avito_internship/blob/master/config/config.yaml)

Для нормальной работы сервиса порты, указанные в них, должны совпадать.

### Уровень логгирования
По умолчанию сервис логгирует информацию уровня "debug" и выше. Уровень логгирования задается в файле:
[./config/config.yaml](https://github.com/adepte-myao/avito_internship/blob/master/config/config.yaml)
Возможные значения: "debug", "info", "warning", "error", "fatal", "panic".

### Заполнение базы данных
Для проверки функций get_account_statement и get_accountant_report базу данных необходимо наполнить. 
Для этой цели в проекте существует [sql-скрипт](https://github.com/adepte-myao/avito_internship/blob/master/build/docker/db/ps_fill.sql)
По умолчанию он не включен в процесс сборки. Для включения в конец 
[Dockerfile базы данных](https://github.com/adepte-myao/avito_internship/blob/master/build/docker/db/Dockerfile) добавьте строчку
```Dockerfile
COPY ps_fill.sql /docker-entrypoint-initdb.d/11-init.sql
```

### Запуск
```ShellSession
docker-compose up
```

## Тестирование

### Особенность тестирования

Для тестирования [одного из компонентов приложения](https://github.com/adepte-myao/avito_internship/tree/master/internal/storage) необходима работающая база данных.
Есть два варианта организации базы данных:

1.  Использовать существующую базу данных. Для этого в [файле](https://github.com/adepte-myao/avito_internship/blob/master/internal/storage/storage_test.go) нужно
    заменить строку
    ```go
    databaseURL = "postgres://balancer:superpass@localhost:5434/tests?sslmode=disable"
    ```
    на строку
    ```go
    databaseURL = "postgres://balancer:superpass@localhost:5434/userbalances?sslmode=disable"
    ```
    Побочный эффект: почти все таблицы в основной базе данных после тестирования будут содержать мусор, полезные данные из них удалятся.

2.  Создать новую базу данных в том же контейнере с именем tests, и выполнить в ней 
    [скрипт](https://github.com/adepte-myao/avito_internship/blob/master/build/docker/db/init.sql).

### Запуска тестирования

Запустить тестирование из корневого каталога можно 2-мя способами:
```ShellSession
make test
```
Если не работает указанный выше способ:
```ShellSession
go test -v -race -timeout 30s ./...
```

## Конечные точки

### /ping (get)
Проверяет доступность сервиса извне.

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200500328-37d238d2-e4d5-4cf4-b722-a7de0a26f785.png)

### /swagger/index.html (open in browser)

Swagger-документация для сервиса.

![image](https://user-images.githubusercontent.com/106271382/200500741-57276dff-ea0e-4267-b9ae-09f2b57a2d1d.png)

### /balance/deposit (post)

Пополняет баланс существующего аккаунта. Если аккаунта не существует, создает новый и пополняет его. В случае успеха возвращает код 204 (NoContent).

Пример запрос-ответа (аккаунта не существовало):
![image](https://user-images.githubusercontent.com/106271382/200501364-edea53d9-1711-4b74-8eaa-9e93cb642674.png)

Пример запрос-ответа (аккаунт существует):
![image](https://user-images.githubusercontent.com/106271382/200501491-c1146e82-002f-4761-8aa6-d59f9660c88c.png)

### /balance/withdraw (post)

Уменьшает баланс, если аккаунт существует и денег на нем достаточно. В случае успеха возвращает код 204 (NoContent).

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200502002-37586d94-678c-471c-8f01-d087115f794e.png)

### /balance/get (get)

Возвращает баланс для указанного аккаунта.

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200502451-658360d0-fab5-4195-94b1-5e17d3469af3.png)

### /balance/transfer (post)

Переводит деньги между счетами, если оба аккаунта существуют и денег на аккаунте отправителя достаточно. В случае успеха возвращает код 204 (NoContent).

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200502903-89fe6daf-df38-4bef-a910-2c347fed6cc6.png)

### /balance/statement (get)

Формирует и отправляет список транзакций пользователя, отсортированный по указанным критериям. 
Если указан только первый критерий, сортировка гарантирована только по нему.
Если критерии не указаны, сортировка будет только по времени фиксирования транзакции.
Доступные критерии: "record_time" (время фиксирования транзакции), "amount" (сумма транзакции).

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200504216-8309d0d5-652d-475e-b4bb-4086dc2d2922.png)

### /reservation/make (post)

Резервирует деньги с указанного аккаунта при выполнении следующих условий:
- указанный аккаунт существует
- резервации с такими же параметрами-идентификаторами не существует
- сумма, указанная в теле запроса, не больше баланса пользователя

В случае успеха возвращает код 204 (NoContent).

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200504912-96668e22-6e21-4acb-bbc6-071b029129c9.png)

### /reservation/accept (post)

Принимает резервацию при следующих условиях:
- существует запись о резервации с указанными параметрами
- не существует записи об отклонении резервации с указанными параметрами

В случае успеха возвращает код 204 (NoContent).

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200505598-2c957854-93fa-43b5-bf6a-5daa8f414a24.png)

### /reservation/cancel (post)

Отклоняет резервацию и возвращает деньги на аккаунт пользователя при следующих условиях:
- существует запись о резервации с указанными параметрами
- не существует записи о принятии резервации с указанными параметрами

В случае успеха возвращает код 204 (NoContent).

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200505827-0afe9cd2-3fd9-4c55-8fbf-a8f6483c767b.png)


### /accountant-report (get)

Формирует отчет в формате "название сервиса" - "суммарный доход" за указанный месяц в году.
Пока что возвращает непосредственно данные, планируется возвращать ссылку на .csv файл с данными.

Пример запрос-ответа:
![image](https://user-images.githubusercontent.com/106271382/200510640-de417591-3b9d-46d6-b983-7c7c5fb8ef0a.png)


