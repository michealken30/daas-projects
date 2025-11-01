# 📚 Week 2 Assignment: Security Hardening & Best Practices

## Overview
This assignment focuses on implementing Kubernetes security best practices, proper image versioning, pod security standards, and service account security for our microservices application.

---

## 🎯 Assignment 1: "Secure the JWT Secret"

### Current State
JWT_SECRET is stored in a ConfigMap (insecure)

### Task
Convert to use Kubernetes Secret for security best practices

### Requirements
- [ ] Create a Kubernetes Secret for JWT_SECRET
- [ ] Update deployment manifests to reference the Secret instead of ConfigMap
- [ ] Test that the application still works with the new Secret
- [ ] Verify the secret is not visible in plain text when queried

### Validation Steps
1. `kubectl get secrets` - Verify secret exists
2. `kubectl describe secret app-secrets` - Confirm secret data is encoded
3. Test login/register functionality still works
4. Verify JWT tokens are properly generated and validated

---

## 🎯 Assignment 2: "Implement Semantic Versioning"

### Current State
Images use generic tags like `daasnigeria/daasrepo:auth-service`

### Task
Update all microservice images to use Semantic Versioning (SemVer)

### Requirements
- [ ] Update auth-service image to use proper semantic versioning
- [ ] Update upload-service image to use proper semantic versioning
- [ ] Update api-gateway image to use proper semantic versioning
- [ ] Change `imagePullPolicy` from `Always` to `IfNotPresent`
- [ ] Document the versioning strategy

---

## 🎯 Assignment 3: "Implement Pod Security Standards"

### Current State
Deployments do not follow Kubernetes security best practices

### Task
Harden the pods by implementing Pod Security Standards

### Requirements
- [ ] Implement **runAsNonRoot** security context
- [ ] Set **readOnlyRootFilesystem** where possible
- [ ] Configure **allowPrivilegeEscalation: false**
- [ ] Add **securityContext** with appropriate user ID
- [ ] Implement **resource limits** and **requests** (already done)
- [ ] Add **Pod Security Policy** or **Pod Security Standards**
- [ ] Drop all unnecessary capabilities
- [ ] Set appropriate filesystem group ownership

---

## 🎯 Assignment 4: "Service Account Security & RBAC"

### Current State
Deployments use the default service account with excessive permissions

### Task
Implement proper service account security and Role-Based Access Control (RBAC)

### Requirements
- [ ] Create dedicated service accounts for each microservice
- [ ] Implement least-privilege RBAC policies
- [ ] Create ClusterRoles/Roles with minimal required permissions
- [ ] Create RoleBindings/ClusterRoleBindings for service accounts
- [ ] Disable automountServiceAccountToken where not needed
- [ ] Update deployments to use dedicated service accounts
- [ ] Verify services function with restricted permissions

---

## 📋 Deliverables

### 1. Updated Kubernetes Manifests
- [ ] `00-common.yaml` - With Kubernetes Secret
- [ ] `01-auth-service.yaml` - With security context, versioned image, and service account
- [ ] `02-upload-service.yaml` - With security context, versioned image, and service account
- [ ] `03-api-gateway.yaml` - With security context, versioned image, and service account
- [ ] `05-rbac.yaml` - New file with service accounts, roles, and role bindings

### 2. Documentation
- [ ] **Security-Changes.md** - Document all security improvements made
- [ ] **Versioning-Strategy.md** - Explain your versioning approach
- [ ] **RBAC-Design.md** - Document service account and RBAC implementation
- [ ] **Testing-Results.md** - Proof that application works after changes

### 3. Testing Evidence
- [ ] Screenshots of `kubectl get secrets`
- [ ] Screenshots of `kubectl get serviceaccounts`
- [ ] Screenshots of `kubectl get roles,rolebindings`
- [ ] Screenshots of successful application functionality
- [ ] Pod security compliance verification

---

## 🧪 Testing Checklist

### Functionality Testing
- [ ] User registration works
- [ ] User login works  
- [ ] File upload works
- [ ] JWT tokens are properly validated
- [ ] All services communicate correctly

### Security Testing
- [ ] Secrets are base64 encoded
- [ ] Pods run as non-root user
- [ ] Root filesystem is read-only (where applicable)
- [ ] No privilege escalation possible
- [ ] Resource limits are enforced
- [ ] Service accounts have minimal permissions
- [ ] Default service account token is not mounted unnecessarily

### Deployment Testing
- [ ] All pods start successfully
- [ ] Health checks pass
- [ ] Services are accessible
- [ ] No security warnings in pod descriptions
- [ ] RBAC policies are properly enforced

---


## 📚 Resources

### Documentation
- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [Semantic Versioning](https://semver.org/)
- [Security Context](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#securitycontext-v1-core)
- [RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
- [Service Accounts](https://kubernetes.io/docs/concepts/security/service-accounts/)

### Best Practices
- Never store sensitive data in ConfigMaps
- Always use specific image versions in production
- Run containers as non-root users when possible
- Implement defense-in-depth security strategies
- Follow principle of least privilege for service accounts
- Regularly audit and rotate secrets

---

## ⏰ Submission Guidelines

### Deadline
Submit by Thursday 6th November at 11:59PM

### Submission Format
1. **Git Repository**: Push all changes to your branch
2. **Documentation**: Include all required markdown files
3. **Demo**: Be prepared to demonstrate the working application

### Evaluation Criteria
- **Security Implementation** (35%)
- **RBAC Implementation** (25%)
- **Functionality Preservation** (25%)
- **Documentation Quality** (15%)

---

**Good luck! 🚀**

*Remember: Security is not a feature, it's a foundation.*