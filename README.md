# MonitorServer
Create "settings-dev.json" with the following format

    {
        "app": {
            "serverName": "test",
            "portRun": "8001"
        },
        
        "postgresParams": {
            "server": "127.0.0.1",
            "port": "5432",
            "user": "postgres",
            "password": "admin",
            "database": "beha"
        },
        "botParams": {
            "token": "1857766717:AAGOdDVdYgYbj9yFBa5imAc9sUZR1Y7ZfL8",
            "chat_id": "@ServerParamStatus"
        }
        "periods": {
            "default_notification": 3600,
            "default_ticker": 60,
            "default_check": 30
        }
    }


ROUTES

router.POST("/addserver", addServer)

Request_body:

    {
        "host":     "127.0.0.1:2281",
        "login":    "root",
        "password": "admin",
        "os"":       "CentOS",
        "version":  "7.0",
    }

Данная программа позволяет добавлять сервера для проверки их состояния, а также проверки активности сервисов запущенных на всех добавленных серверах.
Таким образом в телеграм бот "@ServerParamStatus" будут поступать следующие уведомления:
    - база данных недоступна;
    - порт недоступен;
    - ip не пингуется;
    - сервис неактивен;
    - превышен допустимый лимит загруженности процессора;
    - превышен допустимый лимит загруженности ОЗУ;
    - превышен допустимый лимит загруженности жесткого диска;