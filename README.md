# Tutorial: How To Write a Reconciler Using K8s Controller-Runtime!

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime)

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

## Step-By-Step Guide

Once the [dev env setup](#local-dev) is ready, each step will use a [git tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging). 

There are 7 steps, which match to 7 tags:

<img width="458" alt="Screen Shot 2022-10-26 at 4 07 05 PM" src="https://user-images.githubusercontent.com/407675/198126519-28475fc7-6ef8-44c6-b380-4287631032cd.png">

To move from one tag to another, checkout the tag and create a new branch from it.

For example, move to the first step:
```bash
git checkout tags/s1 -b s1
```

Then, run the tests: 
```bash
make test
```

And so on for each step!