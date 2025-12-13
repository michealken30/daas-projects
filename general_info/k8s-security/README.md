# Kubernetes Security: Loan Application Use Case

## Table of Contents

1. [Introduction](#introduction)
2. [Prerequisites & Setup](#prerequisites--setup)
3. [The 4 C's of Cloud Native Security](#the-4-cs-of-cloud-native-security)
4. [RBAC (Role-Based Access Control)](#rbac-role-based-access-control)
5. [Pod Security Standards](#pod-security-standards)
6. [Container Image Scanning](#container-image-scanning)
7. [Policy Management](#policy-management)

---

## Introduction

This guide teaches Kubernetes security using a real-world loan application deployed on minikube. We'll cover security from code to cloud, implementing best practices and using industry-standard tools.

### Learning Objectives

By the end of this course, you will:
- Understand the 4 C's of cloud native security
- Implement RBAC for fine-grained access control
- Apply Pod Security Standards to secure workloads
- Configure and analyze Kubernetes audit logs
- Scan container images for vulnerabilities
- Enforce policies using OPA/Kyverno

### Application Overview

Our loan application consists of:
- **PostgreSQL Database**: Stores loan and user data
- **Government API**: Internal API for loan processing
- **Loan Validator Portal**: Public-facing web application
- **OTel Collector**: Observability component

---


## Initial Security Assessment

Before we secure the application, let's identify security issues:

```bash
# Check current pod security contexts
kubectl get pods -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.securityContext}{"\n"}{end}'

# Check service accounts
kubectl get serviceaccounts

# Check RBAC
kubectl get roles,rolebindings,clusterroles,clusterrolebindings

# Check for secrets in plain text
kubectl get secrets -o yaml
```

**Security Issues Identified:**
- No explicit security contexts
- Default service accounts (overprivileged)
- No RBAC policies
- Secrets stored in plain YAML
- No network policies
- No resource limits enforcement
- No image scanning

---

## The 4 C's of Cloud Native Security

The 4 C's framework provides a layered approach to security:

```
┌─────────────────────────────────────┐
│         Cloud (Infrastructure)      │
│  - IAM/RBAC                          │
│  - Network Security                  │
│  - Encryption                        │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│      Cluster (Kubernetes)            │
│  - API Server Security               │
│  - Network Policies                  │
│  - Pod Security Standards            │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│      Container (Runtime)             │
│  - Image Scanning                    │
│  - Least Privilege                  │
│  - Resource Limits                  │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│         Code (Application)           │
│  - Secure Coding                     │
│  - Dependency Scanning               │
│  - SAST/DAST                        │
└─────────────────────────────────────┘
```

### 1. Code Security

**Best Practices:**
- Use secure coding practices
- Scan dependencies for vulnerabilities
- Implement input validation
- Use least privilege in application logic
- Encrypt sensitive data

**Example: Scanning Go Dependencies**

```bash
# Install gosec for Go security scanning
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Scan the loan application
cd general_info/projects/Project_Loan/government_api/backend

# First, ensure dependencies are downloaded
go mod download
go mod tidy

# Scan with gosec 
gosec ./...

# Scan loan validator portal
cd ../../loan_validator_portal/backend

# Ensure dependencies are downloaded
go mod download
go mod tidy

# Scan excluding instrumented file to avoid duplicate declaration errors
gosec -exclude-dir=instrumented ./...

# Or scan only main.go if instrumented file causes issues
gosec main.go
```

### 2. Container Security

**Best Practices:**
- Use minimal base images (Alpine, Distroless)
- Implement multi-stage builds
- Run as non-root user
- Scan images before deployment
- Use specific image tags (avoid `latest`)

**Example: Secure Dockerfile**

```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app

# Final minimal image
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/app .
USER 65534:65534
ENTRYPOINT ["./app"]
```

### 3. Cluster Security

**Best Practices:**
- Enable RBAC
- Use Pod Security Standards
- Implement Network Policies
- Enable audit logging
- Use admission controllers
- Encrypt secrets at rest
- Use external secret management (Vault, AWS Secrets Manager)

**Example: Using External Secrets**

```yaml
apiVersion: external-secrets.io/v1
kind: ExternalSecret
metadata:
  name: postgres-secret
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: postgres-secret
    creationPolicy: Owner
  data:
  - secretKey: POSTGRES_PASSWORD
    remoteRef:
      key: database/postgres
      property: password
```

We'll cover these in detail in subsequent sections.

### 4. Cloud Security

**Best Practices:**
- Enable encryption in transit (TLS)
- Implement network segmentation
- Use IAM for cloud resources

---

## RBAC (Role-Based Access Control)

RBAC provides fine-grained access control to Kubernetes resources.

### Core Concepts

- **ServiceAccount**: Identity for pods
- **Role**: Permissions within a namespace
- **RoleBinding**: Grants Role to ServiceAccount/User
- **ClusterRole**: Permissions cluster-wide
- **ClusterRoleBinding**: Grants ClusterRole cluster-wide

### Lab 1: Implement RBAC for Loan Application

**Objective:** Create least-privilege service accounts for each component.

**Steps:**

1. Create ServiceAccounts for each component
2. Create Roles with minimal required permissions
3. Bind Roles to ServiceAccounts
4. Update deployments to use ServiceAccounts

See [examples/rbac/](examples/rbac/) for complete examples.

**Example: ServiceAccount for Government API**

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: government-api-sa
  namespace: default
```

**Example: Role for Government API**

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: government-api-role
  namespace: default
rules:
- apiGroups: [""]
  resources: ["configmaps", "secrets"]
  resourceNames: ["government-api-config", "government-api-secret"]
  verbs: ["get"]
```

**Example: RoleBinding**

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: government-api-binding
  namespace: default
subjects:
- kind: ServiceAccount
  name: government-api-sa
  namespace: default
roleRef:
  kind: Role
  name: government-api-role
  apiGroup: rbac.authorization.k8s.io
```

**Update Deployment:**

```yaml
spec:
  template:
    spec:
      serviceAccountName: government-api-sa
      # ... rest of spec
```

**Verification:**

```bash
# Apply RBAC resources
kubectl apply -f examples/rbac/

# Verify service accounts
kubectl get serviceaccounts

# Test access (should fail without proper permissions)
kubectl auth can-i get pods --as=system:serviceaccount:default:government-api-sa
```

### Best Practices

1. **Principle of Least Privilege**: Grant minimum required permissions
2. **Separate ServiceAccounts**: One per application component
3. **Namespace Isolation**: Use Roles instead of ClusterRoles when possible
4. **Regular Audits**: Review RBAC policies regularly
5. **Avoid Default SA**: Never use `default` service account

---

## Pod Security Standards

Pod Security Standards define three policies: **Privileged**, **Baseline**, and **Restricted**.

### Policy Levels

| Profile | Description | Use Case |
|--------|-------------|----------|
| **Privileged** | Unrestricted | System components only |
| **Baseline** | Prevents known privilege escalations | Most applications |
| **Restricted** | Hardened security | Production workloads |

### Lab 2: Apply Pod Security Standards

**Objective:** Enforce Pod Security Standards using a dedicated namespace for pod security demonstrations.

**Steps:**

1. Enable Pod Security Admission
2. Create pod-security namespace with policy labels
3. Update pods to comply with policy
4. Test policy enforcement

See [examples/pod-security/](examples/pod-security/) for examples.

**Create Pod Security Namespace:**

```bash
# Create the pod-security namespace with Restricted policy
kubectl apply -f examples/pod-security/namespace.yaml

# Verify namespace was created
kubectl get namespace pod-security

# Verify labels
kubectl get namespace pod-security -o yaml | grep pod-security.kubernetes.io
```

**Alternative: Label Existing Namespace:**

```bash
# Apply Baseline policy (warnings only)
kubectl label namespace pod-security pod-security.kubernetes.io/enforce=baseline \
  pod-security.kubernetes.io/audit=baseline \
  pod-security.kubernetes.io/warn=baseline

# Or apply Restricted policy (blocks violations)
kubectl label namespace pod-security pod-security.kubernetes.io/enforce=restricted \
  pod-security.kubernetes.io/audit=restricted \
  pod-security.kubernetes.io/warn=restricted
```

**Secure Pod Configuration:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: government-api
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
        seccompProfile:
          type: RuntimeDefault
      containers:
      - name: government-api
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
          runAsNonRoot: true
          runAsUser: 1000
```

**Verification:**

```bash
# Check pod security context
kubectl get pod <pod-name> -n pod-security -o jsonpath='{.spec.securityContext}'

# Test policy enforcement
kubectl apply -f examples/pod-security/violating-pod.yaml
# Should be rejected if enforce=restricted

# Check pods in pod-security namespace
kubectl get pods -n pod-security
```

---

## Container Image Scanning

Scan container images for vulnerabilities before deployment.

### Tools

- **Trivy**: Open-source, comprehensive scanning
- **Aqua Security**: Enterprise-grade scanning
- **Snyk**: Developer-friendly scanning

### Lab 3: Scan Loan Application Images

**Objective:** Scan container images for vulnerabilities and integrate scanning into the deployment process.

**Steps:**

1. Install Trivy
2. Scan application images
3. Integrate scanning into CI/CD pipeline
4. Create policies to enforce image security best practices (e.g., block latest tags)

See [examples/trivy-scan/](examples/trivy-scan/) for scripts and results.

**Scan Images:**

```bash
# Scan government API image
trivy image daasnigeria/daasrepo:government-api

# Scan with JSON output
trivy image -f json -o scan-results.json daasnigeria/daasrepo:government-api

# Scan with exit code on vulnerabilities
trivy image --exit-code 1 --severity HIGH,CRITICAL daasnigeria/daasrepo:government-api
```

**Integrate with CI/CD:**

---

## Policy Management

Policy engines enforce security policies at admission time.

### OPA/Gatekeeper vs Kyverno

| Feature | OPA/Gatekeeper | Kyverno |
|---------|----------------|---------|
| Language | Rego | YAML |
| Learning Curve | Steep | Gentle |
| Performance | Fast | Fast |
| Community | Large | Growing |

We'll use **Kyverno** for its simplicity and YAML-based policies.

### Lab 4: Enforce Policies with Kyverno

**Objective:** Install Kyverno and create policies for the loan application.

**Steps:**

1. Install Kyverno
2. Create validation policies
3. Create mutation policies
4. Test policy enforcement

See [examples/kyverno-policies/](examples/kyverno-policies/) for examples.

**Install Kyverno:**

```bash
Use any of the best installation methods for Kyverno applicable to you.
https://kyverno.io/docs/installation/methods/

```

**Example Policy: Require Resource Limits**

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-resource-limits
spec:
  validationFailureAction: enforce
  rules:
  - name: check-resource-limits
    match:
      resources:
        kinds:
        - Pod
    validate:
      message: "CPU and memory limits are required"
      pattern:
        spec:
          containers:
          - resources:
              limits:
                memory: "?*"
                cpu: "?*"
```

**Example Policy: Block Privileged Containers**

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: disallow-privileged
spec:
  validationFailureAction: enforce
  rules:
  - name: check-privileged
    match:
      resources:
        kinds:
        - Pod
    validate:
      message: "Privileged containers are not allowed"
      pattern:
        spec:
          containers:
          - securityContext:
              privileged: "false"
```

**Example Policy: Require Specific Labels**

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-labels
spec:
  validationFailureAction: enforce
  rules:
  - name: check-labels
    match:
      resources:
        kinds:
        - Pod
    validate:
      message: "Labels 'app' and 'version' are required"
      pattern:
        metadata:
          labels:
            app: "?*"
            version: "?*"
```

**Test Policy Enforcement:**

```bash
# Apply policies
kubectl apply -f examples/kyverno-policies/

# Try to create a pod without resource limits (should fail)
kubectl apply -f examples/kyverno-policies/test-violating-pod.yaml

# Create compliant pod (should succeed)
kubectl apply -f examples/kyverno-policies/test-compliant-pod.yaml
```

---
