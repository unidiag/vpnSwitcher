# vpnSwitcher
Программа автоматически поддерживает активное интернет-соединение на сервере, отдавая преимущество нахождение в VPN сети. Таким образом, обеспечивается доступ к данному серверу из интернета, если у него нет статического IP-адреса.

## Компиляция

`go build -a -ldflags '-extldflags "-static"' -o vpnSwitcher`

## Запуск
`./vpnSwitcher`

По умолчанию пингуемый домен VPN-сервера берётся из файла `/etc/openvpn/client.conf`. Если у вас другой файл, то можно его указать во втором параметре. Например:

`./vpnSwitcher /etc/openvpn/another.conf`
