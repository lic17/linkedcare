#Step 1:
```
istioctl operator init
istioctl manifest apply --set addonComponents.grafana.enabled=true --set addonComponents.kiali.enabled=true --set addonComponents.tracing.enabled=true --set values.kiali.hub=registry.cn-hangzhou.aliyuncs.com/linkedcare
```
#Step2
```
kubectl apply -f kube-app-manager-controller.yaml
```
#Step3
```
创建服务治理和灰度发布的CRD
kubectl apply -f  servicemesh_v1alpha2_servicepolicy.yaml
kubectl apply -f  servicemesh_v1alpha2_strategy.yaml
```
#Step4
```
创建应用管理后端和服务治理相关operator
kubectl apply -f linkedcare-system.yaml
```
