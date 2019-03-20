# temperature-backend
[![Build Status](https://travis-ci.org/ShakirzyanovArsen/temperature-backend.svg?branch=master)](https://travis-ci.org/ShakirzyanovArsen/temperature-backend)

Запуск приложения без контейнеризации: _go run main.go_

Запуск тестов(если передать флаг short то запустятся все тесты кроме
 функциональных из пакета functional_test): 
_go test ./..._  

Сборка и запуск приложения в docker контейнере:
1. _docker build -t temperature-backend ._
2. _docker run -p 8080:8080 temperature-backend_

Сборка и запуск тестов в docker контейнере:
1. _docker build -t temp-back-tests -f DockerfileTest ._
2. _docker run temp-back-tests_

В директории requests есть примеры запросов к сервису.