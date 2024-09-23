# vpnSwitcher
Программа автоматически поддерживает интернет-соединение на сервере, отдавая преимущество нахождения в VPN сети.

## Запуск
`go run main`

или компилируем в статический бинарный файл:

`go run build -a -ldflags '-extldflags "-static"' -o vpnSwitcher`
