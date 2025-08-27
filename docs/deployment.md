# Deployment Guide

This guide covers deploying the Cognize server to various production environments.

## Quick Deploy Options

### 1. Docker Compose (Recommended for VPS)

The easiest way to deploy to a VPS or dedicated server.

```bash
# Clone the repository
git clone https://github.com/Cognize-AI/server-cognize.git
cd server-cognize

# Create production environment file
cp .env.example .env.prod
# Edit .env.prod with production values

# Deploy with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f
```

### 2. Railway

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template)

1. Click the Railway button above
2. Connect your GitHub account
3. Configure environment variables
4. Deploy automatically

### 3. Render

1. Connect your GitHub repository to Render
2. Choose "Web Service"
3. Set build command: `go build -o main .`
4. Set start command: `./main`
5. Configure environment variables
6. Deploy

### 4. DigitalOcean App Platform

1. Create a new app from your GitHub repository
2. Configure the app spec (see App Platform section below)
3. Set environment variables
4. Deploy

## Detailed Deployment Instructions

### Production Environment Variables

Create a production `.env` file with the following variables:

```bash
# Server Configuration
PORT=4000
ENVIRONMENT=prod

# Database Configuration (use your production database URL)
DB_STRING=postgres://username:password@your-db-host:5432/cognize?sslmode=require

# Security
JWT_SECRET=your-very-long-and-secure-jwt-secret-for-production
ENC_SECRET=your-32-character-encryption-key-prod

# Google OAuth (production credentials)
GOOGLE_OAUTH_CLIENT_ID=your-production-google-oauth-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-production-google-oauth-client-secret
GOOGLE_OAUTH_REDIRECT_URL=https://your-frontend-domain.com/auth/callback

# Logging (recommended for production)
AXIOM_TOKEN=your-axiom-token
AXIOM_ORG=your-axiom-org
AXIOM_DATASET=cognize-prod

# Optional: Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-production-email@yourdomain.com
SMTP_PASSWORD=your-app-password
```

### Docker Deployment

#### Production Docker Compose

Create a `docker-compose.prod.yml` file:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:17
    container_name: cognize_postgres_prod
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    networks:
      - cognize_network
    # Don't expose ports externally in production
    # ports:
    #   - "5432:5432"

  backend:
    build: .
    container_name: cognize_backend_prod
    restart: unless-stopped
    ports:
      - "80:4000"  # Or use a reverse proxy
    env_file:
      - .env.prod
    depends_on:
      - postgres
    networks:
      - cognize_network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:4000/"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Optional: Add nginx reverse proxy
  nginx:
    image: nginx:alpine
    container_name: cognize_nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - backend
    networks:
      - cognize_network

volumes:
  postgres_data:

networks:
  cognize_network:
    driver: bridge
```

#### Deploy with Production Compose

```bash
# Build and start services
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Update deployment
git pull
docker-compose -f docker-compose.prod.yml up -d --build

# Backup database
docker exec cognize_postgres_prod pg_dump -U $POSTGRES_USER cognize > backup_$(date +%Y%m%d_%H%M%S).sql
```

### VPS Deployment

#### Prerequisites

- Ubuntu 20.04+ or CentOS 8+ server
- Docker and Docker Compose installed
- Domain name pointing to your server
- SSL certificate (Let's Encrypt recommended)

#### Setup Steps

1. **Server Setup**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   sudo usermod -aG docker $USER
   
   # Install Docker Compose
   sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   
   # Install nginx (if not using Docker nginx)
   sudo apt install nginx certbot python3-certbot-nginx -y
   ```

2. **Clone and Configure**
   ```bash
   git clone https://github.com/Cognize-AI/server-cognize.git
   cd server-cognize
   
   # Create production environment
   cp .env.example .env.prod
   nano .env.prod  # Edit with production values
   ```

3. **SSL Certificate**
   ```bash
   # Get SSL certificate
   sudo certbot --nginx -d your-domain.com
   
   # Auto-renewal
   sudo crontab -e
   # Add: 0 12 * * * /usr/bin/certbot renew --quiet
   ```

4. **Deploy**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Cloud Platform Deployments

#### Railway

1. Connect GitHub repository
2. Set environment variables in Railway dashboard
3. Configure custom domain
4. Deploy automatically on push

**Railway App Configuration:**
```json
{
  "build": {
    "builder": "NIXPACKS"
  },
  "deploy": {
    "startCommand": "go run main.go",
    "restartPolicyType": "ON_FAILURE"
  }
}
```

#### Render

**render.yaml:**
```yaml
services:
  - type: web
    name: cognize-server
    env: go
    buildCommand: go build -o main .
    startCommand: ./main
    envVars:
      - key: PORT
        value: 4000
      - key: ENVIRONMENT
        value: prod
      # Add other environment variables in Render dashboard
    healthCheckPath: /
```

#### DigitalOcean App Platform

**app.yaml:**
```yaml
name: cognize-server
services:
  - name: api
    source_dir: /
    github:
      repo: Cognize-AI/server-cognize
      branch: main
    run_command: ./main
    build_command: go build -o main .
    environment_slug: go
    instance_count: 1
    instance_size_slug: basic-xxs
    routes:
      - path: /
    health_check:
      http_path: /
    envs:
      - key: PORT
        value: "4000"
      - key: ENVIRONMENT
        value: "prod"
      # Add other environment variables

databases:
  - name: cognize-db
    engine: PG
    version: "17"
    size_slug: db-s-1vcpu-1gb
```

#### Google Cloud Run

```bash
# Build and push image
docker build -t gcr.io/YOUR_PROJECT_ID/cognize-server .
docker push gcr.io/YOUR_PROJECT_ID/cognize-server

# Deploy to Cloud Run
gcloud run deploy cognize-server \
  --image gcr.io/YOUR_PROJECT_ID/cognize-server \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars PORT=8080,ENVIRONMENT=prod
```

#### AWS ECS

**task-definition.json:**
```json
{
  "family": "cognize-server",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "cognize-server",
      "image": "YOUR_ACCOUNT.dkr.ecr.REGION.amazonaws.com/cognize-server:latest",
      "portMappings": [
        {
          "containerPort": 4000,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {"name": "PORT", "value": "4000"},
        {"name": "ENVIRONMENT", "value": "prod"}
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/cognize-server",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### Database Setup

#### Managed Database Services

**PostgreSQL Options:**
- **AWS RDS**: Fully managed PostgreSQL
- **Google Cloud SQL**: Managed PostgreSQL on GCP
- **DigitalOcean Managed Databases**: Simple managed PostgreSQL
- **Supabase**: PostgreSQL with additional features
- **Neon**: Serverless PostgreSQL

**Connection String Format:**
```bash
# SSL required for most managed services
DB_STRING=postgres://username:password@host:port/database?sslmode=require
```

#### Self-Hosted Database

```bash
# Using Docker for production database
docker run --name cognize-postgres-prod \
  -e POSTGRES_USER=cognize \
  -e POSTGRES_PASSWORD=secure_password \
  -e POSTGRES_DB=cognize \
  -v /var/lib/postgresql/data:/var/lib/postgresql/data \
  -v /backups:/backups \
  --restart unless-stopped \
  -p 5432:5432 \
  postgres:17
```

### Monitoring and Logging

#### Health Checks

The application provides a health check endpoint:
```http
GET /
Response: "Welcome to cognize"
```

#### Logging

Production logging is handled by Axiom. Configure these environment variables:
```bash
AXIOM_TOKEN=your-token
AXIOM_ORG=your-org
AXIOM_DATASET=cognize-prod
```

#### Monitoring Tools

- **Uptime Monitoring**: UptimeRobot, Pingdom
- **Application Monitoring**: New Relic, Datadog
- **Error Tracking**: Sentry
- **Metrics**: Prometheus + Grafana

### Backup and Recovery

#### Database Backups

```bash
# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec cognize-postgres-prod pg_dump -U cognize cognize > /backups/cognize_$DATE.sql
gzip /backups/cognize_$DATE.sql

# Keep only last 30 days of backups
find /backups -name "cognize_*.sql.gz" -mtime +30 -delete
```

#### Cron Job for Backups

```bash
# Add to crontab
0 2 * * * /path/to/backup-script.sh
```

### Security Considerations

1. **Environment Variables**: Never commit secrets to Git
2. **SSL/TLS**: Always use HTTPS in production
3. **Database**: Use strong passwords and SSL connections
4. **CORS**: Configure for your specific frontend domain
5. **API Keys**: Rotate regularly and use strong encryption
6. **Updates**: Keep dependencies and base images updated

### Performance Optimization

1. **Database Connection Pooling**: Configured in GORM
2. **Caching**: Consider Redis for session storage
3. **CDN**: Use CloudFlare or AWS CloudFront for static assets
4. **Load Balancing**: Use nginx or cloud load balancers for high traffic
5. **Horizontal Scaling**: Deploy multiple instances behind a load balancer

### Troubleshooting

#### Common Production Issues

1. **Database Connection Timeout**
   - Check connection string format
   - Verify database is accessible
   - Check network security groups

2. **OAuth Redirect Issues**
   - Verify redirect URL matches Google OAuth config
   - Check HTTPS vs HTTP in URLs

3. **Environment Variable Issues**
   - Use `docker exec -it container env` to check variables
   - Verify .env file is being loaded correctly

4. **Performance Issues**
   - Monitor database query performance
   - Check memory and CPU usage
   - Review logs for errors

#### Logs and Debugging

```bash
# View application logs
docker logs cognize_backend_prod -f

# Check database logs
docker logs cognize_postgres_prod -f

# Access application container
docker exec -it cognize_backend_prod /bin/sh
```

---

For more deployment options or specific platform questions, please check our [GitHub issues](https://github.com/Cognize-AI/server-cognize/issues) or create a new issue.