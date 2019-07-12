apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    app: webhook
  name: webhook
  namespace: foghorn
spec:
  rules:
  - host: {{ (datasource "config").host }}
    http:
      paths:
      - backend:
          serviceName: webhook
          servicePort: 80
        path: /
