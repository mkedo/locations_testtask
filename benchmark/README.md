| Архитектура | Чтение 100%,мс | Запись 100%,мс |
|---|---|----|
|tarantool + vinyl | 57 | 56 | 
|postgres + redis (cache) | 60 | 933 |
|postgres (no cache) | 810 | 760 |
|postgres + redis (persistence) | 372 | 72 |

- Запись на **postgres** медленнее из-за бОльших накладных расходов: как минимум
проверка конститентности и парсинг SQL-запроса. 
К тому же запись происходит в однопоточном режиме,
 чтобы избежать чрезмерного количества ошибок сериализации 
(```ab``` делает множество конкурентных запросов к одному объявлению. 
По идее, одно объявление так часто не обновляют и такое должно случатся намного реже.
 Есть место для оптимизации).


Запись в многопоточном режиме (со случайными id объявлений)
```
Concurrency Level:      100
Time taken for tests:   4.212 seconds
Complete requests:      3000
Failed requests:        0
Total transferred:      375000 bytes
Total body sent:        579000
HTML transferred:       27000 bytes
Requests per second:    712.21 [#/sec] (mean)
Time per request:       140.408 [ms] (mean)
Time per request:       1.404 [ms] (mean, across all concurrent requests)
Transfer rate:          86.94 [Kbytes/sec] received
                        134.23 kb/s sent
                        221.17 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       4
Processing:     5  137 338.8     45    3488
Waiting:        5  137 338.8     44    3488
Total:          5  137 338.8     45    3488

Percentage of the requests served within a certain time (ms)
  50%     45
  66%     81
  75%    116
  80%    142
  90%    264
  95%    450
  98%   1257
  99%   2103
 100%   3488 (longest request)
```


Запись в однопоточном (со случайными id объявлений)
```
Concurrency Level:      100
Time taken for tests:   10.093 seconds
Complete requests:      3000
Failed requests:        0
Total transferred:      375000 bytes
Total body sent:        579000
HTML transferred:       27000 bytes
Requests per second:    297.25 [#/sec] (mean)
Time per request:       336.419 [ms] (mean)
Time per request:       3.364 [ms] (mean, across all concurrent requests)
Transfer rate:          36.29 [Kbytes/sec] received
                        56.02 kb/s sent
                        92.31 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       3
Processing:    10  328  74.9    286     457
Waiting:       10  328  75.0    285     457
Total:         10  329  75.0    286     457

Percentage of the requests served within a certain time (ms)
  50%    286
  66%    394
  75%    415
  80%    421
  90%    428
  95%    432
  98%    438
  99%    443
 100%    457 (longest request)
```


- В связке **postgres + redis (cache)**, по сути, тестируется скорость чтения из **redis**,
т.к ```ab``` обращается только к одному объявлению.
Скорость чтения будет зависеть от попадания в кеш.

- Чтение из **postgres** медленнее, т.к **postgres** приходится пройтись по 2 индексам 
в item_locations и locations, и, возможно, считывать данные с диска.

## tarantool + vinyl

Чтение

```
ab -k -c 100 -n 2000 127.0.0.1:8080/item/5/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /item/5/locations
Document Length:        2448 bytes

Concurrency Level:      100
Time taken for tests:   1.057 seconds
Complete requests:      2000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      5090000 bytes
HTML transferred:       4896000 bytes
Requests per second:    1892.04 [#/sec] (mean)
Time per request:       52.853 [ms] (mean)
Time per request:       0.529 [ms] (mean, across all concurrent requests)
Transfer rate:          4702.38 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       1
Processing:     9   51   6.6     52      57
Waiting:        1   30  14.3     30      56
Total:          9   51   6.6     52      57

Percentage of the requests served within a certain time (ms)
  50%     52
  66%     53
  75%     53
  80%     53
  90%     54
  95%     56
  98%     56
  99%     57
 100%     57 (longest request)
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
Time taken for tests:   1.058 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      250000 bytes
Total body sent:        416000
HTML transferred:       18000 bytes
Requests per second:    1890.25 [#/sec] (mean)
Time per request:       52.903 [ms] (mean)
Time per request:       0.529 [ms] (mean, across all concurrent requests)
Transfer rate:          230.74 [Kbytes/sec] received
                        383.96 kb/s sent
                        614.70 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       1
Processing:     2   51   6.6     52      56
Waiting:        1   29  14.7     29      56
Total:          2   52   6.6     52      56

Percentage of the requests served within a certain time (ms)
  50%     52
  66%     53
  75%     54
  80%     54
  90%     55
  95%     55
  98%     56
  99%     56
 100%     56 (longest request)
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
ab -T application/json -p post.req -c 100 -n 2000 127.0.0.1:8080/item/5/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /item/5/locations
Document Length:        9 bytes

Concurrency Level:      200
Time taken for tests:   8.328 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      250000 bytes
Total body sent:        386000
HTML transferred:       18000 bytes
Requests per second:    240.14 [#/sec] (mean)
Time per request:       832.848 [ms] (mean)
Time per request:       4.164 [ms] (mean, across all concurrent requests)
Transfer rate:          29.31 [Kbytes/sec] received
                        45.26 kb/s sent
                        74.57 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       1
Processing:    67  786 137.7    824     933
Waiting:       61  786 137.8    824     932
Total:         67  787 137.7    825     933

Percentage of the requests served within a certain time (ms)
  50%    825
  66%    832
  75%    837
  80%    839
  90%    849
  95%    900
  98%    926
  99%    930
 100%    933 (longest request)
```

postgres
--
Чтение без кеша

```
ab -k -c 100 -n 2000 127.0.0.1:8080/item/5/locations
...
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /item/5/locations
Document Length:        2573 bytes

Concurrency Level:      100
Time taken for tests:   1.315 seconds
Complete requests:      2000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      5340000 bytes
HTML transferred:       5146000 bytes
Requests per second:    1520.82 [#/sec] (mean)
Time per request:       65.754 [ms] (mean)
Time per request:       0.658 [ms] (mean, across all concurrent requests)
Transfer rate:          3965.42 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       1
Processing:     4   63  74.0     51     809
Waiting:        3   56  75.5     34     809
Total:          4   63  74.0     51     810

Percentage of the requests served within a certain time (ms)
  50%     51
  66%     58
  75%     63
  80%     73
  90%    120
  95%    187
  98%    300
  99%    428
 100%    810 (longest request)
```

Запись без кеша

```
ab -T application/json -p post.req -c 100 -n 2000 127.0.0.1:8080/item/5/locations
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /item/5/locations
Document Length:        9 bytes

Concurrency Level:      200
Time taken for tests:   7.192 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      250000 bytes
Total body sent:        386000
HTML transferred:       18000 bytes
Requests per second:    278.07 [#/sec] (mean)
Time per request:       719.241 [ms] (mean)
Time per request:       3.596 [ms] (mean, across all concurrent requests)
Transfer rate:          33.94 [Kbytes/sec] received
                        52.41 kb/s sent
                        86.35 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       1
Processing:    94  678 145.7    724     760
Waiting:       81  677 146.1    724     759
Total:         95  678 145.7    724     760

Percentage of the requests served within a certain time (ms)
  50%    724
  66%    730
  75%    736
  80%    737
  90%    742
  95%    748
  98%    753
  99%    754
 100%    760 (longest request)
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