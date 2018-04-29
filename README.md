ramss
===============
Rest API manipulates systemd service


setup
--------------

```
$ mkdir -p $GOPATH/src/github.com/teddyyy
$ cd $GOPATH/src/github.com/teddyyy
$ git clone https://github.com/teddyyy/ramss.git
$ cd ramss
$ dep ensure
$ go build
```

usage
--------------
```
$ ./ramss -h
Usage of ./ramss:
  -f string
        config file (default "./config.yaml")
  -p string
        listen port (default "8080")
```

run
--------------
```
sudo ./ramss -f sample.yaml -p 10080
```

API
--------------
* GET `/api/v1/systemd/`
  * The status of all systemd services is obtained by GET.
  * Ex) Response body
  ```
  {
    "Systemd": [
      {
        "Id": "nginx",
        "Description": "A high performance web server and a reverse proxy server",
        "LoadState": "loaded",
        "ActiveState": "active",
        "UnitFileState": "enabled",
        "MainPID": 12577
      },
      {
        "Id": "rsyslog",
        "Description": "System Logging Service",
        "LoadState": "loaded",
        "ActiveState": "active",
        "UnitFileState": "enabled",
        "MainPID": 11656
      }
    ]
  }
  ```

* GET `/api/v1/systemd/<service>`
  * The status of systemd service is obtained by GET.
  * The `<service>` is defined in config file and match with unit_name.
  * Ex) Response body
  ```
  {
    "Id": "nginx",
    "Description": "A high performance web server and a reverse proxy server",
    "LoadState": "loaded",
    "ActiveState": "active",
    "UnitFileState": "enabled",
    "MainPID": 12577
  }
  ```

* POST `/api/v1/systemd/<service>`
  * The status of systemd service is manipulated by POST.
  * The `<service>` is defined in config file and match with unit_name.
  * Ex) Request body
  ```
  {
    "action":"start",
    "mode":"replace"
  }
  ```
  * Ex) Response body
  ```
  {
    "Code": 200,
    "Message": "success"
  }
  ```

  * The posible `"action"` are.
    * start
    * stop
    * restart
    * reload
  * The posible `"mode"` are.
    * replace
    * replace-irreversibly
    * isolate
    * flush
