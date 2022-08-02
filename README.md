go rum 可选参数
```shell
  -c int
        并发数量 (default 1)
  -n int
        图片数量 (default 1)
  -s int
        清理生成出来的水印 (default 0)
```

```shell
go run main.go -c 10 -n 100 

NumCPU: 8
水印图：100 张
同步执行时间 = 10.484067834s
无并发数量控制时间 = 2.310320875s
同时并发 10 个时间 = 2.140755708s
```

```php
php main.php

总耗时: 5.07026s
总耗时: 5070.26ms
水印decode耗时: 0.43ms
原图decode耗时: 22.75ms
新图encode耗时: 20.47ms
```
