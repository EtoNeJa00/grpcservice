# grpcservice

gRPC сервис, который позволяет хранить данные, а также изменять и удалять их, используя в качестве хранилища
сервис memcahed и внутреннее хранилище на случай, если memcached окажется недоступен.

## Зависимости
docker==1:20.10.23-1
docker-compose==2.15.1-1

## Запуск
`docker-compose up -d`

## Доступные команды
Для общения с сервисом необходим gRPC клиент.
API можно получить как через рефлексию, так и через импорт protobuf.

### Создать/обновить записи
Команда `set`
`{
    "data":"message",
    "id": "00000000-0000-0000-0000-000000000000"
}`

Если uuid нулевой, то будет создана новая запись, иначе будет обновлена старая.

### Получить записи
Команда `get`
`{
    "id": "7f6cadd4-0443-433d-a790-650f3bbb11ba"
}`

Возвращает запись.

### Удалить записи
Команда `delete`
`{
    "id": "7f6cadd4-0443-433d-a790-650f3bbb11ba"
}`

Удаляет запись.
