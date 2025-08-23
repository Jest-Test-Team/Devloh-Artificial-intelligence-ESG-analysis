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

## 4. Oracle Cloud Infrastructure (OCI) Setup

### Create OCI Account

1. Go to [https://www.oracle.com/cloud/free/](https://www.oracle.com/cloud/free/)
2. Sign up for the Always Free tier

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
4. Add a description
5. Select the parent compartment
6. Click "Create Compartment"
7. Copy the Compartment OCID for later use

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
5. Set appropriate permissions
6. Click "Generate"
7. Save the token securely (you'll need it for GitHub Actions)

### Create Repositories

Create repositories for each service:

1. Click "Create Repository"
2. Name it "api-service" (repeat for "ingestor-service", "processor-service", etc.)
3. Set visibility to public or private as needed
4. Click "Create"

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
```

### Get Connection Information

1. Go to Project Settings → Database
2. Copy the connection string
3. Save it for use with Doppler

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

Add the following secrets:

1. `SUPABASE_URL`: Your Supabase project URL
2. `SUPABASE_KEY`: Your Supabase project API key
3. `DATABASE_URL`: Your Supabase connection string
4. `OCI_TENANCY_OCID`: Your OCI tenancy OCID
5. `OCI_USER_OCID`: Your OCI user OCID
6. `OCI_COMPARTMENT_ID`: Your OCI compartment ID
7. `OCI_REGION`: Your OCI region
8. `DOCKERHUB_USERNAME`: Your Docker Hub username
9. `DOCKERHUB_TOKEN`: Your Docker Hub access token
10. `TLS_CERT_PATH`: Path to TLS certificate (or leave empty for default)
11. `TLS_KEY_PATH`: Path to TLS key (or leave empty for default)

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
   - "logs"
   - "processed-logs"

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
   - For Kubernetes deployments, consider using cert-manager

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

### Weaviate Setup (Alternative)

1. Go to [https://weaviate.io/](https://weaviate.io/)
2. Sign up for Weaviate Cloud
3. Create a cluster
4. Copy API key and URL
5. Add to Doppler:
   - `WEAVIATE_URL`
   - `WEAVIATE_API_KEY`
