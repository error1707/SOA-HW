### Домашняя работа 6
## REST сервис

Для запуска сервеной части воспользуйтесь docker-compose файлом

1. `docker-compose up -d rabbitmq postgres`
2. Подождать пару секунд, пока rabbitmq и postgres не поднимутся
3. `docker-compose run pgmigrate`
4. `docker-compose up -d server worker`

Далее открыть в браузере http://localhost:8080

Интерфейс интуитивно понятен
