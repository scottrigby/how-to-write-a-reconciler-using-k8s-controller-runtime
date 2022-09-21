Commands run so far

1. kubebuilder init --domain kubecon.na --repo github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/projects/cfp
3. kubebuilder create api --group talks --version v1 --kind Speaker
4. kubebuilder create api --group talks --version v1 --kind Proposals
5. Update `api/*_types.go` file and run `make manifests` 
