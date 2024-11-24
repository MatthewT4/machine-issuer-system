# Платформа по управлению вычислительными ресурсами

## Инструкция по поднятию
1. В директорию склонировать backend и frontend части
```git clone https://github.com/SemenovRoman12/kafedra.git```
```git clone https://github.com/MatthewT4/machine-issuer-system.git```
2. Перейти в backend часть
```cd machine-issuer-system```
3. Поднять docker контейнеры
```docker-compose up -d```

### Информация по директориям
api: swagger документация по сервису
infra: Утилиты и серверы для работы с виртуальными машинами
postgres: Миграции и докерфайл для базы

### Инструкция по использованию
1. Фронт доступен по адресу http://localhost:4200
2. Для работы с сервисом нужно авторизоваться
3. По адресу http://localhost:3000 доступны дашборды в grafana