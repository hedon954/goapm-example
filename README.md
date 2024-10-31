# goapm-example

goapm-example is a simple example for goapm.

## Quick Start

Start the base infrastructures and services:

```bash
make docker-up
```

Send some requests to the services:

```bash
# successful request
curl http://127.0.0.1:30001/order/add?uid=1&sku_id=3&num=1

# failed request -> sku not enough
curl http://127.0.0.1:30001/order/add?uid=1&sku_id=3&num=1000000

# failed request -> user not exist
curl http://127.0.0.1:30001/order/add?uid=1000000&sku_id=3&num=1
```

## Some Links

- [Jaeger](http://127.0.0.1:16686)
- [Grafana](http://127.0.0.1:3000)
- [Prometheus](http://127.0.0.1:9090)
- [Kibana](http://127.0.0.1:5601)

## Effect Pictures

### Grafana
![Grafana-MySQL](https://hedonspace.oss-cn-beijing.aliyuncs.com/img/image-20241031185641444.png)

![Grafana-Biz](https://hedonspace.oss-cn-beijing.aliyuncs.com/img/image-20241031185718137.png)

![Application-Trafiic](https://hedonspace.oss-cn-beijing.aliyuncs.com/img/image-20241031185748259.png)

![Application-Runtime](https://hedonspace.oss-cn-beijing.aliyuncs.com/img/image-20241031185807108.png)

![Application-GOX](https://hedonspace.oss-cn-beijing.aliyuncs.com/img/image-20241031185826969.png)

### Jaeger
![Jaeger-Trace](https://hedonspace.oss-cn-beijing.aliyuncs.com/img/image-20241031190013255.png)


### Prometheus
![Prometheus](https://hedonspace.oss-cn-beijing.aliyuncs.com/img/image-20241031190039185.png)
