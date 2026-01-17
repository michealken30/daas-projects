# Network Policy Examples

This directory contains NetworkPolicy examples for the loan application.

## Files

- `deny-all.yaml`: Default deny-all policy (apply first)
- `allow-postgres.yaml`: Restricts PostgreSQL access to authorized apps
- `allow-government-api.yaml`: Restricts Government API access
- `allow-loan-validator.yaml`: Allows external access to Loan Validator Portal

## Network Policy Strategy

1. **Default Deny**: Start with deny-all policy
2. **Whitelist**: Add specific allow policies
3. **Least Privilege**: Only allow necessary traffic

## Usage

```bash
# Apply in order
kubectl apply -f deny-all.yaml
kubectl apply -f allow-postgres.yaml
kubectl apply -f allow-government-api.yaml
kubectl apply -f allow-loan-validator.yaml

# Verify policies
kubectl get networkpolicies

# Test connectivity
kubectl run test-pod --image=busybox --rm -it -- sh
# Try to connect to services
```

## Traffic Flow

```
Internet
  ↓
Loan Validator Portal (8080) ← Allowed
  ↓
Government API (8081) ← Only from Loan Validator
  ↓
PostgreSQL (5432) ← Only from Government API and Loan Validator
```

## Best Practices

1. **Default Deny**: Always start with deny-all
2. **Namespace Isolation**: Use namespace selectors for cross-namespace traffic
3. **Port Specificity**: Only allow specific ports
4. **Label Consistency**: Use consistent labels for pod selection
5. **DNS Access**: Always allow DNS (UDP 53) for egress

