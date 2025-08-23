# DevLog AI - Manual Configuration Guide

## Overview

This guide provides step-by-step instructions for setting up all required components of the DevLog AI platform. Follow these instructions to configure your development environment, cloud services, and deployment pipelines.

## Table of Contents

1. [Basic Setup](#1-basic-setup)
2. [Go Environment Configuration](#2-go-environment-configuration)
3. [Flutter Environment Configuration](#3-flutter-environment-configuration)
4. [Oracle Cloud Infrastructure (OCI) Setup](#4-oracle-cloud-infrastructure-oci-setup)
5. [Docker Hub Setup](#5-docker-hub-setup)
6. [Supabase Setup](#6-supabase-setup)
7. [Doppler Setup](#7-doppler-setup)
8. [Kafka & Redis (Upstash) Setup](#8-kafka--redis-upstash-setup)
9. [GitHub Actions Configuration](#9-github-actions-configuration)
10. [SSL/TLS Certificate Configuration](#10-ssltls-certificate-configuration)
11. [Kubernetes Configuration](#11-kubernetes-configuration)
12. [Vector Database Setup](#12-vector-database-setup)
13. [Embedding API Setup](#13-embedding-api-setup)
14. [Environment Variables Reference](#14-environment-variables-reference)

## 1. Basic Setup

### Clone the Repository

```bash
# Clone the repository
git clone https://github.com/your-organization/devlog-ai.git
cd devlog-ai

# Create TLS certificates
./scripts/generate-certs.sh
```

## 2. Go Environment Configuration

### Install Go

1. Download Go from [https://golang.org/dl/](https://golang.org/dl/)
2. Follow the installation instructions for your operating system
3. Verify installation:

```bash
go version
```

### Configure Go Modules

```bash
cd backend
go mod tidy
go mod download
```

### Run Tests

```bash
go test ./...
```

### Build Services

```bash
# Use the Makefile to build all services
make build

# Or build individual services
go build -o bin/api-service ./cmd/api-service
go build -o bin/ingestor-service ./cmd/ingestor-service
go build -o bin/processor-service ./cmd/processor-service
go build -o bin/esg-calculator ./cmd/esg-calculator
```

## 3. Flutter Environment Configuration

### Install Flutter

1. Follow the Flutter installation guide for your platform: [https://flutter.dev/docs/get-started/install](https://flutter.dev/docs/get-started/install)
2. Run Flutter doctor to ensure your setup is correct:

```bash
flutter doctor
```

### Configure Flutter Project

```bash
cd frontend
flutter pub get
```

### Generate JSON Serialization Code

```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

### Run Flutter Web App Locally

```bash
flutter run -d chrome
```

## 4. Oracle Cloud Infrastructure (OCI) Setup

### Create OCI Account

1. Go to [https://www.oracle.com/cloud/free/](https://www.oracle.com/cloud/free/)
2. Sign up for the Always Free tier
3. Verify your email and complete account setup

### Generate API Keys for OCI

1. Log in to OCI Console
2. Navigate to User Settings (click profile icon → User Settings)
3. Go to API Keys section
4. Click "Add API Key"
5. Choose "Generate API Key Pair"
6. Download the private key and save it securely
7. Click "Add" and note the key fingerprint

### Configure OCI CLI

1. Install OCI CLI following instructions at [https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm)
2. Run the setup wizard:

```bash
oci setup config
```

3. Enter the requested information:
   - User OCID (found in User Settings)
   - Tenancy OCID (found in Tenancy Details)
   - Region (e.g., us-phoenix-1)
   - Path to the private key file downloaded earlier
   - Fingerprint from the added key

### Create Compartment

1. In OCI Console, go to "Identity & Security" → "Compartments"
2. Click "Create Compartment"
3. Name it "DevLogAI"
4. Add a description and select the parent compartment
5. Click "Create Compartment"
6. Copy the Compartment OCID for later use in Terraform

### Deploy Infrastructure with Terraform

1. Navigate to the infrastructure directory:

```bash
cd infra/production
```

2. Initialize Terraform:

```bash
terraform init
```

3. Create terraform.tfvars file with required variables:

```
region           = "us-phoenix-1"
compartment_id   = "ocid1.compartment.oc1..aaaaaaaxxxxxxxx"
node_pool_image_id = "ocid1.image.oc1.phx.aaaaaaaxxxxxxxx"
db_password      = "securePassword123"
```

4. Apply the Terraform configuration:

```bash
terraform plan
terraform apply
```

5. After successful deployment, Terraform will output important information like the Kubernetes cluster endpoint

## 5. Docker Hub Setup

### Create Docker Hub Account

1. Go to [https://hub.docker.com/](https://hub.docker.com/)
2. Sign up for an account
3. Verify your email address

### Create Access Token

1. Log in to Docker Hub
2. Go to Account Settings → Security
3. Click "New Access Token"
4. Name it "DevLogAI-CI"
5. Set appropriate permissions (Read, Write, Delete)
6. Click "Generate"
7. Save the token securely (you'll need it for GitHub Actions)

### Create Repositories

Create repositories for each service:

1. Click "Create Repository"
2. Name it "api-service" (repeat for "ingestor-service", "processor-service", etc.)
3. Set visibility to public or private as needed
4. Click "Create"

### Configure Docker Hub Secrets

Add your Docker Hub credentials to both Doppler and GitHub Actions:

```bash
# Add to Doppler (see Doppler setup section below)
doppler secrets set DOCKERHUB_USERNAME=your_username
doppler secrets set DOCKERHUB_TOKEN=your_token

# Add to GitHub (see GitHub Actions section below)
```

## 6. Supabase Setup

### Create Supabase Account

1. Go to [https://supabase.com/](https://supabase.com/)
2. Sign up for an account
3. Verify your email address

### Create Project

1. Click "New Project"
2. Choose an organization or create a new one
3. Enter project details:
   - Name: "DevLogAI"
   - Database Password: (create a secure password)
   - Region: (choose a region close to your users)
4. Click "Create New Project"
5. Wait for project to initialize

### Configure Database Schema

1. Go to the SQL Editor
2. Create the necessary tables:

```sql
-- Create users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email TEXT UNIQUE NOT NULL,
  name TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create projects table
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  owner_id UUID REFERENCES users(id) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create logs table for storing log metadata
CREATE TABLE logs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  project_id UUID REFERENCES projects(id) NOT NULL,
  log_id TEXT NOT NULL,
  timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
  level TEXT NOT NULL,
  message TEXT NOT NULL,
  source TEXT NOT NULL,
  service_name TEXT NOT NULL,
  metadata JSONB,
  vector_id TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
```

### Get Connection Information

1. Go to Project Settings → Database
2. Copy the connection string
3. Store it for use with Doppler

## 7. Doppler Setup

### Create Doppler Account

1. Go to [https://www.doppler.com/](https://www.doppler.com/)
2. Sign up for an account
3. Verify your email address

### Create Project & Environment

1. Click "Create Project"
2. Name it "DevLogAI"
3. Select "Development" as initial environment
4. Click "Create Project"

### Configure Secrets

Using the Doppler web interface, add the following secrets:

#### OCI Secrets

- `OCI_TENANCY_OCID`: Your OCI tenancy OCID
- `OCI_USER_OCID`: Your OCI user OCID
- `OCI_COMPARTMENT_ID`: Your OCI compartment ID
- `OCI_REGION`: Your OCI region

#### Docker Hub Secrets

- `DOCKERHUB_USERNAME`: Your Docker Hub username
- `DOCKERHUB_TOKEN`: Your Docker Hub access token

#### Database Secrets

- `SUPABASE_URL`: Your Supabase project URL
- `SUPABASE_KEY`: Your Supabase project API key
- `DATABASE_URL`: Your Supabase connection string

#### Kafka Secrets

- `USE_REAL_KAFKA`: Set to "true" for production
- `KAFKA_BROKERS`: Comma-separated list of Kafka brokers (e.g., "kafka-1.upstash.io:9092,kafka-2.upstash.io:9092")
- `KAFKA_USERNAME`: Your Upstash Kafka username
- `KAFKA_PASSWORD`: Your Upstash Kafka password
- `KAFKA_TOPIC`: Topic name for logs (default: "logs")
- `KAFKA_GROUP_ID`: Consumer group ID (default: "devlogai-processor")

#### Redis Secrets

- `REDIS_URL`: Your Upstash Redis connection string

#### Vector DB Secrets

- `VECTORDB_TYPE`: "pinecone" or "weaviate"
- `PINECONE_API_KEY`: Your Pinecone API key
- `PINECONE_ENVIRONMENT`: Your Pinecone environment
- `PINECONE_INDEX_NAME`: Name of your Pinecone index (default: "devlogai-logs")
- `WEAVIATE_API_KEY`: Your Weaviate API key
- `WEAVIATE_URL`: Your Weaviate URL

#### Embedding API Secrets

- `EMBEDDING_TYPE`: "openai" or "gemini"
- `OPENAI_API_KEY`: Your OpenAI API key
- `OPENAI_EMBEDDING_MODEL`: Embedding model name (default: "text-embedding-ada-002")
- `GEMINI_API_KEY`: Your Google Gemini API key
- `GEMINI_EMBEDDING_MODEL`: Embedding model name (default: "embedding-001")

#### TLS/Security Secrets

- `TLS_CERT_PATH`: Path to TLS certificate (or leave empty for default path)
- `TLS_KEY_PATH`: Path to TLS key (or leave empty for default path)
- `ENABLE_PII_ANONYMIZATION`: Set to "true" to enable PII detection and anonymization (default: true)

### Install Doppler CLI

1. Install the Doppler CLI following instructions at [https://docs.doppler.com/docs/cli](https://docs.doppler.com/docs/cli)
2. Login to your Doppler account:

```bash
doppler login
```

3. Configure your project and environment:

```bash
doppler setup --project devlogai --config development
```

### Using Doppler with CLI Commands

To run services with Doppler environment variables:

```bash
doppler run -- go run ./cmd/api-service/main.go
```

### Get Service Token

1. Go to Access → Service Tokens
2. Click "Generate"
3. Name it "GitHub Actions"
4. Select appropriate permissions
5. Click "Generate"
6. Copy the token and save it (you'll add it as a GitHub secret)

## 8. Kafka & Redis (Upstash) Setup

### Create Upstash Account

1. Go to [https://upstash.com/](https://upstash.com/)
2. Sign up for an account
3. Verify your email address

### Create Kafka Cluster

1. Go to the Kafka section
2. Click "Create Cluster"
3. Configure:
   - Name: "devlogai-kafka"
   - Region: (choose one close to your OCI region)
   - Type: Single (Free)
4. Click "Create"

### Create Kafka Topics

1. Select your cluster
2. Go to Topics tab
3. Create topics:
   - "logs" (for log ingestion)
   - "processed-logs" (for processed logs)

### Create Redis Database

1. Go to the Redis section
2. Click "Create Database"
3. Configure:
   - Name: "devlogai-redis"
   - Region: (choose one close to your OCI region)
   - Type: Free
4. Click "Create"

### Get Connection Details

1. For Kafka, copy the endpoints and credentials
2. For Redis, copy the connection string
3. Add these to Doppler:
   - `KAFKA_BROKERS`
   - `KAFKA_USERNAME`
   - `KAFKA_PASSWORD`
   - `REDIS_URL`

## 9. GitHub Actions Configuration

### Add Repository Secrets

In your GitHub repository:

1. Go to Settings → Secrets and variables → Actions
2. Add the following secrets:
   - `DOPPLER_TOKEN`: Your Doppler service token
   - `DOCKERHUB_USERNAME`: Your Docker Hub username
   - `DOCKERHUB_TOKEN`: Your Docker Hub access token
   - `OCI_CLI_FINGERPRINT`: Your OCI API key fingerprint
   - `OCI_CLI_KEY_CONTENT`: Your OCI private key (entire contents)
   - `OCI_CLI_REGION`: Your OCI region
   - `OCI_CLI_TENANCY`: Your OCI tenancy OCID
   - `OCI_CLI_USER`: Your OCI user OCID

### Add Secrets via CLI

You can also add secrets using the GitHub CLI:

```bash
# Install GitHub CLI if you haven't already
# https://cli.github.com/

# Log in to GitHub
gh auth login

# Set secrets
gh secret set DOCKERHUB_USERNAME --body "your_dockerhub_username"
gh secret set DOCKERHUB_TOKEN --body "your_dockerhub_token"
gh secret set DOPPLER_TOKEN --body "your_doppler_token"
```

### Update Workflow File

Make sure your CI/CD workflow file (.github/workflows/ci-cd.yaml) is properly configured:

1. Ensure the Docker Hub secrets are referenced correctly:

```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v2
  with:
    username: ${{ secrets.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}
```

2. Make sure the Doppler token is used for retrieving secrets:

```yaml
- name: Set up Doppler
  uses: dopplerhq/cli-action@v1
  with:
    token: ${{ secrets.DOPPLER_TOKEN }}
```

## 10. SSL/TLS Certificate Configuration

### Development Certificates

For development, we've already created self-signed certificates with:

```bash
./scripts/generate-certs.sh
```

### Production Certificates

For production, you should obtain certificates from a trusted certificate authority. Options include:

1. **Let's Encrypt (free):**

   - Use [certbot](https://certbot.eff.org/) to obtain certificates
   - Place them in the `certs` directory or configure paths in environment variables

2. **OCI Certificate Service:**

   - Go to OCI Console → Security → Certificates
   - Create a certificate
   - Download and place in appropriate location

3. **Kubernetes TLS:**
   - For Kubernetes deployments, consider using cert-manager:

```bash
# Install cert-manager using Helm
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.13.0 \
  --set installCRDs=true

# Create an Issuer or ClusterIssuer
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF

# Then use it in your Ingress resources
```

## 11. Kubernetes Configuration

### Install kubectl

Follow the instructions at [https://kubernetes.io/docs/tasks/tools/install-kubectl/](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Connect to OKE Cluster

After creating the cluster with Terraform:

```bash
oci ce cluster create-kubeconfig \
    --cluster-id <cluster-id> \
    --file ~/.kube/config \
    --region <region> \
    --token-version 2.0.0
```

### Deploy Services

1. Update the environment variables in the K8s manifests
2. Apply the manifests:

```bash
kubectl apply -f k8s/base/api-deployment.yaml
kubectl apply -f k8s/base/ingestor-deployment.yaml
kubectl apply -f k8s/base/processor-deployment.yaml
kubectl apply -f k8s/base/esg-calculator-deployment.yaml
```

### Configure Kubernetes Secrets

Use Doppler to sync secrets to Kubernetes:

```bash
doppler run -- \
  kubectl create secret generic devlogai-secrets \
  --from-literal=KAFKA_BROKERS="$KAFKA_BROKERS" \
  --from-literal=KAFKA_USERNAME="$KAFKA_USERNAME" \
  --from-literal=KAFKA_PASSWORD="$KAFKA_PASSWORD" \
  --from-literal=REDIS_URL="$REDIS_URL" \
  --from-literal=OPENAI_API_KEY="$OPENAI_API_KEY" \
  --from-literal=SUPABASE_URL="$SUPABASE_URL" \
  --from-literal=SUPABASE_KEY="$SUPABASE_KEY"
```

## 12. Vector Database Setup

### Pinecone Setup

1. Go to [https://www.pinecone.io/](https://www.pinecone.io/)
2. Sign up for an account
3. Create an index:
   - Name: "devlogai-logs"
   - Dimensions: 1536 (for OpenAI embeddings)
   - Metric: "cosine"
4. Copy API key and environment
5. Add to Doppler:
   - `PINECONE_API_KEY`
   - `PINECONE_ENVIRONMENT`
   - `PINECONE_INDEX_NAME`

### Weaviate Setup (Alternative)

1. Go to [https://weaviate.io/](https://weaviate.io/)
2. Sign up for Weaviate Cloud
3. Create a cluster:
   - Name: "devlogai"
   - Choose a region close to your OCI region
4. Copy API key and URL
5. Add to Doppler:
   - `WEAVIATE_URL`
   - `WEAVIATE_API_KEY`

## 13. Embedding API Setup

### OpenAI Setup

1. Go to [https://platform.openai.com/](https://platform.openai.com/)
2. Sign up for an account
3. Navigate to "API Keys" and create a new secret key
4. Copy the API key and add to Doppler as `OPENAI_API_KEY`

### Google Gemini Setup (Alternative)

1. Go to [Google AI Studio](https://makersuite.google.com/)
2. Sign up and get access to Gemini API
3. Create an API key from the console
4. Add to Doppler as `GEMINI_API_KEY`

## 14. Environment Variables Reference

### Core Configuration

- `PORT`: Port for the service to listen on (default varies by service)
- `DEBUG`: Enable debug logging (true/false)
- `USE_REAL_KAFKA`: Use real Kafka instead of mock (true/false)

### Security Configuration

- `ENABLE_PII_ANONYMIZATION`: Enable PII detection and anonymization (true/false)
- `TLS_CERT_PATH`: Path to TLS certificate
- `TLS_KEY_PATH`: Path to TLS private key

### Kafka Configuration

- `KAFKA_BROKERS`: Comma-separated list of Kafka brokers
- `KAFKA_USERNAME`: Kafka username for authentication
- `KAFKA_PASSWORD`: Kafka password for authentication
- `KAFKA_TOPIC`: Kafka topic name for logs
- `KAFKA_GROUP_ID`: Consumer group ID

### Database Configuration

- `DATABASE_URL`: PostgreSQL connection string
- `SUPABASE_URL`: Supabase project URL
- `SUPABASE_KEY`: Supabase API key

### Redis Configuration

- `REDIS_URL`: Redis connection string

### Vector Database Configuration

- `VECTORDB_TYPE`: Vector database type (pinecone/weaviate)
- `VECTORDB_INDEX_NAME`: Name of the vector database index
- `PINECONE_API_KEY`: Pinecone API key
- `PINECONE_ENVIRONMENT`: Pinecone environment
- `PINECONE_NAMESPACE`: Pinecone namespace (optional)
- `WEAVIATE_API_KEY`: Weaviate API key
- `WEAVIATE_URL`: Weaviate URL

### Embedding Configuration

- `EMBEDDING_TYPE`: Embedding provider (openai/gemini)
- `EMBEDDING_BATCH_SIZE`: Number of texts to embed in one batch
- `OPENAI_API_KEY`: OpenAI API key
- `OPENAI_EMBEDDING_MODEL`: OpenAI embedding model name
- `GEMINI_API_KEY`: Google Gemini API key
- `GEMINI_EMBEDDING_MODEL`: Gemini embedding model name

### OCI Configuration

- `OCI_REGION`: OCI region
- `OCI_TENANCY_OCID`: OCI tenancy OCID
- `OCI_USER_OCID`: OCI user OCID
- `OCI_COMPARTMENT_ID`: OCI compartment ID
