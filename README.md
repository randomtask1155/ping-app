# ping-app


Simple application that will ping an app using its normal domain name as well as ping the app directly through the gorouter.   Useful for testing common network issues between an app and load balancer.


## example manifest

```
---
applications:
  - name: ping-app
    memory: 128M
    instances: 1
    buildpack: go_buildpack
    command: ping-app
    env:
      INTERVAL: 1m
      REMOTE_APP_HOSTNAME: ping-app-b.mydomain.com
      GOROUTER_ADDRESS: "10.10.10.10:10.10.10.11"
```

## push app A

edit manifest variables

* REMOTE_APP_HOSTNAME - should be the hostname of app b
* GOROUTER_ADDRESS    - should be the ip address of one or more IP address delminted by ":"


```
cf push ping-app-a
```


## push app B

edit manifest variables

* REMOTE_APP_HOSTNAME - should be the hostname of app a
* GOROUTER_ADDRESS    - should be the ip address of one or more IP address delminted by ":"


```
cf push ping-app-b
```



## example logs from sending app

```
   2022-05-02T12:20:34.57-0500 [APP/PROC/WEB/0] OUT Sending POST to ping-app-b.mydomain.com
   2022-05-02T12:20:34.59-0500 [APP/PROC/WEB/0] OUT HTTP/1.1 200 OK
   2022-05-02T12:20:34.59-0500 [APP/PROC/WEB/0] OUT Connection: close
   2022-05-02T12:20:34.59-0500 [APP/PROC/WEB/0] OUT Content-Length: 19
   2022-05-02T12:20:34.59-0500 [APP/PROC/WEB/0] OUT Content-Type: text/html; charset=utf-8
   2022-05-02T12:20:34.59-0500 [APP/PROC/WEB/0] OUT Date: Mon, 02 May 2022 17:20:34 GMT
   2022-05-02T12:20:34.59-0500 [APP/PROC/WEB/0] OUT X-Vcap-Request-Id: b411d8ea-817f-4ee7-6375-ed9aff10d26e
   2022-05-02T12:20:34.59-0500 [APP/PROC/WEB/0] OUT <html>Hello!</html>
   2022-05-02T12:20:34.59-0500 [APP/PROC/WEB/0] OUT Sending POST to 10.10.10.10
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT HTTP/1.1 200 OK
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT Connection: close
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT Content-Length: 19
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT Content-Type: text/html; charset=utf-8
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT Date: Mon, 02 May 2022 17:20:34 GMT
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT X-Vcap-Request-Id: 48fa551c-f528-4878-4142-1a1641db1f57
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT <html>Hello!</html>
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT Sending POST to 10.10.10.11
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT Failed sending post to gorouter: 10.225.28.140: An Error Occured sending POST Post "https://10.10.10.11": dial tcp 10.10.10.11:443: connect: connection refused
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT domain stats: error=0 success=7
   2022-05-02T12:20:34.61-0500 [APP/PROC/WEB/0] OUT gorouter stats: error=7 success=7
   2022-05-02T12:20:37.43-0500 [APP/PROC/WEB/0] ERR 2022/05/02 17:20:37 "{\"data\": \"TCoaNatyy\"}"
```


## example logs from receiving app

```
2022-04-27T15:07:43.96-0500 [APP/PROC/WEB/0] ERR 2022/04/27 20:07:43 "{\"data\": \"CMRAjWwhT\"}"
```


## update variables dynamically

* interval - sets the interval using format defined in https://pkg.go.dev/time#ParseDuration
* remoteapp - change which app to pinng to
* postbodysize - change the amount of asci characters sent to remote app in each request
* gorouter - change the ip of gorouter app will send posts to


```
curl 'ping-app-a.mydomain.com/config?interval=1m&remoteapp=ping-app-b.mydomain.com&postbodysize=100&gorouter=10.10.10.11'
```


## reset stats counters

after each ping iteration the app will logs stats for each flow

```
   2022-04-28T12:30:54.46-0500 [APP/PROC/WEB/0] OUT domain stats: error=0 success=7
   2022-04-28T12:30:54.46-0500 [APP/PROC/WEB/0] OUT gorouter stats: error=6 success=1
```

these stats can be reset using the /resetcounters endpoint

```
curl ping-app-a.mydomain.com/resetcounters
```