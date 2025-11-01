# 📚 Week 2 Assignment: Security Hardening & Best Practices

## Overview
This assignment focuses on implementing Kubernetes security best practices, proper image versioning, and pod security standards for our microservices application.

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

```

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
- [ ] Update auth-service image to: `daasnigeria/daasrepo:auth-service-vx.x.x`
- [ ] Update upload-service image to: `daasnigeria/daasrepo:upload-service-vx.x.x`
- [ ] Update api-gateway image to: `daasnigeria/daasrepo:api-gateway-vx.x.x`
- [ ] Change `imagePullPolicy` from `Always` to `IfNotPresent`
- [ ] Document the versioning strategy

### Semantic Versioning Format
```
MAJOR.MINOR.PATCH
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)
```


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


### Pod Security Standards Levels
- **Privileged**: No restrictions (current state)
- **Baseline**: Minimal restrictions
- **Restricted**: Heavily restricted (target)

---

## 📋 Deliverables

### 1. Updated Kubernetes Manifests
- [ ] `00-common.yaml` - With Kubernetes Secret
- [ ] `01-auth-service.yaml` - With security context and versioned image
- [ ] `02-upload-service.yaml` - With security context and versioned image  
- [ ] `03-api-gateway.yaml` - With security context and versioned image

### 2. Documentation
- [ ] **Security-Changes.md** - Document all security improvements made
- [ ] **Versioning-Strategy.md** - Explain your versioning approach
- [ ] **Testing-Results.md** - Proof that application works after changes

### 3. Testing Evidence
- [ ] Screenshots of `kubectl get secrets`
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

### Deployment Testing
- [ ] All pods start successfully
- [ ] Health checks pass
- [ ] Services are accessible
- [ ] No security warnings in pod descriptions

---

## 📚 Resources

### Documentation
- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [Semantic Versioning](https://semver.org/)
- [Security Context](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#securitycontext-v1-core)

### Best Practices
- Never store sensitive data in ConfigMaps
- Always use specific image versions in production
- Run containers as non-root users when possible
- Implement defense-in-depth security strategies

---

## ⏰ Submission Guidelines

### Deadline
Submit by Thursday 6th November at 11:59PM

### Submission Format
1. **Git Repository**: Push all changes to your branch
2. **Documentation**: Include all required markdown files
3. **Demo**: Be prepared to demonstrate the working application

### Evaluation Criteria
- **Security Implementation** (40%)
- **Functionality Preservation** (30%)
- **Documentation Quality** (20%)
- **Best Practices Adherence** (10%)

---

**Good luck! 🚀**

*Remember: Security is not a feature, it's a foundation.*