# Kubernetes Security Assignment

**Total Points: 10 marks**  
---

## Assignment Overview

You will build an **end-to-end secure CI/CD pipeline** that automatically scans, validates, and deploys a containerized application to a Kubernetes cluster while enforcing security best practices throughout the entire software delivery lifecycle.

We are using the loan application as our example application. But must be deployed in different namespaces.
- government-api namespace
- loan-validator-portal namespace
- database namespace
- observability namespace

Each application must have a different ArgoCD application.
---

## Assignment Objectives

Implement a fully automated, security-hardened deployment pipeline that demonstrates:
- Automated security scanning in CI/CD
- GitOps-based deployment practices
- Policy enforcement and compliance validation
- Least privilege access controls

---

## Technical Requirements

### 1. **CI/CD Pipeline Implementation - GitHub Actions** (3 marks)

Create a GitHub Actions workflow that includes:

**a) Code Security Scanning (1 marks)**
- Implement SAST (Static Application Security Testing) using any tool of your choice 
- Pipeline must fail if HIGH or CRITICAL vulnerabilities are detected
- Generate security scan reports as artifacts

**b) Container Image Security (1 marks)**
- Build Docker images with security best practices:
  - Use minimal base images (alpine, distroless, or scratch)
  - Run containers as non-root user
  - Implement multi-stage builds
- Scan images with **Trivy** for vulnerabilities
- Sign images using **Cosign** (bonus: implement image signing verification)
- Push images only if scan passes defined thresholds

**c) Kubernetes Manifest Validation (1 marks)**
- Scan Kubernetes manifests using **Kubesec**
- Validate against Pod Security Standards (restricted profile)
- Check for common misconfigurations (missing resource limits, privileged containers, etc.)

---

### 2. **GitOps Deployment with ArgoCD** (2 marks)
- Configure ArgoCD application for automatic sync with self-heal enabled
- Implement sync waves for proper deployment ordering
- Add health checks and sync hooks (PreSync, PostSync)
- Use ArgoCD Image Updater or implement automated image tag updates in manifests
- Ensure deployment happens automatically after pipeline pushes updated manifests

---

### 3. **Kubernetes Security Controls** (3 marks)

**a) RBAC Implementation (1 mark)**
- Create namespace-specific ServiceAccounts
- Implement least privilege RBAC:
- **Document your RBAC strategy** and justify permission choices

**b) Pod Security Standards & Policies (1 mark)**
- Enforce Pod Security Admission (PSA) at namespace level:
  - Apply `restricted` profile to application namespaces
- Implement policy enforcement using **OPA Gatekeeper** or **Kyverno**:
  - Require all containers to have resource limits
  - Block privileged containers
  - Enforce image pull policies
  - Require non-root containers
  - Block host namespace sharing
  - Mandate security contexts

**c) Network Policies (1 mark)**
- Implement default-deny NetworkPolicies
- Create specific ingress/egress rules for:
  - Application pods (only necessary communication)

---

### 4. **Documentation & Security Posture** (2 mark)

Provide comprehensive documentation including:

**Architecture Diagram**
- Complete pipeline flow from code commit to deployment
- Security checkpoints and validation gates
- Component interactions

---

---

## Bonus Points (Optional - up to +2 marks)

- **Image signing & verification** with Cosign/Notary (+0.5)
- **Secrets management** with Sealed Secrets or External Secrets Operator (+1)
- **Supply chain security** with SLSA provenance or SBOM generation (+0.5)

---

## Evaluation Criteria

Your solution will be evaluated on:

1. **Functionality** (40%) - Does the pipeline work end-to-end?
2. **Security Posture** (40%) - Are security controls properly implemented?
3. **Automation** (10%) - Is the process fully automated?
4. **Documentation** (10%) - Is the solution well-documented and reproducible?

---

## Helpful Resources

- GitHub Actions Security Hardening: https://docs.github.com/en/actions/security-guides
- ArgoCD Best Practices: https://argo-cd.readthedocs.io/
- Pod Security Standards: https://kubernetes.io/docs/concepts/security/pod-security-standards/
- Trivy Documentation: https://aquasecurity.github.io/trivy/
- OPA Gatekeeper: https://open-policy-agent.github.io/gatekeeper/website/docs/
- Kyverno: https://kyverno.io/docs/

---
