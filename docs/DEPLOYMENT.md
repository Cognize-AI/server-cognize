# Deployment Guide

This guide covers various deployment strategies for the Cognize API server in production environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Docker Deployment](#docker-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Manual Deployment](#manual-deployment)
- [Database Setup](#database-setup)
- [Monitoring and Logging](#monitoring-and-logging)
- [Security Considerations](#security-considerations)
- [Performance Optimization](#performance-optimization)
- [Backup and Recovery](#backup-and-recovery)

## Prerequisites

### System Requirements

- **CPU**: 2+ cores recommended
- **Memory**: 4GB+ RAM recommended
- **Storage**: 20GB+ available space
- **Network**: HTTPS support (SSL/TLS certificates)

### Dependencies

- Docker & Docker Compose (for containerized deployment)
- PostgreSQL 17+ (managed service recommended)
- Load balancer (for high availability)
- SSL/TLS certificates

## Environment Configuration

### Production Environment Variables

Create a production `.env` file with secure values:

```env
# Server Configuration
PORT=4000
ENVIRONMENT=prod

# Database Configuration (Use managed PostgreSQL service)
DB_STRING=postgres://username:password@db-host:5432/cognize?sslmode=require

# JWT Configuration (Use strong, random secret)
JWT_SECRET=your-super-secure-random-jwt-secret-256-bits-long

# Google OAuth Configuration
GOOGLE_OAUTH_CLIENT_ID=your-production-google-oauth-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-production-google-oauth-client-secret
GOOGLE_OAUTH_REDIRECT_URL=https://yourdomain.com/auth/callback

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=noreply@yourdomain.com
SMTP_PASSWORD=your-app-password

# Logging Configuration (Recommended for production)
AXIOM_TOKEN=your-axiom-token
AXIOM_ORG=your-axiom-org
AXIOM_DATASET=cognize-production-logs

# Encryption
ENC_SECRET=your-encryption-secret-key-256-bits
```

### Security Best Practices

1. **Never commit `.env` files** to version control
2. **Use strong, random secrets** for JWT and encryption
3. **Enable SSL/TLS** for all connections
4. **Use managed database services** with encryption at rest
5. **Implement proper access controls** and network security

## Docker Deployment

### Single Server Deployment

1. **Prepare the server:**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   
   # Install Docker Compose
   sudo apt install docker-compose -y
   ```

2. **Create deployment directory:**
   ```bash
   mkdir /opt/cognize
   cd /opt/cognize
   ```

3. **Create production docker-compose.yml:**
   ```yaml
   version: '3.8'
   
   services:
     postgres:
       image: postgres:17
       container_name: cognize_postgres
       restart: unless-stopped
       environment:
         POSTGRES_USER: ${DB_USER}
         POSTGRES_PASSWORD: ${DB_PASSWORD}
         POSTGRES_DB: ${DB_NAME}
       ports:
         - "5432:5432"
       volumes:
         - postgres_data:/var/lib/postgresql/data
         - ./backup:/backup
       networks:
         - cognize_network
   
     app:
       image: cognize-api:latest
       container_name: cognize_api
       restart: unless-stopped
       ports:
         - "4000:4000"
       environment:
         - PORT=4000
         - DB_STRING=postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
         - JWT_SECRET=${JWT_SECRET}
         - GOOGLE_OAUTH_CLIENT_ID=${GOOGLE_OAUTH_CLIENT_ID}
         - GOOGLE_OAUTH_CLIENT_SECRET=${GOOGLE_OAUTH_CLIENT_SECRET}
         - GOOGLE_OAUTH_REDIRECT_URL=${GOOGLE_OAUTH_REDIRECT_URL}
         - ENVIRONMENT=prod
       depends_on:
         - postgres
       networks:
         - cognize_network
   
     nginx:
       image: nginx:alpine
       container_name: cognize_nginx
       restart: unless-stopped
       ports:
         - "80:80"
         - "443:443"
       volumes:
         - ./nginx.conf:/etc/nginx/nginx.conf
         - ./ssl:/etc/nginx/ssl
       depends_on:
         - app
       networks:
         - cognize_network
   
   volumes:
     postgres_data:
   
   networks:
     cognize_network:
       driver: bridge
   ```

4. **Create Nginx configuration:**
   ```nginx
   events {
       worker_connections 1024;
   }
   
   http {
       upstream app {
           server app:4000;
       }
   
       server {
           listen 80;
           server_name yourdomain.com;
           return 301 https://$server_name$request_uri;
       }
   
       server {
           listen 443 ssl http2;
           server_name yourdomain.com;
   
           ssl_certificate /etc/nginx/ssl/cert.pem;
           ssl_certificate_key /etc/nginx/ssl/key.pem;
           ssl_protocols TLSv1.2 TLSv1.3;
           ssl_ciphers HIGH:!aNULL:!MD5;
   
           location / {
               proxy_pass http://app;
               proxy_set_header Host $host;
               proxy_set_header X-Real-IP $remote_addr;
               proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
               proxy_set_header X-Forwarded-Proto $scheme;
           }
       }
   }
   ```

5. **Deploy:**
   ```bash
   # Build the application image
   docker build -t cognize-api:latest .
   
   # Start services
   docker-compose up -d
   
   # Check status
   docker-compose ps
   ```

### High Availability Deployment

For production environments requiring high availability:

1. **Use multiple app instances:**
   ```yaml
   app:
     image: cognize-api:latest
     restart: unless-stopped
     deploy:
       replicas: 3
   ```

2. **Use external load balancer** (AWS ALB, GCP Load Balancer, etc.)

3. **Use managed database services** (AWS RDS, GCP Cloud SQL, etc.)

## Cloud Deployment

### AWS Deployment

#### Using Amazon ECS

1. **Create ECS Task Definition:**
   ```json
   {
     "family": "cognize-api",
     "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
     "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "512",
     "memory": "1024",
     "containerDefinitions": [
       {
         "name": "cognize-api",
         "image": "your-account.dkr.ecr.region.amazonaws.com/cognize-api:latest",
         "portMappings": [
           {
             "containerPort": 4000,
             "protocol": "tcp"
           }
         ],
         "environment": [
           {
             "name": "PORT",
             "value": "4000"
           }
         ],
         "secrets": [
           {
             "name": "DB_STRING",
             "valueFrom": "arn:aws:secretsmanager:region:account:secret:cognize/db-string"
           }
         ]
       }
     ]
   }
   ```

2. **Create ECS Service with Application Load Balancer**

3. **Use Amazon RDS for PostgreSQL**

4. **Use AWS Secrets Manager for sensitive configuration**

#### Using AWS Lambda (Serverless)

For serverless deployment, you'll need to adapt the application for Lambda:

1. Create a Lambda handler
2. Use Amazon API Gateway
3. Use Amazon RDS Proxy for database connections

### Google Cloud Platform

#### Using Cloud Run

1. **Build and push to Container Registry:**
   ```bash
   # Build image
   docker build -t gcr.io/your-project/cognize-api .
   
   # Push to registry
   docker push gcr.io/your-project/cognize-api
   ```

2. **Deploy to Cloud Run:**
   ```bash
   gcloud run deploy cognize-api \
     --image gcr.io/your-project/cognize-api \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated \
     --port 4000 \
     --set-env-vars ENVIRONMENT=prod
   ```

3. **Use Cloud SQL for PostgreSQL**

4. **Use Secret Manager for configuration**

### Azure Deployment

#### Using Azure Container Instances

1. **Create resource group and container instance**

2. **Use Azure Database for PostgreSQL**

3. **Use Azure Key Vault for secrets**

## Manual Deployment

### Linux Server Deployment

1. **Prepare the server:**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install required packages
   sudo apt install postgresql-client nginx certbot -y
   ```

2. **Install Go:**
   ```bash
   wget https://go.dev/dl/go1.24.6.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.24.6.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **Deploy application:**
   ```bash
   # Create app user
   sudo useradd -m -s /bin/bash cognize
   
   # Create application directory
   sudo mkdir -p /opt/cognize
   sudo chown cognize:cognize /opt/cognize
   
   # Clone and build
   sudo -u cognize git clone https://github.com/Cognize-AI/server-cognize.git /opt/cognize
   cd /opt/cognize
   sudo -u cognize go build -o cognize .
   ```

4. **Create systemd service:**
   ```ini
   [Unit]
   Description=Cognize API Server
   After=network.target
   
   [Service]
   Type=simple
   User=cognize
   WorkingDirectory=/opt/cognize
   ExecStart=/opt/cognize/cognize
   EnvironmentFile=/opt/cognize/.env
   Restart=on-failure
   RestartSec=5
   
   [Install]
   WantedBy=multi-user.target
   ```

5. **Start service:**
   ```bash
   sudo systemctl enable cognize
   sudo systemctl start cognize
   sudo systemctl status cognize
   ```

## Database Setup

### Managed Database Services (Recommended)

#### AWS RDS PostgreSQL

1. Create RDS instance with encryption
2. Configure security groups for application access
3. Set up automated backups
4. Configure monitoring and alerts

#### Google Cloud SQL

1. Create Cloud SQL PostgreSQL instance
2. Configure authorized networks
3. Enable automatic backups
4. Set up monitoring

#### Azure Database for PostgreSQL

1. Create Azure Database for PostgreSQL
2. Configure firewall rules
3. Enable automated backups
4. Set up monitoring

### Self-Managed PostgreSQL

If using self-managed PostgreSQL:

1. **Install PostgreSQL:**
   ```bash
   sudo apt install postgresql postgresql-contrib -y
   ```

2. **Configure PostgreSQL:**
   ```bash
   # Create database and user
   sudo -u postgres createdb cognize
   sudo -u postgres createuser cognize_user
   sudo -u postgres psql -c "ALTER USER cognize_user PASSWORD 'secure_password';"
   sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE cognize TO cognize_user;"
   ```

3. **Configure security:**
   - Enable SSL/TLS
   - Configure pg_hba.conf for secure access
   - Set up regular backups

## Monitoring and Logging

### Application Monitoring

1. **Health Check Endpoint:**
   Add a health check endpoint to your application:
   ```go
   r.GET("/health", func(c *gin.Context) {
       c.JSON(200, gin.H{"status": "healthy"})
   })
   ```

2. **Metrics Collection:**
   - Use Prometheus for metrics collection
   - Set up Grafana for visualization
   - Monitor key metrics: response time, error rate, throughput

3. **Alerting:**
   - Set up alerts for high error rates
   - Monitor database connection health
   - Alert on high response times

### Logging

1. **Centralized Logging:**
   - Use Axiom, ELK stack, or cloud logging services
   - Structure logs with consistent formats
   - Include correlation IDs for request tracing

2. **Log Levels:**
   - INFO: General application flow
   - WARN: Unexpected situations
   - ERROR: Error conditions
   - DEBUG: Detailed debugging (only in development)

## Security Considerations

### Network Security

1. **Use HTTPS everywhere**
2. **Configure firewalls** to allow only necessary ports
3. **Use VPCs/VNets** for network isolation
4. **Implement rate limiting**

### Application Security

1. **Validate all inputs**
2. **Use parameterized queries** (GORM handles this)
3. **Implement proper CORS policies**
4. **Regular security updates**

### Database Security

1. **Use strong passwords**
2. **Enable encryption at rest and in transit**
3. **Regular backups with encryption**
4. **Principle of least privilege** for database access

## Performance Optimization

### Application Performance

1. **Database Connection Pooling:**
   ```go
   // Configure in dbConfig.go
   db.DB().SetMaxOpenConns(25)
   db.DB().SetMaxIdleConns(25)
   db.DB().SetConnMaxLifetime(5 * time.Minute)
   ```

2. **Caching:**
   - Implement Redis for session caching
   - Cache frequently accessed data
   - Use CDN for static assets

3. **Database Optimization:**
   - Add appropriate indexes
   - Optimize queries
   - Use read replicas for read-heavy workloads

### Infrastructure Performance

1. **Load Balancing:**
   - Distribute traffic across multiple instances
   - Use health checks
   - Configure proper timeouts

2. **Auto Scaling:**
   - Set up auto-scaling based on CPU/memory usage
   - Configure appropriate scaling policies

## Backup and Recovery

### Database Backups

1. **Automated Backups:**
   ```bash
   # Daily backup script
   #!/bin/bash
   BACKUP_DIR="/backup/$(date +%Y%m%d)"
   mkdir -p $BACKUP_DIR
   pg_dump -h localhost -U cognize_user cognize > $BACKUP_DIR/cognize_backup.sql
   gzip $BACKUP_DIR/cognize_backup.sql
   ```

2. **Point-in-Time Recovery:**
   - Enable WAL archiving
   - Regular backup testing
   - Document recovery procedures

### Disaster Recovery

1. **Multi-Region Deployment:**
   - Deploy in multiple regions
   - Use database replication
   - Implement automated failover

2. **Recovery Testing:**
   - Regular disaster recovery drills
   - Documented recovery procedures
   - RTO/RPO targets

## Troubleshooting

### Common Issues

1. **Database Connection Issues:**
   ```bash
   # Check database connectivity
   pg_isready -h db_host -p 5432 -U username
   ```

2. **Memory Issues:**
   ```bash
   # Monitor memory usage
   free -h
   docker stats
   ```

3. **Log Analysis:**
   ```bash
   # Check application logs
   docker logs cognize_api
   sudo journalctl -u cognize -f
   ```

### Performance Issues

1. **Slow Queries:**
   - Enable PostgreSQL query logging
   - Use EXPLAIN ANALYZE for query optimization
   - Add appropriate indexes

2. **High CPU Usage:**
   - Profile application with pprof
   - Optimize CPU-intensive operations
   - Scale horizontally

## Maintenance

### Regular Maintenance Tasks

1. **Security Updates:**
   - Regular OS and package updates
   - Update Go and dependencies
   - Security vulnerability scanning

2. **Database Maintenance:**
   - Regular VACUUM and ANALYZE
   - Index maintenance
   - Statistics updates

3. **Log Rotation:**
   - Configure log rotation
   - Archive old logs
   - Monitor disk usage

### Deployment Updates

1. **Blue-Green Deployment:**
   - Deploy to new environment
   - Test thoroughly
   - Switch traffic gradually

2. **Rolling Updates:**
   - Update instances one by one
   - Health check before proceeding
   - Rollback plan ready

This deployment guide provides a comprehensive overview of deploying the Cognize API server in various environments. Choose the approach that best fits your infrastructure requirements and organizational constraints.