# how-to-write-a-reconciler-using-k8s-controller-runtime
How To Write a Reconciler Using k8s Controller-Runtime!


# Using the provided vagrant File

We provide a vagrantfile that will create a VM with all the necessary tools installed.
To use it, simply run `vagrant up` and then `vagrant ssh` to get into the VM.
You can then run `make` to build the project.

Please not that building the vm will take a while, as it will download all the necessary tools.

## Prerequisites

This will work with virtual box and vagrant installed on your machine.

for vmware fusion, you will need to install the vagrant-vmware-fusion plugin.
https://www.vagrantup.com/vmware/downloads

## Steps

Into a terminal run the following commands:

1. Add the provider
```sh
⋊> ~ vagrant box add hashicorp/bionic64
```
For mac users, make sure you have the latest version of virtual box installed.
You may have to check your macbook's setting --> security, as default, it blocks the application
running from oracle, just allow and approve it, this problem is gone.

2. Start the vagrant box
```sh
⋊> ~ cd ./dev/vagrant/    
⋊> ~/dev/vagrant vagrant up
```

3. SSH into the vagrant box
```sh
⋊> ~/dev/vagrant vagrant ssh
```

The repository is synced to the vagrant box at `/home/vagrant/app`. You can edit
the files in your favorite editor and run the commands in the vagrant box.
