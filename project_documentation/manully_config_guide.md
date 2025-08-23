# DevLog AI - Manual Configuration Guide

This guide provides step-by-step instructions for configuring the external services and environments required by DevLog AI.

## Table of Contents

1. [Oracle Cloud Infrastructure (OCI) Setup](#1-oracle-cloud-infrastructure-oci-setup)
2. [Docker Hub Setup](#2-docker-hub-setup)
3. [Supabase Setup](#3-supabase-setup)
4. [Doppler Setup](#4-doppler-setup)
5. [Snyk Setup](#5-snyk-setup)
6. [GitHub Secrets Configuration](#6-github-secrets-configuration)
7. [Upstash (Kafka & Redis) Setup](#7-upstash-kafka--redis-setup)
8. [Vector Database Setup](#8-vector-database-setup)
9. [LLM Provider Setup](#9-llm-provider-setup)

---

## 1. Oracle Cloud Infrastructure (OCI) Setup

### Sign Up for OCI Free Tier

1. Go to [Oracle Cloud Free Tier](https://www.oracle.com/cloud/free/)
2. Click "Start for free"
3. Fill out the registration form and follow the verification steps
4. Complete the account setup process

### Generate API Keys

1. Log in to your Oracle Cloud account
2. Click on the profile icon in the top right and select "User Settings"
3. Under "Resources", click on "API Keys"
4. Click "Add API Key"
5. Select "Generate API Key Pair" and click "Download Private Key"
6. Save the private key securely and click "Add"
7. Note down the fingerprint displayed

### Configure OCI CLI

1. Install OCI CLI on your local machine following the [official documentation](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm)
2. Run the CLI setup wizard:
   ```bash
   oci setup config
   ```
3. Enter the required information:
   - User OCID (found in User Settings)
   - Tenancy OCID (found in Tenancy Details)
   - Region (e.g., us-phoenix-1)
   - Path to the private key file you downloaded
   - The key's fingerprint from earlier

### Create Compartment

1. In the OCI Console, navigate to "Identity & Security" > "Compartments"
2. Click "Create Compartment"
3. Name it "DevLogAI" (or your preferred name)
4. Add a description and select the parent compartment
5. Click "Create Compartment"
6. Note down the Compartment OCID for later use in Terraform

### Configure VCN Network Security

1. Navigate to "Networking" > "Virtual Cloud Networks"
2. Use the Terraform scripts in the project to create your VCN, or create it manually:
   - Click "Start VCN Wizard" and select "VCN with Internet Connectivity"
   - Configure CIDR blocks and subnet settings as needed
   - Ensure the security lists allow necessary ports (80, 443, 6443 for Kubernetes)

---

## 2. Docker Hub Setup

### Create Docker Hub Account

1. Go to [Docker Hub](https://hub.docker.com/)
2. Click "Sign Up" and create an account
3. Verify your email address

### Create Access Token

1. Log in to your Docker Hub account
2. Click on your username in the top right, then "Account Settings"
3. In the left menu, go to "Security"
4. Click "New Access Token"
5. Give it a name (e.g., "DevLogAI-CI")
6. Select appropriate permissions (Read, Write, Delete)
7. Click "Generate" and copy the token
8. Store this token securely; you'll need it for GitHub Actions

### Create Repository (Optional)

1. On your Docker Hub dashboard, click "Create Repository"
2. Name your repository (e.g., "devlog-ai")
3. Add a description and set visibility (Public/Private)
4. Click "Create"

---

## 3. Supabase Setup

### Create Supabase Account

1. Go to [Supabase](https://supabase.com/)
2. Click "Start your project" and sign up
3. Verify your email address

### Create a New Project

1. After logging in, click "New Project"
2. Select an organization or create a new one
3. Enter a name for your project (e.g., "DevLogAI")
4. Choose a database password (save this securely)
5. Select your region (choose one close to your OCI region)
6. Click "Create new project"

### Set Up Database Schema

1. In your project dashboard, go to the "SQL Editor"
2. Create tables for users and projects:

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

1. Go to "Settings" > "Database"
2. Under "Connection string", select "URI" and copy the connection string
3. Store this connection string securely; you'll need it for Doppler

---

## 4. Doppler Setup

### Create Doppler Account

1. Go to [Doppler](https://www.doppler.com/)
2. Click "Get Started" and sign up
3. Verify your email address

### Create Project

1. After logging in, click "Create Project"
2. Name your project (e.g., "DevLogAI")
3. Select "Development" as the initial environment
4. Click "Create Project"

### Configure Secrets

1. In your project dashboard, add the following secrets:
   - `SUPABASE_URL`: Your Supabase project URL
   - `SUPABASE_KEY`: Your Supabase project API key
   - `DATABASE_URL`: Your Supabase connection string
   - `OCI_TENANCY_OCID`: Your OCI tenancy OCID
   - `OCI_USER_OCID`: Your OCI user OCID
   - `OCI_COMPARTMENT_ID`: Your OCI compartment ID
   - `OCI_REGION`: Your OCI region
   - `DOCKERHUB_USERNAME`: Your Docker Hub username
   - `DOCKERHUB_TOKEN`: Your Docker Hub access token

### Get Doppler Token

1. Go to "Access" > "Service Tokens"
2. Click "Generate" and create a token for your CI/CD
3. Copy the token and store it securely
4. You'll add this as a GitHub secret later

---

## 5. Snyk Setup

### Create Snyk Account

1. Go to [Snyk](https://snyk.io/)
2. Click "Sign Up" and create an account
3. Verify your email address

### Connect to GitHub

1. After logging in, click "Add projects"
2. Select "GitHub" and authorize Snyk
3. Choose the DevLog AI repository
4. Click "Add selected repositories"

### Generate API Token

1. Go to your account settings by clicking on your profile picture
2. Select "General" and scroll to "API Token"
3. Click "Show" and copy your token
4. Store this token securely; you'll add it as a GitHub secret

---

## 6. GitHub Secrets Configuration

### Configure GitHub Repository Secrets

1. Go to your GitHub repository
2. Click on "Settings" > "Secrets and variables" > "Actions"
3. Add the following secrets:
   - `DOPPLER_TOKEN`: Your Doppler service token
   - `DOCKERHUB_USERNAME`: Your Docker Hub username
   - `DOCKERHUB_TOKEN`: Your Docker Hub access token
   - `SNYK_TOKEN`: Your Snyk API token
   - `OCI_CLI_FINGERPRINT`: Your OCI API key fingerprint
   - `OCI_CLI_KEY_CONTENT`: Your OCI private key (entire contents)
   - `OCI_CLI_REGION`: Your OCI region
   - `OCI_CLI_TENANCY`: Your OCI tenancy OCID
   - `OCI_CLI_USER`: Your OCI user OCID

---

## 7. Upstash (Kafka & Redis) Setup

### Create Upstash Account

1. Go to [Upstash](https://upstash.com/)
2. Click "Sign Up" and create an account
3. Verify your email address

### Create Kafka Cluster

1. After logging in, go to the Kafka tab
2. Click "Create Cluster"
3. Configure your cluster:
   - Name: "devlogai-kafka"
   - Region: Choose one close to your OCI region
   - Plan: Select "Free" for development
4. Click "Create"

### Create Kafka Topics

1. Select your new cluster
2. Go to "Topics" and click "Create Topic"
3. Create the following topics:
   - "logs" (for log ingestion)
   - "processed-logs" (for processed logs)
4. Configure partitions and retention as needed

### Create Redis Database

1. Go to the Redis tab
2. Click "Create Database"
3. Configure your database:
   - Name: "devlogai-redis"
   - Region: Choose one close to your OCI region
   - Plan: Select "Free" for development
4. Click "Create"

### Get Connection Details

1. For Kafka:
   - Select your cluster and go to "Details"
   - Copy the endpoints and credentials
2. For Redis:
   - Select your database and go to "Details"
   - Copy the connection string and REST URL
3. Add these credentials to Doppler as:
   - `KAFKA_BROKERS`
   - `KAFKA_USERNAME`
   - `KAFKA_PASSWORD`
   - `REDIS_URL`

---

## 8. Vector Database Setup

### Option 1: Pinecone Setup

1. Go to [Pinecone](https://www.pinecone.io/)
2. Sign up for an account and verify your email
3. Create a new index:
   - Name: "devlogai-logs"
   - Dimensions: 1536 (for OpenAI embeddings)
   - Metric: "cosine"
4. Copy the API key and environment from the dashboard
5. Add to Doppler as `PINECONE_API_KEY` and `PINECONE_ENVIRONMENT`

### Option 2: Weaviate Setup

1. Go to [Weaviate Cloud](https://weaviate.io/pricing)
2. Sign up for the free Sandbox tier
3. Create a new cluster:
   - Name: "devlogai"
   - Choose a region close to your OCI region
4. Once created, get the cluster URL and API key
5. Add to Doppler as `WEAVIATE_URL` and `WEAVIATE_API_KEY`

---

## 9. LLM Provider Setup

### OpenAI Setup

1. Go to [OpenAI](https://platform.openai.com/)
2. Sign up for an account if you don't have one
3. Navigate to "API Keys" and create a new secret key
4. Copy the API key and add to Doppler as `OPENAI_API_KEY`

### Google Gemini Setup (Alternative)

1. Go to [Google AI Studio](https://makersuite.google.com/)
2. Sign up and get access to Gemini API
3. Create an API key from the console
4. Add to Doppler as `GEMINI_API_KEY`

---

## Deployment Process

After configuring all the above services, you can deploy DevLog AI by following these steps:

1. **Initialize Infrastructure:**

   ```bash
   cd infra/production
   terraform init
   terraform plan
   terraform apply
   ```

2. **Update Kubernetes Config:**

   ```bash
   # Get Kubeconfig from OCI
   oci ce cluster create-kubeconfig --cluster-id <OKE_CLUSTER_ID> --file ~/.kube/config --region <REGION> --token-version 2.0.0
   ```

3. **Deploy Services to Kubernetes:**

   ```bash
   kubectl apply -f k8s/base/api-deployment.yaml
   kubectl apply -f k8s/base/ingestor-deployment.yaml
   ```

4. **Verify Deployment:**
   ```bash
   kubectl get pods
   kubectl get services
   ```

## Troubleshooting

### Common Issues and Solutions

1. **OCI Authentication Errors:**

   - Verify your API key and configuration in `~/.oci/config`
   - Check that the key file permissions are set to 600

2. **Docker Hub Authentication Failures:**

   - Ensure your Docker Hub token has the correct permissions
   - Try logging in manually with `docker login` to verify credentials

3. **Kubernetes Connectivity Issues:**

   - Check that your kubeconfig is correctly set up
   - Verify network security rules allow traffic to/from the cluster

4. **Database Connection Problems:**

   - Confirm your Supabase connection string is correct
   - Check for any IP restrictions on your database

5. **Kafka/Redis Connectivity:**
   - Verify your Upstash credentials
   - Check that your services are trying to connect to the correct endpoints
