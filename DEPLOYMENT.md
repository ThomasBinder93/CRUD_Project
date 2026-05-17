# Deployment Guide

This guide explains how to deploy the CRUD application to various environments.

## Local Development

### With Make
```bash
make run
```

### With Go
```bash
go mod tidy
go run main.go
```

### With Docker Compose
```bash
docker-compose up
```

The application will be available at `http://localhost:8080`

## Production Deployment

### Docker Standalone

1. **Build the image**:
```bash
docker build -t crud-app:1.0.0 .
```

2. **Run the container**:
```bash
docker run -d \
  --name crud-app \
  -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e LOG_LEVEL=info \
  -v crud-data:/data \
  crud-app:1.0.0
```

3. **Check health**:
```bash
curl http://localhost:8080/health
```

### Docker Compose (Production)

Create `docker-compose.prod.yml`:

```yaml
version: "3.9"

services:
  app:
    image: crud-app:1.0.0
    ports:
      - "8080:8080"
    environment:
      SERVER_PORT: "8080"
      SERVER_READ_TIMEOUT: "30"
      SERVER_WRITE_TIMEOUT: "30"
      DB_PATH: "/data/crud.db"
      LOG_LEVEL: "warn"
    volumes:
      - crud-data:/data
    restart: always
    container_name: crud-app
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - backend

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    networks:
      - backend

volumes:
  crud-data:

networks:
  backend:
```

Deploy:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

## Kubernetes Deployment

### Prerequisites
- kubectl configured
- Kubernetes cluster running
- Docker image pushed to registry

### 1. Create Namespace
```bash
kubectl create namespace crud-app
```

### 2. Create ConfigMap
```bash
kubectl create configmap crud-config \
  --from-literal=SERVER_PORT=8080 \
  --from-literal=LOG_LEVEL=info \
  -n crud-app
```

### 3. Create Persistent Volume
```yaml
# pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: crud-data-pv
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/mnt/data"

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: crud-data-pvc
  namespace: crud-app
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

Deploy:
```bash
kubectl apply -f pv.yaml
```

### 4. Deploy Application

Update `deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: crud-app
  namespace: crud-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: crud-app
  template:
    metadata:
      labels:
        app: crud-app
    spec:
      containers:
      - name: crud-app
        image: your-registry/crud-app:1.0.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_PORT
          value: "8080"
        - name: LOG_LEVEL
          value: "info"
        - name: DB_PATH
          value: "/data/crud.db"
        volumeMounts:
        - name: data
          mountPath: /data
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: crud-data-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: crud-app
  namespace: crud-app
spec:
  selector:
    app: crud-app
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: crud-app-hpa
  namespace: crud-app
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: crud-app
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

Deploy:
```bash
kubectl apply -f deployment.yaml
```

### 5. Verify Deployment
```bash
# Check pods
kubectl get pods -n crud-app

# Check service
kubectl get svc -n crud-app

# View logs
kubectl logs -n crud-app deployment/crud-app

# Port forward for testing
kubectl port-forward -n crud-app svc/crud-app 8080:80
```

## Environment Variables

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `SERVER_PORT` | 8080 | Server port | 3000 |
| `SERVER_READ_TIMEOUT` | 15 | Read timeout (seconds) | 30 |
| `SERVER_WRITE_TIMEOUT` | 15 | Write timeout (seconds) | 30 |
| `DB_PATH` | ./crud.db | Database path | /data/crud.db |
| `LOG_LEVEL` | info | Log level | debug, info, warn, error |

## Scaling Considerations

### Vertical Scaling (More Powerful Machine)
- Increase CPU/memory limits in Kubernetes
- Increase connection pool size

### Horizontal Scaling (More Instances)
- Use Kubernetes HPA (Horizontal Pod Autoscaler)
- Use load balancer (Nginx, HAProxy)
- Ensure stateless operations
- Consider database scaling

### Database Scaling
For production with high load:
1. Migrate from SQLite to PostgreSQL
2. Implement connection pooling (pgbouncer)
3. Set up read replicas
4. Consider sharding for very large datasets

## Monitoring and Logging

### Health Checks
The application exposes a health endpoint:
```bash
curl http://your-app:8080/health
```

### Structured Logging
Logs are output in structured format:
```
time=2024-04-26T10:30:00.123Z level=INFO msg="request handled" method=GET path=/api/items status=200 duration=5ms
```

### Log Collection
For production, integrate with:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Splunk
- DataDog
- New Relic

Example filebeat config:
```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/crud-app/*.log

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

## Backup and Recovery

### Database Backup
```bash
# Backup SQLite database
cp /data/crud.db /backups/crud.db.backup

# Or with Kubernetes
kubectl exec -n crud-app crud-app-pod -- cp /data/crud.db /data/backup/crud.db
```

### Automated Backups
Create a CronJob:
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: crud-backup
  namespace: crud-app
spec:
  schedule: "0 2 * * *"  # 2 AM daily
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: alpine:latest
            command: ["sh", "-c", "cp /data/crud.db /backups/crud.db.$(date +%s)"]
            volumeMounts:
            - name: data
              mountPath: /data
            - name: backups
              mountPath: /backups
          volumes:
          - name: data
            persistentVolumeClaim:
              claimName: crud-data-pvc
          - name: backups
            hostPath:
              path: /var/backups/crud
          restartPolicy: OnFailure
```

## SSL/TLS Configuration

### With Docker
```bash
docker run -d \
  -p 80:8080 \
  -p 443:8443 \
  -v /etc/ssl/certs:/app/certs:ro \
  crud-app:1.0.0
```

### With Kubernetes (cert-manager)
```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml

# Create certificate
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: crud-app-cert
  namespace: crud-app
spec:
  secretName: crud-app-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
  - crud-app.example.com
EOF
```

## Troubleshooting

### Application Won't Start
```bash
# Check logs
docker logs crud-app
# or
kubectl logs -n crud-app deployment/crud-app

# Check configuration
docker inspect crud-app
# or
kubectl describe pod -n crud-app <pod-name>
```

### High Memory Usage
- Check for memory leaks
- Reduce database connection pool
- Implement request timeout

### High CPU Usage
- Check log level (reduce if verbose)
- Profile the application
- Add caching layer

### Database Locked
- Ensure single writer
- Check for long transactions
- Monitor database usage

## Rollback

### Docker
```bash
docker stop crud-app
docker rm crud-app
docker run -d <previous-image:tag>
```

### Kubernetes
```bash
# Check rollout history
kubectl rollout history deployment/crud-app -n crud-app

# Rollback to previous version
kubectl rollout undo deployment/crud-app -n crud-app

# Rollback to specific revision
kubectl rollout undo deployment/crud-app --to-revision=2 -n crud-app
```

## Performance Optimization

1. **Enable caching**:
   - Add Redis for frequently accessed items
   - Implement HTTP caching headers

2. **Database optimization**:
   - Add indexes on frequently queried columns
   - Implement query caching

3. **API optimization**:
   - Implement pagination
   - Add request/response compression
   - Use connection pooling

4. **Application optimization**:
   - Profile for bottlenecks
   - Optimize hot paths
   - Use buffering for I/O operations

## Maintenance

### Regular Tasks
- Monitor disk space
- Check error logs for issues
- Review performance metrics
- Update dependencies
- Test disaster recovery procedures

### Scheduled Maintenance
- Run database maintenance (vacuum, analyze)
- Update application to latest version
- Test and validate backups
- Security patches

## Support and Updates

For updates and patches:
```bash
# Check for new version
git fetch origin

# Pull latest changes
git pull origin main

# Build and deploy new version
make docker-build
# Deploy with your chosen method
```

---

For additional help, see README.md, API_DOCUMENTATION.md, and ARCHITECTURE.md
