# Примеры Kubernetes манифестов для BirthdayWish

## Обзор перехода от Docker Compose к Kubernetes

Docker Compose отлично подходит для разработки, но для продакшена нужен Kubernetes. Основные отличия:

| Docker Compose | Kubernetes |
|----------------|------------|
| `services` | `Deployments` + `Services` |
| `ports` | `Services` + `Ingress` |
| `volumes` | `PersistentVolumes` + `PersistentVolumeClaims` |
| `environment` | `ConfigMaps` + `Secrets` |
| `networks` | `Services` (автоматически) |
| `depends_on` | `initContainers` + `readinessProbes` |

## Структура Kubernetes манифестов

```
k8s/
├── namespaces/
│   └── birthdaywish-namespace.yaml
├── configmaps/
│   ├── app-config.yaml
│   └── monitoring-config.yaml
├── secrets/
│   └── database-secrets.yaml
├── databases/
│   ├── postgres-auth-deployment.yaml
│   ├── postgres-auth-service.yaml
│   ├── postgres-wish-deployment.yaml
│   ├── postgres-wish-service.yaml
│   ├── minio-deployment.yaml
│   └── minio-service.yaml
├── services/
│   ├── auth-deployment.yaml
│   ├── auth-service.yaml
│   ├── wishlister-deployment.yaml
│   ├── wishlister-service.yaml
│   ├── filer-deployment.yaml
│   ├── filer-service.yaml
│   ├── gateway-deployment.yaml
│   └── gateway-service.yaml
├── ingress/
│   └── birthdaywish-ingress.yaml
└── monitoring/
    ├── prometheus-deployment.yaml
    ├── prometheus-service.yaml
    ├── grafana-deployment.yaml
    └── grafana-service.yaml
```

## Namespace

**k8s/namespaces/birthdaywish-namespace.yaml**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: birthdaywish
  labels:
    name: birthdaywish
    app: birthdaywish
```

## ConfigMaps

**k8s/configmaps/app-config.yaml**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: birthdaywish
data:
  # Database configuration
  POSTGRES_DB_USERS: "users_service"
  POSTGRES_DB_WISHES: "wish_service"
  
  # Service ports
  AUTH_PORT: "50051"
  WISHLISTER_PORT: "50052"
  FILER_PORT: "50053"
  GATEWAY_PORT: "3030"
  
  # MinIO configuration
  MINIO_BUCKET: "birthdaywish-files"
  
  # Logging
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
```

## Secrets

**k8s/secrets/database-secrets.yaml**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: database-secrets
  namespace: birthdaywish
type: Opaque
data:
  # Base64 encoded values
  postgres-user: cG9zdGdyZXM=  # postgres
  postgres-password: c2VjdXJlX3Bhc3N3b3JkXzEyMw==  # secure_password_123
  minio-user: bWluaW9hZG1pbg==  # minioadmin
  minio-password: c2VjdXJlX21pbmlvX3Bhc3N3b3Jk  # secure_minio_password
```

## PostgreSQL для Auth Service

**k8s/databases/postgres-auth-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres-auth
  namespace: birthdaywish
  labels:
    app: postgres-auth
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres-auth
  template:
    metadata:
      labels:
        app: postgres-auth
    spec:
      containers:
      - name: postgres
        image: postgres:17-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: database-secrets
              key: postgres-user
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: database-secrets
              key: postgres-password
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: POSTGRES_DB_USERS
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        livenessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - postgres
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - postgres
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-auth-pvc
```

**k8s/databases/postgres-auth-service.yaml**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: postgres-auth
  namespace: birthdaywish
spec:
  selector:
    app: postgres-auth
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
```

**k8s/databases/postgres-auth-pvc.yaml**
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-auth-pvc
  namespace: birthdaywish
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: standard
```

## Auth Service

**k8s/services/auth-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: birthdaywish
  labels:
    app: auth-service
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: birthdaywish/auth-service:latest
        ports:
        - containerPort: 50051
        env:
        - name: DB_HOST
          value: "postgres-auth"
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: database-secrets
              key: postgres-user
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: database-secrets
              key: postgres-password
        - name: DB_NAME
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: POSTGRES_DB_USERS
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: LOG_LEVEL
        livenessProbe:
          exec:
            command:
            - /bin/sh
            - -c
            - "nc -z localhost 50051"
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - /bin/sh
            - -c
            - "nc -z localhost 50051"
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
```

**k8s/services/auth-service.yaml**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: birthdaywish
spec:
  selector:
    app: auth-service
  ports:
  - port: 50051
    targetPort: 50051
    protocol: TCP
  type: ClusterIP
```

## Gateway Service

**k8s/services/gateway-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway-service
  namespace: birthdaywish
  labels:
    app: gateway-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gateway-service
  template:
    metadata:
      labels:
        app: gateway-service
    spec:
      containers:
      - name: gateway-service
        image: birthdaywish/gateway-service:latest
        ports:
        - containerPort: 3030
        env:
        - name: AUTH_SERVICE_URL
          value: "auth-service:50051"
        - name: WISHLISTER_SERVICE_URL
          value: "wishlister-service:50052"
        - name: FILER_SERVICE_URL
          value: "filer-service:50053"
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: LOG_LEVEL
        livenessProbe:
          httpGet:
            path: /health
            port: 3030
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 3030
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
```

**k8s/services/gateway-service.yaml**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: gateway-service
  namespace: birthdaywish
spec:
  selector:
    app: gateway-service
  ports:
  - port: 3030
    targetPort: 3030
    protocol: TCP
  type: ClusterIP
```

## Ingress

**k8s/ingress/birthdaywish-ingress.yaml**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: birthdaywish-ingress
  namespace: birthdaywish
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - api.birthdaywish.com
    secretName: birthdaywish-tls
  rules:
  - host: api.birthdaywish.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: gateway-service
            port:
              number: 3030
```

## Prometheus для мониторинга

**k8s/monitoring/prometheus-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
  namespace: birthdaywish
  labels:
    app: prometheus
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
      - name: prometheus
        image: prom/prometheus:v2.47.2
        ports:
        - containerPort: 9090
        args:
        - '--config.file=/etc/prometheus/prometheus.yml'
        - '--storage.tsdb.path=/prometheus/'
        - '--web.console.libraries=/etc/prometheus/console_libraries'
        - '--web.console.templates=/etc/prometheus/consoles'
        - '--storage.tsdb.retention.time=30d'
        - '--web.enable-lifecycle'
        volumeMounts:
        - name: prometheus-config
          mountPath: /etc/prometheus
        - name: prometheus-storage
          mountPath: /prometheus
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
      volumes:
      - name: prometheus-config
        configMap:
          name: prometheus-config
      - name: prometheus-storage
        persistentVolumeClaim:
          claimName: prometheus-pvc
```

## Команды для развертывания

### Создание ресурсов

```bash
# Создание namespace
kubectl apply -f k8s/namespaces/

# Создание конфигурации
kubectl apply -f k8s/configmaps/
kubectl apply -f k8s/secrets/

# Создание баз данных
kubectl apply -f k8s/databases/

# Создание сервисов
kubectl apply -f k8s/services/

# Создание ingress
kubectl apply -f k8s/ingress/

# Создание мониторинга
kubectl apply -f k8s/monitoring/
```

### Проверка статуса

```bash
# Проверка всех ресурсов
kubectl get all -n birthdaywish

# Проверка подов
kubectl get pods -n birthdaywish

# Проверка сервисов
kubectl get services -n birthdaywish

# Проверка ingress
kubectl get ingress -n birthdaywish

# Логи подов
kubectl logs -f deployment/auth-service -n birthdaywish
```

### Масштабирование

```bash
# Увеличение количества реплик
kubectl scale deployment gateway-service --replicas=5 -n birthdaywish

# Автоматическое масштабирование
kubectl autoscale deployment gateway-service --cpu-percent=70 --min=2 --max=10 -n birthdaywish
```

## Helm Charts (продвинутый уровень)

Для более удобного управления можно использовать Helm:

**Chart.yaml**
```yaml
apiVersion: v2
name: birthdaywish
description: BirthdayWish microservices application
version: 0.1.0
appVersion: "1.0.0"
```

**values.yaml**
```yaml
replicaCount: 1

image:
  repository: birthdaywish
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 3030

ingress:
  enabled: true
  className: nginx
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
  hosts:
    - host: api.birthdaywish.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: birthdaywish-tls
      hosts:
        - api.birthdaywish.com

resources:
  limits:
    cpu: 200m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

## Мониторинг и логирование

### Prometheus + Grafana

```bash
# Установка Prometheus Operator
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack -n monitoring --create-namespace
```

### ELK Stack для логов

```bash
# Установка Elasticsearch, Logstash, Kibana
helm repo add elastic https://helm.elastic.co
helm install elasticsearch elastic/elasticsearch -n logging --create-namespace
helm install kibana elastic/kibana -n logging
```

## Безопасность

### Network Policies

**k8s/security/network-policy.yaml**
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: birthdaywish-network-policy
  namespace: birthdaywish
spec:
  podSelector:
    matchLabels:
      app: birthdaywish
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: birthdaywish
    ports:
    - protocol: TCP
      port: 3030
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: birthdaywish
    ports:
    - protocol: TCP
      port: 50051
    - protocol: TCP
      port: 50052
    - protocol: TCP
      port: 50053
```

### RBAC

**k8s/security/rbac.yaml**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: birthdaywish-sa
  namespace: birthdaywish
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: birthdaywish
  name: birthdaywish-role
rules:
- apiGroups: [""]
  resources: ["pods", "services", "configmaps", "secrets"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: birthdaywish-rolebinding
  namespace: birthdaywish
subjects:
- kind: ServiceAccount
  name: birthdaywish-sa
  namespace: birthdaywish
roleRef:
  kind: Role
  name: birthdaywish-role
  apiGroup: rbac.authorization.k8s.io
```

Этот набор манифестов обеспечивает полный переход от Docker Compose к Kubernetes с сохранением всей функциональности и добавлением возможностей оркестрации, масштабирования и мониторинга.
