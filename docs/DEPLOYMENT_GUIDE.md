# ABT Dashboard Deployment Guide

## Overview
This guide covers deploying the ABT Dashboard from development to production environments. The dashboard is a lightweight Go application with static web assets that can be deployed in various configurations.

## Prerequisites

### System Requirements
- **Go**: Version 1.19 or higher
- **Memory**: Minimum 512MB RAM (1GB+ recommended for large datasets)
- **Storage**: 100MB+ free space
- **Network**: Port 8080 available (configurable)

### Development Dependencies
```bash
# Check Go version
go version

# Verify Git is available
git --version

# For building from source
make --version  # Optional, for Makefile builds
```

---

## Quick Start Deployment

### 1. Clone and Build
```bash
# Clone repository
git clone <repository-url>
cd abt-dashboard

# Download dependencies
go mod download

# Build application
go build -o abt-dashboard ./cmd/api

# Run with sample data
./abt-dashboard
```

### 2. Verify Installation
```bash
# Check server is running
curl http://localhost:8080/api/revenue/countries?limit=1

# Open dashboard in browser
open http://localhost:8080
```

---

## Production Deployment

### Environment Configuration

#### Environment Variables
Create a `.env` file or set environment variables:
```bash
# Server Configuration
ABT_PORT=8080                     # Server port
ABT_HOST=0.0.0.0                  # Bind address
ABT_ENV=production                # Environment
ABT_DATA_PATH=/data/dataset.csv   # CSV data file path

# Performance Settings
ABT_CACHE_TTL=300                 # Cache TTL in seconds
ABT_MAX_CONCURRENT=100            # Max concurrent requests
ABT_TIMEOUT=30                    # Request timeout in seconds
ABT_ENABLE_GZIP=true              # Enable gzip compression

# Logging
ABT_LOG_LEVEL=info                # Log level: debug, info, warn, error
ABT_LOG_FORMAT=json               # Log format: text, json
ABT_ACCESS_LOG=true               # Enable access logging

# Security
ABT_TRUSTED_PROXIES=10.0.0.0/8,172.16.0.0/12,192.168.0.0/16
ABT_RATE_LIMIT=100                # Requests per minute per IP
```

#### Configuration File (config.yaml)
```yaml
server:
  port: 8080
  host: "0.0.0.0"
  timeout: 30s
  read_timeout: 10s
  write_timeout: 10s

data:
  path: "/data/dataset.csv"
  cache_ttl: "5m"

performance:
  enable_gzip: true
  max_concurrent: 100
  cache_headers: true

logging:
  level: "info"
  format: "json"
  access_log: true

security:
  rate_limit: 100
  trusted_proxies:
    - "10.0.0.0/8"
    - "172.16.0.0/12"
    - "192.168.0.0/16"
```

---

## Deployment Methods

### 1. Binary Deployment

#### Build for Target Platform
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o abt-dashboard-linux ./cmd/api

# Linux ARM64 (for ARM servers)
GOOS=linux GOARCH=arm64 go build -o abt-dashboard-arm64 ./cmd/api

# Windows
GOOS=windows GOARCH=amd64 go build -o abt-dashboard.exe ./cmd/api

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o abt-dashboard-darwin ./cmd/api

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o abt-dashboard-darwin-arm64 ./cmd/api
```

#### Deployment Script
```bash
#!/bin/bash
# deploy.sh

set -e

# Configuration
APP_NAME="abt-dashboard"
DEPLOY_USER="abt"
DEPLOY_HOST="your-server.com"
DEPLOY_PATH="/opt/abt-dashboard"
SERVICE_NAME="abt-dashboard"

# Build for production
echo "Building application..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${APP_NAME} ./cmd/api

# Copy files to server
echo "Deploying to ${DEPLOY_HOST}..."
scp ${APP_NAME} ${DEPLOY_USER}@${DEPLOY_HOST}:${DEPLOY_PATH}/
scp dataset.csv ${DEPLOY_USER}@${DEPLOY_HOST}:${DEPLOY_PATH}/
scp -r web/ ${DEPLOY_USER}@${DEPLOY_HOST}:${DEPLOY_PATH}/

# Restart service
echo "Restarting service..."
ssh ${DEPLOY_USER}@${DEPLOY_HOST} "sudo systemctl restart ${SERVICE_NAME}"

echo "Deployment complete!"
```

---

### 2. Docker Deployment

#### Dockerfile
```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o abt-dashboard ./cmd/api

# Final stage
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh abt

# Set working directory
WORKDIR /app

# Copy binary and static files
COPY --from=builder /app/abt-dashboard .
COPY --from=builder /app/web ./web
COPY --from=builder /app/dataset.csv .

# Change ownership
RUN chown -R abt:abt /app

# Switch to non-root user
USER abt

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/revenue/countries?limit=1 || exit 1

# Run application
CMD ["./abt-dashboard"]
```

#### Docker Compose
```yaml
version: '3.8'

services:
  abt-dashboard:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ABT_ENV=production
      - ABT_LOG_LEVEL=info
      - ABT_CACHE_TTL=300
    volumes:
      - ./data:/data:ro  # Mount data directory
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/revenue/countries?limit=1"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # Optional: Nginx reverse proxy
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - abt-dashboard
    restart: unless-stopped
```

#### Build and Deploy
```bash
# Build image
docker build -t abt-dashboard:latest .

# Run container
docker run -d \
  --name abt-dashboard \
  -p 8080:8080 \
  -e ABT_ENV=production \
  -v $(pwd)/data:/data:ro \
  abt-dashboard:latest

# Using Docker Compose
docker-compose up -d

# Check logs
docker logs abt-dashboard

# Update deployment
docker-compose pull
docker-compose up -d
```

---

### 3. Kubernetes Deployment

#### Deployment Manifest
```yaml
# abt-dashboard-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: abt-dashboard
  labels:
    app: abt-dashboard
spec:
  replicas: 3
  selector:
    matchLabels:
      app: abt-dashboard
  template:
    metadata:
      labels:
        app: abt-dashboard
    spec:
      containers:
      - name: abt-dashboard
        image: abt-dashboard:latest
        ports:
        - containerPort: 8080
        env:
        - name: ABT_ENV
          value: "production"
        - name: ABT_LOG_LEVEL
          value: "info"
        - name: ABT_CACHE_TTL
          value: "300"
        volumeMounts:
        - name: data-volume
          mountPath: /data
          readOnly: true
        resources:
          limits:
            memory: "512Mi"
            cpu: "500m"
          requests:
            memory: "256Mi"
            cpu: "250m"
        livenessProbe:
          httpGet:
            path: /api/revenue/countries?limit=1
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /api/revenue/countries?limit=1
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: data-volume
        configMap:
          name: abt-data

---
apiVersion: v1
kind: Service
metadata:
  name: abt-dashboard-service
spec:
  selector:
    app: abt-dashboard
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: abt-dashboard-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - dashboard.yourdomain.com
    secretName: abt-dashboard-tls
  rules:
  - host: dashboard.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: abt-dashboard-service
            port:
              number: 80
```

#### ConfigMap for Data
```yaml
# abt-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: abt-data
data:
  dataset.csv: |
    # Your CSV data here
    # This would typically be loaded from a volume or external storage
```

#### Deploy to Kubernetes
```bash
# Create namespace
kubectl create namespace abt-dashboard

# Apply configurations
kubectl apply -f abt-configmap.yaml -n abt-dashboard
kubectl apply -f abt-dashboard-deployment.yaml -n abt-dashboard

# Check deployment
kubectl get pods -n abt-dashboard
kubectl logs -f deployment/abt-dashboard -n abt-dashboard

# Scale deployment
kubectl scale deployment abt-dashboard --replicas=5 -n abt-dashboard
```

---

## System Service Setup

### Systemd Service (Linux)

#### Service File
```ini
# /etc/systemd/system/abt-dashboard.service
[Unit]
Description=ABT Analytics Dashboard
After=network.target
Wants=network.target

[Service]
Type=simple
User=abt
Group=abt
WorkingDirectory=/opt/abt-dashboard
ExecStart=/opt/abt-dashboard/abt-dashboard
Restart=always
RestartSec=5

# Environment
Environment=ABT_ENV=production
Environment=ABT_LOG_LEVEL=info
Environment=ABT_PORT=8080

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/abt-dashboard/logs

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

#### Setup Commands
```bash
# Create user
sudo useradd -r -s /bin/false abt

# Create directories
sudo mkdir -p /opt/abt-dashboard
sudo mkdir -p /opt/abt-dashboard/logs

# Copy files
sudo cp abt-dashboard /opt/abt-dashboard/
sudo cp -r web /opt/abt-dashboard/
sudo cp dataset.csv /opt/abt-dashboard/

# Set permissions
sudo chown -R abt:abt /opt/abt-dashboard

# Install and start service
sudo systemctl daemon-reload
sudo systemctl enable abt-dashboard
sudo systemctl start abt-dashboard

# Check status
sudo systemctl status abt-dashboard
```

---

## Load Balancer Configuration

### Nginx Configuration
```nginx
# /etc/nginx/sites-available/abt-dashboard
upstream abt_backend {
    least_conn;
    server 127.0.0.1:8080 max_fails=3 fail_timeout=30s;
    server 127.0.0.1:8081 max_fails=3 fail_timeout=30s;
    server 127.0.0.1:8082 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name dashboard.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name dashboard.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/ssl/certs/dashboard.crt;
    ssl_certificate_key /etc/ssl/private/dashboard.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;

    # Security Headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Gzip Configuration
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/json;

    # Static Files
    location /web/ {
        alias /opt/abt-dashboard/web/;
        expires 1h;
        add_header Cache-Control "public, immutable";
    }

    # API Endpoints
    location /api/ {
        proxy_pass http://abt_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
        
        # Buffering
        proxy_buffering on;
        proxy_buffer_size 128k;
        proxy_buffers 4 256k;
        proxy_busy_buffers_size 256k;
    }

    # Root
    location / {
        proxy_pass http://abt_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Health Check
    location /health {
        proxy_pass http://abt_backend/api/revenue/countries?limit=1;
        access_log off;
    }
}
```

### HAProxy Configuration
```
# /etc/haproxy/haproxy.cfg
global
    daemon
    maxconn 4096
    log stdout local0 info

defaults
    mode http
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms
    option httplog
    option dontlognull

frontend abt_frontend
    bind *:80
    bind *:443 ssl crt /etc/ssl/certs/dashboard.pem
    redirect scheme https if !{ ssl_fc }
    
    # Rate limiting
    stick-table type ip size 100k expire 30s store http_req_rate(10s)
    http-request track-sc0 src
    http-request reject if { sc_http_req_rate(0) gt 20 }
    
    default_backend abt_backend

backend abt_backend
    balance roundrobin
    option httpchk GET /api/revenue/countries?limit=1
    http-check expect status 200
    
    server app1 127.0.0.1:8080 check
    server app2 127.0.0.1:8081 check
    server app3 127.0.0.1:8082 check
```

---

## Monitoring and Logging

### Application Logging
```go
// Enhanced logging configuration
type LogConfig struct {
    Level  string `yaml:"level"`  // debug, info, warn, error
    Format string `yaml:"format"` // text, json
    File   string `yaml:"file"`   // log file path
    Stdout bool   `yaml:"stdout"` // log to stdout
}

// Structured logging example
log.WithFields(log.Fields{
    "endpoint": "/api/revenue/countries",
    "duration": "250ms",
    "status": 200,
    "remote_ip": "192.168.1.100",
}).Info("Request processed")
```

### Health Check Endpoint
```bash
# Add health check endpoint
curl -f http://localhost:8080/health || exit 1

# Expected response
{
  "status": "healthy",
  "timestamp": "2025-08-22T19:00:00Z",
  "version": "1.0.0",
  "uptime": "24h30m15s"
}
```

### Monitoring with Prometheus
```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'abt-dashboard'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

---

## Security Hardening

### Firewall Configuration
```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw deny 8080/tcp   # Block direct app access
sudo ufw enable

# iptables
iptables -A INPUT -p tcp --dport 22 -j ACCEPT
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -s 127.0.0.1 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP
```

### SSL/TLS Configuration
```bash
# Generate self-signed certificate (development)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

# Let's Encrypt (production)
sudo certbot --nginx -d dashboard.yourdomain.com

# Update certificate renewal
sudo crontab -e
# Add: 0 2 * * * /usr/bin/certbot renew --quiet
```

### Rate Limiting
```go
// Rate limiting middleware
func rateLimitMiddleware(limit int) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Limit(limit), limit*2)
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Backup and Recovery

### Data Backup Strategy
```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/abt-dashboard"

# Create backup directory
mkdir -p ${BACKUP_DIR}

# Backup data file
cp /opt/abt-dashboard/dataset.csv ${BACKUP_DIR}/dataset_${DATE}.csv

# Backup configuration
cp /opt/abt-dashboard/config.yaml ${BACKUP_DIR}/config_${DATE}.yaml

# Backup logs (last 7 days)
find /opt/abt-dashboard/logs -name "*.log" -mtime -7 -exec cp {} ${BACKUP_DIR}/ \;

# Create archive
tar -czf ${BACKUP_DIR}/abt-dashboard_backup_${DATE}.tar.gz -C ${BACKUP_DIR} .

# Cleanup old backups (keep 30 days)
find ${BACKUP_DIR} -name "*.tar.gz" -mtime +30 -delete

echo "Backup completed: ${BACKUP_DIR}/abt-dashboard_backup_${DATE}.tar.gz"
```

### Disaster Recovery
```bash
#!/bin/bash
# restore.sh

BACKUP_FILE=$1
RESTORE_DIR="/opt/abt-dashboard"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file.tar.gz>"
    exit 1
fi

# Stop service
sudo systemctl stop abt-dashboard

# Create restore directory
sudo mkdir -p ${RESTORE_DIR}/restore

# Extract backup
sudo tar -xzf ${BACKUP_FILE} -C ${RESTORE_DIR}/restore

# Restore files
sudo cp ${RESTORE_DIR}/restore/dataset_*.csv ${RESTORE_DIR}/dataset.csv
sudo cp ${RESTORE_DIR}/restore/config_*.yaml ${RESTORE_DIR}/config.yaml

# Set permissions
sudo chown -R abt:abt ${RESTORE_DIR}

# Start service
sudo systemctl start abt-dashboard

# Verify restoration
curl -f http://localhost:8080/api/revenue/countries?limit=1

echo "Restoration completed"
```

---

## Performance Tuning

### Go Application Tuning
```bash
# Environment variables for production
export GOGC=100              # Garbage collection target
export GOMAXPROCS=4          # Number of CPU cores
export GOMEMLIMIT=1GiB       # Memory limit

# Build with optimizations
go build -ldflags="-s -w" -o abt-dashboard ./cmd/api
```

### System Tuning
```bash
# /etc/sysctl.d/99-abt-dashboard.conf
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 120
net.ipv4.tcp_keepalive_intvl = 30
net.ipv4.tcp_keepalive_probes = 3

# Apply changes
sudo sysctl -p /etc/sysctl.d/99-abt-dashboard.conf
```

### Database Connection Tuning (if applicable)
```go
// Database connection pool settings
db.SetMaxOpenConns(25)                 // Maximum open connections
db.SetMaxIdleConns(10)                 // Maximum idle connections
db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
```

---

## Troubleshooting

### Common Issues

#### High Memory Usage
```bash
# Check memory usage
ps aux | grep abt-dashboard
top -p $(pgrep abt-dashboard)

# Monitor with Go profiling
go tool pprof http://localhost:8080/debug/pprof/heap
```

#### Slow Response Times
```bash
# Check for blocking operations
go tool pprof http://localhost:8080/debug/pprof/goroutine

# Monitor request duration
curl -w "@curl-format.txt" http://localhost:8080/api/revenue/countries
```

#### Connection Issues
```bash
# Check listening ports
netstat -tlnp | grep 8080
ss -tlnp | grep 8080

# Test connectivity
telnet localhost 8080
```

### Log Analysis
```bash
# Real-time log monitoring
tail -f /opt/abt-dashboard/logs/app.log

# Error analysis
grep -i error /opt/abt-dashboard/logs/app.log | tail -20

# Performance analysis
grep "duration" /opt/abt-dashboard/logs/app.log | sort -k5 -n | tail -10
```

---

## Deployment Checklist

### Pre-deployment
- [ ] Go version 1.19+ installed
- [ ] Application builds successfully
- [ ] All tests pass
- [ ] Data file (dataset.csv) available
- [ ] Configuration reviewed
- [ ] SSL certificates ready (production)
- [ ] Firewall rules configured
- [ ] Monitoring setup prepared

### Deployment Steps
- [ ] Build application for target platform
- [ ] Create deployment user and directories
- [ ] Copy application files
- [ ] Set proper file permissions
- [ ] Configure system service
- [ ] Start application service
- [ ] Configure reverse proxy (if applicable)
- [ ] Test all endpoints
- [ ] Verify performance requirements
- [ ] Setup monitoring and logging
- [ ] Configure backups

### Post-deployment
- [ ] Health checks passing
- [ ] All features functional
- [ ] Performance within targets (<10s)
- [ ] Logs properly configured
- [ ] Monitoring alerts setup
- [ ] Backup strategy tested
- [ ] Documentation updated
- [ ] Team access configured

---

## Support and Maintenance

### Regular Maintenance Tasks
```bash
# Weekly log rotation
sudo logrotate /etc/logrotate.d/abt-dashboard

# Monthly performance review
curl -w "@curl-format.txt" http://localhost:8080/api/revenue/countries > performance_$(date +%Y%m).log

# Security updates
sudo apt update && sudo apt upgrade
go version  # Check for Go updates
```

### Emergency Procedures
```bash
# Quick restart
sudo systemctl restart abt-dashboard

# Emergency stop
sudo systemctl stop abt-dashboard
sudo pkill -f abt-dashboard

# Quick health check
curl -f http://localhost:8080/health || echo "Service down"
```

---

**ABT Dashboard - Enterprise-Ready Deployment Guide**
