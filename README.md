# MonitorServer
Create "settings-dev.json" with the following format 
{
  "app" : {
    "serverName": "test",
    "portRun": "8001"
  },

  "postgresParams": {
    "server": "127.0.0.1",
    "port": "",
    "user": "postgres",
    "password": "",
    "database": ""
  }
}

ROUTES

router.POST("/addserver", addServer)
body:
{
    Host:     "127.0.0.1:2281",
    Login:    "root",
    Password: "password",
    OS:       "CentOS",
    Version:  "7.0",
}