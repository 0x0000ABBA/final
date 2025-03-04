## How to use
```shell
$ git clone github.com/0x0000abba/final
$ cd final
$ docker-compose up -d
```
## Makefile
```shell
  - make build - для сборки приложения;
  - make test - для запуска unit-тестов;
  - make docker-build - для сборки Docker-образа с приложением;
  - make run - для запуска приложения;
  - make lint - для запуска линтера;
```