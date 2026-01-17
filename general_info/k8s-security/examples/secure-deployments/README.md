# Secure Deployment Examples

This directory contains fully hardened deployment configurations for the loan application.

## Security Features Applied

### 1. RBAC
- ✅ Dedicated ServiceAccounts for each component
- ✅ Least-privilege Roles
- ✅ RoleBindings connecting SAs to Roles

### 2. Pod Security Standards
- ✅ Non-root containers
- ✅ Read-only root filesystem
- ✅ Dropped capabilities
- ✅ Seccomp profiles
- ✅ Resource limits

### 3. Network Policies
- ✅ Default deny-all
- ✅ Whitelist specific traffic
- ✅ Namespace isolation

### 4. Image Security
- ✅ Scanned images
- ✅ Specific image tags (no latest)
- ✅ Minimal base images

### 5. Secrets Management
- ✅ Encrypted secrets
- ✅ External secret management ready

## Deployment Order

1. **RBAC Resources**
   ```bash
   kubectl apply -f ../rbac/
   ```

2. **Network Policies**
   ```bash
   kubectl apply -f ../network-policies/
   ```

3. **Kyverno Policies**
   ```bash
   kubectl apply -f ../kyverno-policies/
   ```

4. **Secure Deployments**
   ```bash
   kubectl apply -f secure-government-api.yaml
   kubectl apply -f secure-loan-validator.yaml
   ```

## Verification

```bash
# Check security contexts
kubectl get pods -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.securityContext}{"\n"}{end}'

# Check service accounts
kubectl get pods -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.serviceAccountName}{"\n"}{end}'

# Check network policies
kubectl get networkpolicies

# Check Kyverno policies
kubectl get clusterpolicies
```

## Before vs After

| Component | Before | After |
|-----------|---------|-------|
| ServiceAccount | default | Dedicated SA |
| Security Context | None | Non-root, read-only FS |
| Network Policy | None | Restrictive policies |
| Resource Limits | Present | Enforced by policy |
| Image Scanning | None | Pre-deployment |
| Audit Logging | Disabled | Enabled |

