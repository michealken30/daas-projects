# Pod Security Standards Examples

This directory contains examples demonstrating Pod Security Standards.

## Files

- `namespace.yaml`: Creates the `pod-security` namespace with Pod Security Standards labels
- `secure-government-api.yaml`: Hardened Government API deployment
- `secure-loan-validator.yaml`: Hardened Loan Validator Portal deployment
- `violating-pod.yaml`: Example pod that violates security standards (for testing)
- `namespace-policy.yaml`: Instructions for applying Pod Security Standards

## Security Features

### Pod-Level Security Context
- `runAsNonRoot: true`: Prevents running as root
- `runAsUser: 1000`: Runs as non-root user
- `fsGroup: 1000`: Sets file system group
- `seccompProfile: RuntimeDefault`: Enables seccomp filtering

### Container-Level Security Context
- `allowPrivilegeEscalation: false`: Prevents privilege escalation
- `readOnlyRootFilesystem: true`: Makes root filesystem read-only
- `capabilities.drop: ["ALL"]`: Drops all Linux capabilities
- `runAsNonRoot: true`: Container runs as non-root
- `runAsUser: 1000`: Specific non-root user ID

## Usage

### Create Namespace with Pod Security Standards

```bash
# Create the pod-security namespace with Restricted policy
kubectl apply -f namespace.yaml

# Verify namespace was created
kubectl get namespace pod-security

# Verify labels
kubectl get namespace pod-security -o yaml | grep pod-security.kubernetes.io
```

### Apply Pod Security Standards to Existing Namespace

```bash
# Apply Baseline policy (warnings)
kubectl label namespace pod-security pod-security.kubernetes.io/enforce=baseline \
  pod-security.kubernetes.io/audit=baseline \
  pod-security.kubernetes.io/warn=baseline

# Apply Restricted policy (blocks violations)
kubectl label namespace pod-security pod-security.kubernetes.io/enforce=restricted \
  pod-security.kubernetes.io/audit=restricted \
  pod-security.kubernetes.io/warn=restricted
```

### Test Policy Enforcement

```bash
# Try to create violating pod (should fail with restricted policy)
kubectl apply -f violating-pod.yaml

# Apply secure deployment (should succeed)
kubectl apply -f secure-government-api.yaml

# Check pods in pod-security namespace
kubectl get pods -n pod-security
```

## Common Violations

| Violation | Fix |
|-----------|-----|
| Running as root | Add `runAsNonRoot: true`, `runAsUser: 1000` |
| Privileged container | Set `privileged: false` |
| Host network | Remove `hostNetwork: true` |
| Host PID | Remove `hostPID: true` |
| Default capabilities | Drop all: `capabilities.drop: ["ALL"]` |
| No resource limits | Add `resources.limits` |
| Read-write root FS | Set `readOnlyRootFilesystem: true` |

