# Лабораторная работа 7 — CI/CD сервиса билетов

![CI](https://github.com/q1wf3/devops-lab7/actions/workflows/ci.yml/badge.svg)
![Docker](https://github.com/q1wf3/devops-lab7/actions/workflows/docker.yml/badge.svg)

Вариант 8: multi-arch образ linux/amd64 и linux/arm64. Корень репозитория — каталог `lab_7`, а Go-модуль — `ticket-service`.

## Проверка локально

```sh
cd ticket-service
go mod download
go mod verify
go vet ./...
go build ./...
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
cd ..
docker build -t ticket-service:v1 .
docker run --rm -p 127.0.0.1:18082:8080 ticket-service:v1
```

Зависимости зафиксированы для учебной матрицы Go 1.21.13/1.22.12 из методички. Это учебные версии, не рекомендация для production.

## Подключение GitHub и Docker Hub

Репозиторий: https://github.com/q1wf3/devops-lab7. CI запускается при push и pull request в main. Docker выполняется после проверок CI.

1. Создать или выбрать доступный репозиторий GitHub; загрузить **содержимое** `lab_7` в его корень. `.github/workflows` должен лежать именно в корне репозитория.
2. Создать репозиторий `ticket-service` в Docker Hub и токен с правом записи в него.
3. В GitHub → Settings → Secrets and variables → Actions добавить `DOCKERHUB_USERNAME` и `DOCKERHUB_TOKEN`. Токен не добавлять в файлы проекта.
4. В Variables добавить `DOCKERHUB_NAMESPACE`, равный логину Docker Hub; это нужно, если логины GitHub и Docker Hub различаются.
5. Дождаться успешных запусков CI и Docker. Docker job имеет `needs: checks` и не публикует образ при упавших тестах; в PR публикация отключена.
6. Настроить защиту `main`: Require a pull request before merging и Require status checks to pass. Выбрать фактически появившиеся проверки Build and Test для обеих версий Go и Lint. Проверить правила на PR с ошибкой, затем исправить его и слить.
7. Скачать artifacts `coverage-go-1.21.13` и `coverage-go-1.22.12`. Сравнить длительность двух одинаковых запусков по журналам Actions; локальные замеры не являются замерами GitHub runner.
8. После создания репозитория заменить `q1wf3/devops-lab7` в шаблонах бейджей:

```md
![CI](https://github.com/q1wf3/devops-lab7/actions/workflows/ci.yml/badge.svg)
![Docker](https://github.com/q1wf3/devops-lab7/actions/workflows/docker.yml/badge.svg)
```

Теги образа: `main`, `sha-<commit>` и `1.0.0` для Git-тега `v1.0.0`. `v*` — фильтр Git-тегов, а semver-тег Docker не содержит начальную `v`.

Для публикации вручную после входа:

```sh
docker login
# Заменить YOUR_DOCKERHUB_NAME на собственный логин
docker tag ticket-service:v1 YOUR_DOCKERHUB_NAME/ticket-service:v1
docker push YOUR_DOCKERHUB_NAME/ticket-service:v1
docker pull YOUR_DOCKERHUB_NAME/ticket-service:v1
```

## Проверка защиты main

В PR #1 намеренно добавлена локальная переменная `unused`, которая не использовалась. Проверки Build and Test и Lint завершились ошибкой; кнопка Merge pull request была недоступна. Ошибка удалена следующим коммитом. Правило main требует успешных проверок Go 1.21.13, Go 1.22.12 и Lint, в том числе для администратора.

При отсутствии секретов Docker Hub workflow выполняет сборку для двух архитектур, но явно отмечает пропуск публикации. Для публикации добавьте секреты и запустите Docker через Run workflow на main.
