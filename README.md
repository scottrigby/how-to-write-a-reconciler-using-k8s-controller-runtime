# Tutorial: How To Write a Reconciler Using K8s Controller-Runtime!

Git repo for talk at KubeCon NA 2022. Scott Rigby, Somtochi Onyekwere, Niki Manoledaki & Soulé Ba, Weaveworks; Amine Hilaly, Amazon Web Services. <https://sched.co/182Hg>

> This tutorial walks you through building your own controller using controller runtime, the set of common libraries on which core controllers are built. We'll use Kubebuilder, a framework for building APIs using custom resource definitions (CRDs). We'll also explain lesser-documented best practices and conventions for writing controllers that the community has developed through trial and error learning, through projects such as Flux and Cluster API.

## How this repo is organized

Since we'll be building a controller with reconcilers, we'll need something to reconcile.

For this, we've built a simple API for submitting KubeCon CFP proposals, using CRDs.

- CFP API code is in [`/cfp-api`](cfp-api/README.md)
- controller code is in [`/cfp`](cfp/README.md)

## Local Dev

Fork this repository, then clone it:
```bash
git clone git@github.com:<your username>/how-to-write-a-reconciler-using-k8s-controller-runtime.git
```

Dependencies:
- [Go (1.19+)](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)

For the local dev setup, there are a few options:
1. **GitPod**. GitPod is the preferred method, to ensure all users are running in the same environment regardless of their local machine OS. Please click on the button "Open in Gitpod" at the top of this README. Click through the default settings until you arrive at a page that looks like VSCode. Then, wait while the dependencies load (approximately ~6 minutes).
2. **Vagrant**. As a backup, there is a Vagrantfile with instructions [here](dev/vagrant/README.md).
3. **DIY**. If you have all the required dependencies, you can go ahead and spinup a dev environment:
```bash
cd cfp
make docker-build
make setup-kind
make dev-deploy
export KUBECONFIG=/tmp/cfp-api-test-kubeconfig
```