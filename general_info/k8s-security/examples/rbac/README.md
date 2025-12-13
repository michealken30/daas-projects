# RBAC Examples

This directory contains RBAC (Role-Based Access Control) examples for the loan application.

## Files

- `serviceaccount-*.yaml`: ServiceAccounts for each component
- `role-*.yaml`: Roles defining permissions for each component
- `rolebinding-*.yaml`: RoleBindings connecting ServiceAccounts to Roles

## Usage

```bash
# Apply all RBAC resources
kubectl apply -f examples/rbac/

# Verify ServiceAccounts
kubectl get serviceaccounts

# Verify Roles and RoleBindings
kubectl get roles,rolebindings

# Test permissions
kubectl auth can-i get pods --as=system:serviceaccount:default:government-api-sa
```

## Security Principles Applied

1. **Least Privilege**: Each component only has permissions it needs
2. **Namespace Isolation**: All permissions are namespace-scoped (Roles, not ClusterRoles)
3. **Resource-Specific**: Permissions limited to specific resource names where possible
4. **Separate ServiceAccounts**: Each component has its own ServiceAccount

## Updating Deployments

After applying these RBAC resources, update your deployments to use the ServiceAccounts:

```yaml
spec:
  template:
    spec:
      serviceAccountName: government-api-sa  # or loan-validator-portal-sa, postgres-sa
```

