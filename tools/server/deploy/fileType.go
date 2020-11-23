package main

type FileModel struct {
}

const (
	Deployment = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: {{Name}}
  namespace: {{Namespace}}
spec:
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: {{replicas}}
  template:
    metadata:
      namespace: {{Namespace}}
      labels:
        app: {{Name}}
    spec:
      containers:
      - name: {{Name}}
        image: [url]/{{author}}/{{Name}}:{{version}}
        imagePullPolicy: Always
        resources:
          # keep request = limit to keep this container in guaranteed class
          limits:
            cpu: [cpuLimit]
            memory: [memoryLimit]
          requests:
            cpu: [cpuRequest]
            memory: [memoryRequest]
        volumeMounts:
        - name: logs
          mountPath: [logPath]
        - name: config
          mountPath: /config
        command: ["/bin/sh", "-c"]
        args: ["[cmdArgs]"]
      volumes:
      - name: logs
        hostPath:
          path: [logTargetPath]
      - name: config
        configMap:
          name: [configMapName]
          items:
          - key: [configName]
            path: default.json
`

	Ingress string = `---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name:  {{Name}}
  namespace: {{Namespace}}
  annotations:
    kubernetes.io/ingress.class: [proxyStatus]
spec:
  rules:
  - host:  [domain]
    http:
      paths:
      - path: /
        backend:
          serviceName: {{Name}}
          servicePort: [port]
`

	Service string = `---
apiVersion: v1
kind: Service
metadata:
  name: {{Name}}
  namespace: {{Namespace}}
  annotations:
    prometheus.io/scrape: 'true'
    prometheus.io/path:   /metrics
    prometheus.io/port:   '8081'
  labels:
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "{{Name}}"
spec:
  type: ClusterIP
  ports:
    - name: {{Name}}
      port: [port]
      targetPort: [servicePort]
      protocol: TCP
  selector:
    app: {{Name}}
  sessionAffinity: ClientIP
`

	Dockerfile_Go string = `FROM reg.hoper.xyz/library/centos:7
MAINTAINER "Hoper Team <erp@hyx.com>"

#修改容器时区
RUN rm /etc/localtime && \
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /{{Name}}

ADD ./build/main /{{Name}}/

ADD ./build/ /{{Name}}/build/

RUN mkdir /data && cd /data && mkdir logs

CMD ["./{{Name}}", "-conf", "config"]
`

	Dockerfile_Jar string = `FROM reg.hoper.xyz/library/centos7-jre8:0.1.0
#FROM reg.hoper.xyz/yuwanchat/jdk:v1.8
MAINTAINER "Hoper Aura Dts <dts@hyx.com>"

#修改容器时区
RUN rm /etc/localtime && \
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /{{Name}}

ADD ./build /{{Name}}

#RUN java -jar {{Name}}.jar
`
)
