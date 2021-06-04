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
