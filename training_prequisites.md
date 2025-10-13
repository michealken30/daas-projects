# Training Prerequisites

Before starting this daas Kubernetes bootcamp, it is important to have foundational knowledge in several key areas. Below are the essential prerequisites and what you should know in each category:

---

## 1. Linux
- **Why:** Most cloud-native and DevOps tools run on Linux.
- **What to Learn:**
  - Linux filesystem structure (`/home`, `/etc`, `/var`, etc.)
  - Basic commands: `ls`, `cd`, `pwd`, `cp`, `mv`, `rm`, `cat`, `less`, `grep`, `find`, `chmod`, `chown`, `ps`, `top`, `kill`, `nano` or `vim`
  - File permissions and ownership
  - How to install packages (e.g., `apt`, `yum`)
  - Navigating and editing files in the terminal

## 2. Git
- **Why:** Version control is essential for collaboration and code management.
- **What to Learn:**
  - Cloning a repository: `git clone <repo_url>`
  - Creating and switching branches: `git branch`, `git checkout -b <branch>`
  - Adding and committing changes: `git add .`, `git commit -m "message"`
  - Pushing and pulling: `git push`, `git pull`
  - Resolving merge conflicts
  - Viewing history: `git log`, `git status`

## 3. Basic kubectl Commands
- **Why:** `kubectl` is the main tool for interacting with Kubernetes clusters.
- **What to Learn:**
  - Listing resources: `kubectl get pods`, `kubectl get services`, `kubectl get nodes`
  - Describing resources: `kubectl describe pod <pod_name>`
  - Creating resources from manifest files: `kubectl apply -f <file.yaml>`
  - Deleting resources: `kubectl delete pod <pod_name>`
  - Configuring and exposing services
  - Viewing logs: `kubectl logs <pod_name>`

## 4. Docker
- **Why:** Containers are the foundation of Kubernetes.
- **What to Learn:**
  - Building images: `docker build -t <image_name> .`
  - Running containers: `docker run -d -p 8080:80 <image_name>`
  - Listing containers and images: `docker ps`, `docker images`
  - Stopping and removing containers: `docker stop <container_id>`, `docker rm <container_id>`
  - Removing images: `docker rmi <image_id>`
  - Viewing logs: `docker logs <container_id>`

## 5. Networking
- **Why:** Understanding networking is crucial for troubleshooting and designing distributed systems.
- **What to Learn:**
  - DNS basics (domain names, name resolution)
  - TCP and UDP protocols (differences, use cases)
  - Common ports (e.g., 80 for HTTP, 443 for HTTPS, 22 for SSH)
  - IP addressing and subnets
  - Firewalls and security groups
  - Tools: `ping`, `curl`, `netstat`, `traceroute`, `nslookup`

## 6. Additional Recommended Knowledge
- **YAML:** Learn how to read and write YAML files, as they are used for Kubernetes manifests.
- **Cloud Basics:** Familiarity with cloud concepts (IaaS, PaaS, SaaS, regions, zones).
- **CI/CD Concepts:** Understanding the basics of continuous integration and deployment.
- **Scripting:** Basic shell scripting (bash) for automation.

---

**Tip:** Strengthening these foundational skills will make your Kubernetes and Cloud learning journey much smoother and more effective.
