tarantool + vinyl
--
Чтение

```
ab -k -c 100 -n 2000 127.0.0.1:8080/item/5/locations

This is ApacheBench, Version 2.3 <$Revision: 1638069 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


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

This is ApacheBench, Version 2.3 <$Revision: 1638069 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


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


postgres + redis
--

Чтение

```
ab -k -c 100 -n 2000 127.0.0.1:8080/getItemLocations?ItemId=5


This is ApacheBench, Version 2.3 <$Revision: 1638069 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /getItemLocations?ItemId=5
Document Length:        2447 bytes

Concurrency Level:      100
Time taken for tests:   39.863 seconds
Complete requests:      2000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      5088000 bytes
HTML transferred:       4894000 bytes
Requests per second:    50.17 [#/sec] (mean)
Time per request:       1993.128 [ms] (mean)
Time per request:       19.931 [ms] (mean, across all concurrent requests)
Transfer rate:          124.65 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      1       7
Processing:     3 1887 1995.0    813    6697
Waiting:        3 1846 1976.6    688    6681
Total:          4 1887 1995.0    814    6697
ERROR: The median and mean for the initial connection time are more than twice the standard
       deviation apart. These results are NOT reliable.

Percentage of the requests served within a certain time (ms)
  50%    814
  66%   2955
  75%   3944
  80%   4366
  90%   4929
  95%   5184
  98%   5470
  99%   5602
 100%   6697 (longest request)
```

Запись

```
ab -T application/json -p post.req -c 100 -n 2000 127.0.0.1:8080/item/5/locations

This is ApacheBench, Version 2.3 <$Revision: 1638069 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /item/5/locations
Document Length:        9 bytes

Concurrency Level:      100
Time taken for tests:   65.255 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      250000 bytes
Total body sent:        416000
HTML transferred:       18000 bytes
Requests per second:    30.65 [#/sec] (mean)
Time per request:       3262.763 [ms] (mean)
Time per request:       32.628 [ms] (mean, across all concurrent requests)
Transfer rate:          3.74 [Kbytes/sec] received
                        6.23 kb/s sent
                        9.97 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.6      1      21
Processing:     2 3106 5070.4     28   26980
Waiting:        2 3103 5068.4     26   26979
Total:          2 3107 5070.4     29   26981

Percentage of the requests served within a certain time (ms)
  50%     29
  66%    163
  75%   7530
  80%   8211
  90%   9448
  95%  14885
  98%  18015
  99%  19027
 100%  26981 (longest request)
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
