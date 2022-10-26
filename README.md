# Tutorial: How To Write a Reconciler Using K8s Controller-Runtime!

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime)

Git repo for talk at KubeCon NA 2022. Scott Rigby, Somtochi Onyekwere, Niki Manoledaki & Soulé Ba, Weaveworks; Amine Hilaly, Amazon Web Services. <https://sched.co/182Hg>

> This tutorial walks you through building your own controller using controller runtime, the set of common libraries on which core controllers are built. We'll use Kubebuilder, a framework for building APIs using custom resource definitions (CRDs). We'll also explain lesser-documented best practices and conventions for writing controllers that the community has developed through trial and error learning, through projects such as Flux and Cluster API.

## How this repo is organized

Since we'll be building a controller with reconcilers, we'll need something to reconcile.

For this, we've built a simple API for submitting KubeCon CFP proposals, using CRDs.

- CFP API code is in [`/cfp-api`](cfp-api/README.md)
- controller code is in [`/cfp`](cfp/README.md)

## Local dev

For local dev, there are two options:

- GitPod is the preferred method, to ensure all users are running in the same environment regardless of their local machine OS. TO-DO
- As a backup, there is a Vagrantfile [with instructions](dev/vagrant/README.md)
