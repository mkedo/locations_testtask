| Архитектура | Чтение rps (99%) | Запись rps (99%) |
|---|---|----|
|tarantool + vinyl | 3490 (50ms) | 1844 (84ms) | 
|postgres + redis (cache) | 1850 (59ms) | 658 (1185ms) |
|postgres (no cache) | 1493 (550ms) | 621 (1168ms) |
|postgres + redis (persistence) | 1037 (262ms) | 1904 (61ms) |

- Запись на **postgres** медленнее из-за бОльших накладных расходов: как минимум
проверка конститентности и парсинг SQL-запроса. 

- В связке **postgres + redis (cache)**, по сути, тестируется скорость чтения из **redis**,
т.к ```ab``` обращается только к одному объявлению.
Скорость чтения будет зависеть от попадания в кеш.

- Чтение из **postgres** медленнее, т.к **postgres** приходится пройтись по 2 индексам 
в item_locations и locations, и сделать JOIN.
Возможно денормализация поможет сделать выборку быстрее.

## tarantool + vinyl

Чтение

```
ab -k -c 100 -n 2000 127.0.0.1:8080/random_item/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /random_item/locations
Document Length:        390 bytes

Concurrency Level:      100
Time taken for tests:   0.573 seconds
Complete requests:      2000
Failed requests:        1998
   (Connect: 0, Receive: 0, Length: 1998, Exceptions: 0)
Keep-Alive requests:    1330
Total transferred:      3277935 bytes
HTML transferred:       3024055 bytes
Requests per second:    3490.20 [#/sec] (mean)
Time per request:       28.652 [ms] (mean)
Time per request:       0.287 [ms] (mean, across all concurrent requests)
Transfer rate:          5586.25 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.3      0       1
Processing:     4   27  10.6     27      54
Waiting:        4   22   7.5     21      46
Total:          5   27  10.6     27      54

Percentage of the requests served within a certain time (ms)
  50%     27
  66%     32
  75%     34
  80%     36
  90%     43
  95%     46
  98%     48
  99%     50
 100%     54 (longest request)
```

Запись

```
ab -T application/json -p post.req -c 150 -n 2000 127.0.0.1:8080/random_item/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /random_item/locations
Document Length:        9 bytes

Concurrency Level:      150
Time taken for tests:   1.084 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      250000 bytes
Total body sent:        396000
HTML transferred:       18000 bytes
Requests per second:    1844.91 [#/sec] (mean)
Time per request:       81.305 [ms] (mean)
Time per request:       0.542 [ms] (mean, across all concurrent requests)
Transfer rate:          225.21 [Kbytes/sec] received
                        356.73 kb/s sent
                        581.94 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       1
Processing:     3   78  12.2     81      85
Waiting:        1   43  22.9     42      85
Total:          3   78  12.2     81      85

Percentage of the requests served within a certain time (ms)
  50%     81
  66%     82
  75%     82
  80%     83
  90%     84
  95%     84
  98%     85
  99%     85
 100%     85 (longest request)
```


## postgres + redis (cache)

Чтение

```
ab -k -c 100 -n 2000 127.0.0.1:8080/item/5/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /item/5/locations
Document Length:        2573 bytes

Concurrency Level:      100
Time taken for tests:   1.081 seconds
Complete requests:      2000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      5340000 bytes
HTML transferred:       5146000 bytes
Requests per second:    1850.03 [#/sec] (mean)
Time per request:       54.053 [ms] (mean)
Time per request:       0.541 [ms] (mean, across all concurrent requests)
Transfer rate:          4823.81 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       2
Processing:     3   52   6.7     53      59
Waiting:        2   29  14.9     29      57
Total:          3   52   6.7     53      60

Percentage of the requests served within a certain time (ms)
  50%     53
  66%     54
  75%     54
  80%     54
  90%     56
  95%     57
  98%     59
  99%     59
 100%     60 (longest request)
```

Запись

```
ab -T application/json -p post.req -c 100 -n 2000 127.0.0.1:8080/random_item/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /random_item/locations
Document Length:        9 bytes

Concurrency Level:      100
Time taken for tests:   3.037 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      250000 bytes
Total body sent:        396000
HTML transferred:       18000 bytes
Requests per second:    658.51 [#/sec] (mean)
Time per request:       151.859 [ms] (mean)
Time per request:       1.519 [ms] (mean, across all concurrent requests)
Transfer rate:          80.38 [Kbytes/sec] received
                        127.33 kb/s sent
                        207.71 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       2
Processing:     5  147 239.9     73    2696
Waiting:        5  146 239.9     72    2696
Total:          6  147 239.9     73    2696

Percentage of the requests served within a certain time (ms)
  50%     73
  66%    122
  75%    166
  80%    203
  90%    304
  95%    478
  98%    962
  99%   1185
 100%   2696 (longest request)
```

postgres
--
Чтение без кеша

```
ab -k -c 100 -n 2000 127.0.0.1:8080/random_item/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /random_item/locations
Document Length:        2573 bytes

Concurrency Level:      100
Time taken for tests:   1.339 seconds
Complete requests:      2000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      5340000 bytes
HTML transferred:       5146000 bytes
Requests per second:    1493.57 [#/sec] (mean)
Time per request:       66.954 [ms] (mean)
Time per request:       0.670 [ms] (mean, across all concurrent requests)
Transfer rate:          3894.36 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       1
Processing:     5   64  90.8     36     782
Waiting:        3   60  91.2     31     782
Total:          6   64  90.8     36     783

Percentage of the requests served within a certain time (ms)
  50%     36
  66%     51
  75%     68
  80%     83
  90%    130
  95%    193
  98%    438
  99%    550
 100%    783 (longest request)
```

Запись без кеша

```
ab -T application/json -p post.req -c 100 -n 2000 127.0.0.1:8080/random_item/locations

Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /random_item/locations
Document Length:        9 bytes

Concurrency Level:      100
Time taken for tests:   3.216 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      250000 bytes
Total body sent:        396000
HTML transferred:       18000 bytes
Requests per second:    621.85 [#/sec] (mean)
Time per request:       160.809 [ms] (mean)
Time per request:       1.608 [ms] (mean, across all concurrent requests)
Transfer rate:          75.91 [Kbytes/sec] received
                        120.24 kb/s sent
                        196.15 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1  11.2      0     501
Processing:     5  131 215.6     59    2505
Waiting:        4  131 215.6     59    2505
Total:          5  132 215.8     60    2505

Percentage of the requests served within a certain time (ms)
  50%     60
  66%    106
  75%    146
  80%    179
  90%    290
  95%    491
  98%    923
  99%   1168
 100%   2505 (longest request)
```

postgres + redis (persistence)
--

Чтение

```
ab -k -c 100 -n 2000 127.0.0.1:8080/item/5/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /item/5/locations
Document Length:        2450 bytes

Concurrency Level:      100
Time taken for tests:   1.928 seconds
Complete requests:      2000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      5094000 bytes
HTML transferred:       4900000 bytes
Requests per second:    1037.22 [#/sec] (mean)
Time per request:       96.412 [ms] (mean)
Time per request:       0.964 [ms] (mean, across all concurrent requests)
Transfer rate:          2579.87 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      1      12
Processing:     6   93  52.1     79     371
Waiting:        5   87  52.5     74     369
Total:          6   93  52.2     80     372
ERROR: The median and mean for the initial connection time are more than twice the standard
       deviation apart. These results are NOT reliable.

Percentage of the requests served within a certain time (ms)
  50%     80
  66%    109
  75%    126
  80%    135
  90%    160
  95%    183
  98%    229
  99%    262
 100%    372 (longest request)

```

Запись

```
ab -T application/json -p post.req -c 100 -n 2000 127.0.0.1:8080/item/5/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /item/5/locations
Document Length:        9 bytes

Concurrency Level:      100
Time taken for tests:   1.050 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      250000 bytes
Total body sent:        386000
HTML transferred:       18000 bytes
Requests per second:    1904.59 [#/sec] (mean)
Time per request:       52.505 [ms] (mean)
Time per request:       0.525 [ms] (mean, across all concurrent requests)
Transfer rate:          232.49 [Kbytes/sec] received
                        358.97 kb/s sent
                        591.47 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      0       1
Processing:     2   51   7.4     52      72
Waiting:        1   33  13.1     36      58
Total:          2   51   7.3     52      72

Percentage of the requests served within a certain time (ms)
  50%     52
  66%     53
  75%     54
  80%     55
  90%     57
  95%     57
  98%     59
  99%     61
 100%     72 (longest request)
```


Используемые ключи
---
```
-k              Use HTTP KeepAlive feature
-n requests     Number of requests to perform
-c concurrency  Number of multiple requests to make at a time
-p postfile     File containing data to POST. Remember also to set -T
-T content-type Content-type header to use for POST/PUT data, eg.
```