---
id: e7ek65s31xm3p1zf3346p2y
title: Self Hands on Course
desc: ''
updated: 1754317844897
created: 1741221353071
---

[!NOTE]

> Prometheus is an open-source monitoring and alerting system developed by SoundCloud and now maintained by the Cloud Native Computing Foundation (CNCF). It is widely used for monitoring cloud-native applications due to its powerful querying capabilities, dimensional data model, and built-in alerting.

## 1. First Steps

### Install Prometheus

```shell
cd pca/practice
docker run -d --name prometheus \
  -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

