# ping-app


Simple application that will ping an app using its normal domain name as well as ping the app directly through the gorouter.   Useful for testing common network issues between an app and load balancer.

## push app A

edit manifest variables

* REMOTE_APP_HOSTNAME - should be the hostname of app b
* GOROUTER_ADDRESS    - should be the ip address of one gorouter


```
cf push ping-app-a
```


## push app B

edit manifest variables

* REMOTE_APP_HOSTNAME - should be the hostname of app a
* GOROUTER_ADDRESS    - should be the ip address of one gorouter


```
cf push ping-app-b
```



## example logs

```
   2022-04-27T14:33:14.80-0500 [APP/PROC/WEB/0] OUT Sending POST to ping-app-b.mydomain.com
   2022-04-27T14:33:14.81-0500 [APP/PROC/WEB/0] OUT HTTP/1.1 200 OK
   2022-04-27T14:33:14.81-0500 [APP/PROC/WEB/0] OUT Content-Length: 19
   2022-04-27T14:33:14.81-0500 [APP/PROC/WEB/0] OUT Content-Type: text/html; charset=utf-8
   2022-04-27T14:33:14.81-0500 [APP/PROC/WEB/0] OUT Date: Wed, 27 Apr 2022 19:33:14 GMT
   2022-04-27T14:33:14.81-0500 [APP/PROC/WEB/0] OUT X-Vcap-Request-Id: a2ca9034-d979-46cb-7b79-830d31b472db
   2022-04-27T14:33:14.81-0500 [APP/PROC/WEB/0] OUT <html>Hello!</html>
   2022-04-27T14:33:14.81-0500 [APP/PROC/WEB/0] OUT Sending POST to 10.213.60.19
   2022-04-27T14:33:14.82-0500 [APP/PROC/WEB/0] OUT HTTP/1.1 200 OK
   2022-04-27T14:33:14.82-0500 [APP/PROC/WEB/0] OUT Content-Length: 19
   2022-04-27T14:33:14.82-0500 [APP/PROC/WEB/0] OUT Content-Type: text/html; charset=utf-8
   2022-04-27T14:33:14.82-0500 [APP/PROC/WEB/0] OUT Date: Wed, 27 Apr 2022 19:33:14 GMT
   2022-04-27T14:33:14.82-0500 [APP/PROC/WEB/0] OUT X-Vcap-Request-Id: 2788c994-0cee-41d9-4a6c-041614ac3566
   2022-04-27T14:33:14.82-0500 [APP/PROC/WEB/0] OUT <html>Hello!</html>
   2022-04-27T14:33:33.03-0500 [APP/PROC/WEB/0] ERR 2022/04/27 19:33:33 {"email":"Toby@example.com","name":"Toby"}
```