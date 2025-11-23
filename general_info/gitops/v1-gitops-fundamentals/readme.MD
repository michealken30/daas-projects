# What is GitOps?

GitOps is an operational model that uses Git as the single source of truth for declaring and managing infrastructure and application state. Changes are made by committing to Git; automated agents reconcile the real environment to the desired state.

![alt text](Gitops-cover.webp)

## Four Commandments/Principles of GitOps
-  Git for Everything
-  V for Versioning
-  Auto-Pull
-  Loop the Loop 


## Advantages
- Faster, more reliable deployments through automation.  
- Clear audit trail and easy rollback via Git history.  
- Improved consistency across environments.  
- Better collaboration using standard Git workflows (PRs, reviews).  
- Easier recovery and troubleshooting because desired state is explicit.


## Examples of GitOps tools

- Argo CD — Declarative continuous delivery for Kubernetes that continuously syncs cluster state from Git.  
- Flux (v2) — GitOps operator for Kubernetes (reconciles manifests, supports Helm/Kustomize).  
- OpenShift GitOps — Red Hat’s Argo CD distribution for OpenShift clusters.  
- Jenkins X — CI/CD platform that uses GitOps workflows (Tekton-based).    
- Atlantis — Automates Terraform via Git pull-request-driven workflows for infra-as-code.

> ## Question?
> 1. How does GitOps workflow differ from traditional CI/CD workflow?
> 2. What is Push vs Pull GitOps model?



# What is Argo CD

Argo CD is a declarative, popular open-source GitOps continuous delivery tool for managing Kubernetes applications.

![alt text](image1-31.png)

## Components of Argo CD
![alt text](argocd_architecture.webp)

1. ArgoCD Server
2. ArgoCD Repo Server
3. Argo CD Application Controller
4. Argo CD ApplicationSet Controller
5. Argo CD Dex Server
6. Argo CD Notification Controller
7. Argo CD Redis


# Adding Clusters to ArgoCD
Clusters are stored as secrets in ArgoCD and each secret must have the label `argocd.argoproj.io/secret-type: cluster`


**Declaratively**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: <secret-name>
  namespace: argocd
  labels:
    argocd.argoproj.io/secret-type: cluster
type: Opaque
stringData:
  name: <cluster-name>
  server: <cluster-url>
  config: |
    {
      "tlsClientConfig": {
        "insecure": false,
        "caData": "<BASE64_ENCODED_CA_CERT>",
        "certData": "<BASE64_ENCODED_CLIENT_CERT>",
        "keyData": "<BASE64_ENCODED_CLIENT_KEY>"
      }
    }
```

**Using ArgoCD CLI**

```bash
kubectl config use-context <cluster-context-name>
argocd cluster add <cluster-context-name> --name <cluster-name>
```


# Adding Repositories to ArgoCD
Similar to Cluster, it is of secret resource type and must have the label, `argocd.argoproj.io/secret-type: repository`

**Declaratively**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: <secret-name>
  namespace: argocd
  labels:
    argocd.argoproj.io/secret-type: repository
stringData:
  type: git
  url: <repo-url>
```

**Using ArgoCD CLI**
```bash
argocd repo add <repo-url> --name <repo-name>
```

## Adding Private Repos
Requires credentials. Credentials can be off the following types
- SSH (key based)
- HTTPS (simple auth)
- GitHub Apps
- Helm 

> ## Question?
> Is there any security risk with declaring SourceRepos and clusters in ArgoCD as k8s secrets?


# Creating your first ArgoCD App

# **What is an ArgoCD App?**
A group of Kubernetes resources as defined by a manifest. This is a Custom Resource Definition (CRD).

## ArgoCD Application Source Types
ArgoCD supports 8 main source types:
- Git Directory - Plain YAML/JSON files
- Kustomize - Overlay-based configuration
- Helm Chart - From Helm repositories
- Helm from Git - Helm charts in Git repos
- Helm OCI - Helm charts in OCI registries
- Jsonnet - Jsonnet templating
- Plugin - Custom config management tools
- Multiple Sources - Combine any of the above (available in ArgoCD 2.6+)

>                       **Source Type != Repository Type**


# Sync Policies and Sync Options

# Resourc Hooks (Sync Phase)
- PreSync: Before the sync operation
- Sync: During normal sync (default)
- PostSync: After all Sync resources are healthy
- SyncFail: If sync fails
- Skip: Never sync this resource
    
Use the annotation: `argocd.argoproj.io/hook: PreSync`

# 🌊 Sync Waves (
This is the order of deployment of resource create by ArgoCD app. The lower the number the higher the priority
> -3 > -2 > -1 > 0 > 1 > 2 > 3


# Assignment 1:
Deploy the k8s application using ArgoCD with the following requirements:

1. Configure sync waves to control resource creation order
2. Enable automatic namespace creation
3. Keep auto-sync disabled (manual sync only)
4. Target your personal branch (not main)
5. Make a change to any resource manifest and demonstrate the diff before syncing

Deliverables:
- ArgoCD Application manifest
- Screenshot/output showing the diff detection
- Brief explanation of your sync wave strategy

# App Project and ApplicationSet

App Project is a way to group your applications and set rules/permission for it especially if you are using multiple teams.

ApplicationSets help automate the creation of many ArgoCD apps across multiple clusters.
It uses a template and a generator (like a list of clusters or Git directories) to automatically mass-produce identical or slightly varied Applications, solving the problem of deploying the same thing to many places. 

## Types Of ApplicationSets
- ApplicationSet List Generator
- ApplicationSet Cluster Generator
- ApplicationSet Git Directory Generator
- ApplicationSet Matrix Generator
- ApplicationSet Merge Generators
- ApplicationSet Pull Request Generators
