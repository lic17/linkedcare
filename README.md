# go build

```
./build.run
```

# dockere build

## controller-manager：

```
  sudo  docker build -t  registry.cn-hangzhou.aliyuncs.com/linkedcare/lc-controller-manager:0.1-v1alpha1 build/controller-manager/
```

## Apiserver：
```
  sudo  docker build -t  registry.cn-hangzhou.aliyuncs.com/linkedcare/lc-apiserver:0.1-v1alpha1 build/apiserver
```

# run apiserver

```
./apiserver --v=8 --kubeconfig=/home/licheng/.kube/config
```

## cert
### list:

```
curl http://127.0.0.1:9090/api/cert.linkedcare.io/v1beta1/certs
```

```
[
 {
  "Name": "/etc/kubernetes/pki/apiserver-etcd-client.crt",
  "Time": "358"
 },
 {
  "Name": "/etc/kubernetes/pki/apiserver-kubelet-client.crt",
  "Time": "358"
 },
 {
  "Name": "/etc/kubernetes/pki/apiserver.crt",
  "Time": "358"
 },
 {
  "Name": "/etc/kubernetes/pki/ca.crt",
  "Time": "3012"
 },
 {
  "Name": "/etc/kubernetes/pki/etcd/ca.crt",
  "Time": "3012"
 },
 {
  "Name": "/etc/kubernetes/pki/etcd/healthcheck-client.crt",
  "Time": "358"
 },
 {
  "Name": "/etc/kubernetes/pki/etcd/peer.crt",
  "Time": "358"
 },
 {
  "Name": "/etc/kubernetes/pki/etcd/server.crt",
  "Time": "358"
 },
 {
  "Name": "/etc/kubernetes/pki/expired/apiserver.crt",
  "Time": "92"
 },
 {
  "Name": "/etc/kubernetes/pki/front-proxy-ca.crt",
  "Time": "3012"
 },
 {
  "Name": "/etc/kubernetes/pki/front-proxy-client.crt",
  "Time": "358"
 }
]
```

## application

### list:

```
  http://127.0.0.1:9090/api/app.k8s.io/v1beta1/namespace/default/application
```

### get:
```
  http://127.0.0.1:9090/api/app.k8s.io/v1beta1/namespace/default/name/bookinfo/application
```

### create
```
curl -H "Content-Type:application/json" -H "Data_Type:msg" -X POST --data '{"name": "lctest", "version": "v2"}' http://127.0.0.1:9090/api/app.k8s.io/v1beta1/namespace/default/application
```
### delete
```
curl -v -X DELETE http://127.0.0.1:9090/api/app.k8s.io/v1beta1/namespace/default/name/lctest/application
```

### restart app
```
curl -H "Content-Type:application/json" -H "Data_Type:msg" -X PUT http://127.0.0.1:9090/api/app.k8s.io/v1beta1/namespace/default/name/bookinfo/application #重启application删除整个app的pod来重启app
```

### restart service
```
curl -H "Content-Type:application/json" -H "Data_Type:msg" -X PUT http://127.0.0.1:9090/api/app.k8s.io/v1beta1/namespace/default/name/ratings/service #重启svc删除这个svc的pod来重启svc
```

## strategy

### list:

```
  http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/strategy
```
### get:

```
  http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/name/review-v2/strategy
```
### create:
```
curl -H "Content-Type:application/json" -H "Data_Type:msg" -X POST --data '{"component":"reviews","governor": "","principal": "v1","application":{"name": "bookinfo","version": "v1"},"template":{"spec": {"hosts": ["reviews"],"http":[{"route": [{"destination": {"host": "reviews","subset": "v1"},"weight": 48},{"destination": {"host": "reviews","subset": "v2"},"weight": 52}]}]}}}' http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/strategy
```

### update:
```
curl -H "Content-Type:application/json" -H "Data_Type:msg" -X PUT --data '{"component":"reviews","governor": "","principal": "v1","application":{"name": "bookinfo","version": "v1"},"template":{"spec": {"hosts": ["reviews"],"http":[{"route": [{"destination": {"host": "reviews","subset": "v1"},"weight": 48},{"destination": {"host": "reviews","subset": "v2"},"weight": 52}]}]}}}' http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/name/reviews/strategy
```

### delete
```
curl -v -X DELETE http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/name/review-v2/strategy
```

## servicepolicy

### list:

```
  http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/servicepolicy
```

### get:

```
  http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/name/review-v2/servicepolicy

```

### create:
```
curl -H "Content-Type:application/json" -H "Data_Type:msg" -X POST --data '{"component":"reviews","application":{"name": "bookinfo","version": "v1"},"template":{"spec":{"host":"reviews","trafficPolicy":{"loadBalancer":{"simple":"LEAST_CONN"}}}}}' http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/servicepolicy
```

### update:
```
curl -H "Content-Type:application/json" -H "Data_Type:msg" -X PUT --data '{"component":"reviews","application":{"name": "bookinfo","version": "v1"},"template":{"spec":{"host":"reviews","trafficPolicy":{"loadBalancer":{"simple":"LEAST_CONN"}}}}}' http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/name/reviews/servicepolicy
```

```
{"component":"reviews","application":{"name": "bookinfo","version": "v1"},"template":{"spec":{"host":"reviews","trafficPolicy":{"loadBalancer":{"simple":"LEAST_CONN"}}}}}
host值和component值相同
```

### delete
```
curl -v -X DELETE http://127.0.0.1:9090/api/servicemesh.linkedcare.io/v1alpha1/namespace/default/name/review-v2/servicepolicy

```

```
{
   "component": "review",
   "governor": "",
   "principal": "v1",
   "application": {
     "name": "bookinfo",
     "version": "v1",
   },"template": {
      "spec": {
          "hosts": [
              "reviews"
          ],
          "http": [
              {
                  "route": [
                      {
                          "destination": {
                              "host": "reviews",
                              "subset": "v1"
                          },
                          "weight": 50
                      },
                      {
                          "destination": {
                              "host": "reviews",
                              "subset": "v2"
                          },
                          "weight": 50
                      }
                  ]
              }
          ]
      }
   },
}
```

```
apiVersion: servicemesh.linkedcare.io/v1alpha1
kind: Strategy
metadata:
  annotations:
    app.kubernetes.io/icon: /assets/bookinfo.svg
    servicemesh.kubesphere.io/workloadReplicas: "1"
    servicemesh.kubesphere.io/workloadType: deployments
  generation: 3
  labels:
    app: details
    app.kubernetes.io/name: bookinfo
    app.kubernetes.io/version: v1
  name: test
  namespace: test
  resourceVersion: "234724779"
  selfLink: /apis/servicemesh.kubesphere.io/v1alpha2/namespaces/test/strategies/test
  uid: 33a887b9-91ce-11ea-9097-00163f00bfd3
spec:
  governor: v1 #v1接管所有流量 只有流量全部接管灰度发布才可以下线
  hosts: details
  principal: v1
  protocol: http
  selector:
    matchLabels:
      app: details
      app.kubernetes.io/name: bookinfo
      app.kubernetes.io/version: v1
  strategyPolicy: WaitForWorkloadReady
  template:
    spec:
      hosts:
      - details
      http:
      - match:
        - headers:
            User-Agent:
              regex: .*(Linux ).*
            cookie:
              regex: q=c
            test:
              exact: "1"
          uri:
            prefix: test
        route:
        - destination:
            host: details
            subset: v2
          weight: 100
      - route:
        - destination:
            host: details
            subset: v1
          weight: 100
  type: Canary
```


```
apiVersion: servicemesh.linkedcare.io/v1alpha1
kind: ServicePolicy
metadata:
  labels:
    app: reviews
    app.kubernetes.io/name: bookinfo
    app.kubernetes.io/version: v1
  name: reviews
  namespace: test
  ownerReferences:
  - apiVersion: app.k8s.io/v1beta1
    blockOwnerDeletion: true
    controller: false
    kind: Application
    name: bookinfo
    uid: b571fbd9-7fb9-11ea-9097-00163f00bfd3
  resourceVersion: "239006786"
  selfLink: /apis/servicemesh.kubesphere.io/v1alpha2/namespaces/test/servicepolicies/reviews
  uid: b0868f1b-95c4-11ea-b6af-325dbbd1d7e6
spec:
  selector:
    matchForLables:
      app: reviews
      app.kubernetes.io/name: bookinfo
      app.kubernetes.io/version: v1
  template:
    labels:
      app: reviews
      app.kubernetes.io/name: bookinfo
      app.kubernetes.io/version: v1
//关注该spec下的内容
    spec:
      host: reviews
      trafficPolicy:
//连接池管理
        connectionPool:
          http:
            http1MaxPendingRequests: 1024
            http2MaxRequests: 1024
            maxRequestsPerConnection: 1
            maxRetries: 3
          tcp:
            connectTimeout: 100ms
            maxConnections: 100
//负载均衡和会话保持算法选择
        loadBalancer:
          simple: ROUND_ROBIN
          #consistentHash:
          #  useSourceIp: true

//熔断器
        outlierDetection:
          baseEjectionTime: 30s
          consecutiveErrors: 5
          interval: 10s
          maxEjectionPercent: 10

```



# TODO
```
1. copy application from one namespace to the other namespace or cluster! 80%

2. move ingress deployment service form helm to app! 0%
```


# 依赖
```
../../github.com/lic17/application

来自github仓库github.com/lic17/application 的linkedcare分支
```

```
k8s.io/client-go

已上传http://git.lc.com/cheng.li/client-go.git
```

```
k8s.io/apimachinery

已上传http://git.lc.com/cheng.li/apimachinery.git
```
