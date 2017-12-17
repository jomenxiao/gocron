## gocron
**this repo replace crontab, and offer http api**

## usages
```
Usage of ./gocron:
  -port int
    	listen port (default 8888)
```


example
```
curl -X POST localhost:8888 -d '{"id": "date_command", "cmd": "date", "args": ["-R"], "interval": 5000}'
```

```
curl -X DELETE localhost:8888 -d '{"id": "date_command"}'
```

