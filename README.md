# how-to-write-a-reconciler-using-k8s-controller-runtime
How To Write a Reconciler Using k8s Controller-Runtime!


## Using the provided Vagrantfile

We provide a vagrantfile that will create a VM with all the necessary tools installed.
To use it, simply run `vagrant up` and then `vagrant ssh` to get into the VM.
You can then run `make` to build the project.

Please not that building the vm will take a while, as it will download all the necessary tools.

### Prerequisites

This will work with virtual box and vagrant installed on your machine.

for vmware fusion, you will need to install the vagrant-vmware-fusion plugin.
https://www.vagrantup.com/vmware/downloads

### Steps

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

## API Documentation

Get the API running:
```
go run main.go
```

### Speakers

Create a Speaker:
```
curl -sd '{"ID":"default/ScottRigby","name":"Scott Rigby","bio":"Scott is a rad dad","email":"scott@email.com","timestamp":"0001-01-01T00:00:00Z"}' \
-H "Content-Type: application/json" \
-X POST localhost:8080/api/speakers | jq
{
  "ID": "default/ScottRigby",
  "Name": "Scott Rigby",
  "Bio": "Scott is a rad dad",
  "Email": "scott@email.com",
  "Timestamp": "0001-01-01T00:00:00Z"
}
```

Get all Speakers:
```
curl -sX GET localhost:8080/api/speakers | jq
[
  {
    "ID": "default/ScottRigby",
    "Name": "NewName",
    "Bio": "Scott is a rad dev",
    "Email": "scott@email.com",
    "Timestamp": "0001-01-01T00:00:00Z"
  }
]
```

Get a Speakerby ID:
```
curl -sX GET localhost:8080/api/speakers/default-ScottRigby | jq
[
  {
    "ID": "default/ScottRigby",
    "Name": "NewName",
    "Bio": "Scott is a rad dev",
    "Email": "scott@email.com",
    "Timestamp": "0001-01-01T00:00:00Z"
  }
]
```

Update a Speaker:
```
curl -sd '{"ID":"default/ScottRigby","name":"NewName","bio":"Scott is a rad dev","email":"scott@email.com","timestamp":"0001-01-01T00:00:00Z"}' \
-H "Content-Type: application/json" \
-X PUT localhost:8080/api/speakers/default-ScottRigby | jq
{
  "ID": "default/ScottRigby",
  "Name": "NewName",
  "Bio": "Scott is a rad dev",
  "Email": "scott@email.com",
  "Timestamp": "0001-01-01T00:00:00Z"
}
```

Delete a Speaker:
```
curl -X DELETE localhost:8080/api/speakers/default-ScottRigby
```