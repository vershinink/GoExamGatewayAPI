### GoExamGatewayAPI

Сервис шлюз проксирует запросы к различным сервисам системы. Практика на курсе "Go-разработчик" от SkillFactory. Часть итогового проекта курса.

Для запуска нужно установить путь к файлу конфига в переменную окружения `GATEWAY_CONFIG_PATH`. Адреса сервисов и другие входные данные указываются в файле конфига.

Сам файл конфига `config.yaml` лежит в каталоге config.

**Сделано:**

- Логирование в stdout через пакет slog стандартной библиотеки Go.
- 3 REST API метода: на получение списка новостных статей с пагинацией из сервиса агрегатора; на получение статьи по ее ID с полным деревом комментариев из агрегатора и сервиса комментариев; добавление нового комментария в сервис комментариев через проверку в сервисе цензурирования.
- Тесты для всех основных пакетов приложения.
- Использование контекстов при работе сервера.
- Использоваие middleware для трассировки запросов и логирования.
- Завершение работы приложения по сигналу прерывания с использованием graceful shutdown.
- Сборка и запуск сервиса в Docker контейнере.