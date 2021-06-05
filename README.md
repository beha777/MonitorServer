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
