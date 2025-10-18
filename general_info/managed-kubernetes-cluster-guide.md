# Running Managed Kubernetes Clusters on EKS, AKS, GKE and Oracle Cloud

To connect and manage a managed Kubernetes cluster (EKS, AKS, GKE or Oracle Cloud), you will need a Linux environment. This can be:
- Native Linux
- Linux on macOS (via Terminal)
- Linux on Windows (using WSL - Windows Subsystem for Linux)

This is necessary because you will need to export the kubeconfig file from your cloud provider and configure it on your local Linux environment to run `kubectl` commands.

---

## 1. Amazon EKS (Elastic Kubernetes Service)
- **Steps:**
  1. Create an EKS cluster using the AWS Console or CLI.
  2. Install and configure the AWS CLI and `kubectl` on your Linux environment.
  3. Update your kubeconfig:
     ```bash
     aws eks --region <region> update-kubeconfig --name <cluster_name>
     ```
  4. Test your connection:
     ```bash
     kubectl get nodes
     ```

## 2. Azure AKS (Azure Kubernetes Service)
- **Steps:**
  1. Create an AKS cluster using the Azure Portal or CLI.
  2. Install and configure the Azure CLI and `kubectl` on your Linux environment.
  3. Get your kubeconfig:
     ```bash
     az aks get-credentials --resource-group <resource_group> --name <cluster_name>
     ```
  4. Test your connection:
     ```bash
     kubectl get nodes
     ```

## 3. Google GKE (Google Kubernetes Engine)
- **Steps:**
  1. Create a GKE cluster using the Google Cloud Console or CLI.
  2. Install and configure the Google Cloud SDK (`gcloud`) and `kubectl` on your Linux environment.
  3. Get your kubeconfig:
     ```bash
     gcloud container clusters get-credentials <cluster_name> --zone <zone> --project <project_id>
     ```
  4. Test your connection:
     ```bash
     kubectl get nodes
     ```

## 4. Oracle Cloud OKE (Oracle Kubernetes Engine)
- **Steps:**
  1. Create an OKE cluster using the Oracle Cloud Console.
  2. Install and configure the Oracle Cloud CLI and `kubectl` on your Linux environment.
  3. Download the kubeconfig file from the Oracle Cloud Console.
  4. Set the KUBECONFIG environment variable:
     ```bash
     export KUBECONFIG=/path/to/your/kubeconfig
     ```
  5. Test your connection:
     ```bash
     kubectl get nodes
     ```

---

**Note:**
- Always ensure your cloud CLI tools and `kubectl` are installed and up to date.
- You may need to authenticate with your cloud provider before running these commands.
- The kubeconfig file allows your local `kubectl` to securely connect to your managed Kubernetes cluster.
