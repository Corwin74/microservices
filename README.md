# Домашнее задание №2

## Задание

 - Создать моно репозиторий (исключительно для удобства) для сервисов из Домашнего задания 1.
 - Для каждого сервиса создать main с http-сервером c liveness/readiness probe (можно использовать как стандартную библиотеку golang net/http так и любой понравившийся вам фреймворк)
 - Написать Dockerfile для каждого сервиса (естественно это все должно собираться).
 - Написать инструкцию или скрипт для того, чтобы можно было поднять все сервисы в контейнерах локально. (Подсказка: для удобства локальной разработки лучше всего воспользоваться docker-compose и Makefile)

 - ⭐ Реализовать стратегии деплоя blue-green и canary с помощью стандартных средств kubernetes.

# Поднятие сервисов локально (Docker)

Для сборки и поднятия сервисов в докере необходимо выполнить следующие команды:

```sh
cd services
docker compose build
docker compose up
```
Сервисы будут доступны на localhost, порты: 8081, 8082, 8083, 8084

# Реализация стратегии деплоя blue-green 

Для демонстрации стратегий нам потребуется установленный minikube. Все последующие команды выполняются из каталога `services`. Выполняем сборку "зеленой" и "синей" версии сервиса:

```sh
docker build -t auth:blue -f auth/build/blue/Dockerfile .
docker build -t auth:green -f auth/build/green/Dockerfile .
```

Затем загружаем эти образы в registry minikube:

```sh
minikube image load auth:blue
minikube image load auth:green
```

Создаем namespace `messenger`

```sh
kubectl apply -f deployment/namespace.yaml
```

Загружаем deployments в minikibe:

```sh
kubectl apply -f auth/deployments/blue.yaml
kubectl apply -f auth/deployments/green.yaml 
```

Убеждаемся, что поды создались и готовы:

```
NAME                          READY   STATUS    RESTARTS   AGE
auth-blue-56c4fbdfcb-pxmlh    1/1     Running   0          2m49s
auth-green-69c69556c4-czddr   1/1     Running   0          53s
```

Создаем сервис `NodePort`, который направляет наш траффик на "синий" сервис:

```sh
kubectl apply -f auth/deployments/service_node_port.yaml
```

Получаем ссылку для подключения к нашему сервису:

```sh
minikube service messenger -n messenger --url
```

Полученную ссылку открываем в браузере и видим надпись: "Auth service say: Hello!" синего цвета. Теперь мы видим что наш "зеленый" сервис готов и можем переключить траффик на него командой:

```sh

kubectl patch service messenger -n messenger -p '{"spec": {"selector": {"version": "green"}}}'
```

Обновляем браузер и видим, что теперь приветственная надпись зеленого цвета, наш входящий траффик идет на "зеленую" версию приложения.


# Реализация стратегии деплоя canary

Для реализации стратегии canary мы  увеличиваем количество реплик "синего" приложения  до 5

```sh
kubectl scale deployment auth-blue --replicas=5 -n messenger
```

```
NAME                          READY   STATUS    RESTARTS   AGE
auth-blue-56c4fbdfcb-hnb69    1/1     Running   0          2m36s
auth-blue-56c4fbdfcb-nd54f    1/1     Running   0          2m36s
auth-blue-56c4fbdfcb-pxmlh    1/1     Running   0          24m
auth-blue-56c4fbdfcb-tzg6s    1/1     Running   0          2m36s
auth-blue-56c4fbdfcb-xnqg5    1/1     Running   0          2m36s
auth-green-69c69556c4-czddr   1/1     Running   0          22m
```
И удаляем метку "version: green" у селектора сервиса:

```sh
kubectl patch service messenger -n messenger --type json -p='[{"op": "remove", "path": "/spec/selector/version"}]'
```

Теперь наш траффик распределяется между всеми шестью репликами сервиса. Можем в этом убедится, обновляя браузер (с очисткой кеша). Чаще будут выводится надписи синего цвета, но и иногда и зеленого. После того как мы убедились, что "зеленая" версия работоспособна, мы можем плавно увеличить количество "зеленых" реплик и уменьшить "синих". Таким образом плавно управляя переходом от одно версии сервиса к другой.
