___
# Lab 0: Getting Started
Welcome to CSE 190! This lab serves as a brief introduction to the following tools we will be using in CSE 190 and as a guide for setting up the lab infrastructure.
We provide a brief overview of the technologies and end with a set of
non-graded assignments. Even though these assignments are not graded,
you still need to finish them in order to be prepared for the
remainder of the project.

First, a high-level overview of our toolchain:

- **Golang**: Efficient and concurrent programming language used for building high-performance applications, server-side software, and distributed systems.
- **Kubernetes**: Automated containerized application management platform used for deploying, scaling, and managing applications in a containerized environment, providing resilience, scalability, and declarative configuration.
- **protobuf**: Efficient and language-agnostic data serialization format used for structured data exchange between different systems, optimizing both message size and processing speed.
- **gRPC**: Fast and efficient microservices communication framework that uses protobuf as the interface definition language (IDL) and enables high-performance, language-independent communication between microservices in a distributed system.

## Assignment 1: Creating a GCP instance
The first thing that we will need to do is set up a computing environment where we can eventually deploy our application. Since this is a cloud computing class, we will be using 🥁... the cloud!

### Creating a Project
We will use credits from the Google Cloud Teaching Credits Program to cover the costs of using Google Cloud Platform (GCP), so your first step is to request these credits, access your GCP account, and set up a shared project for your group.

1. **Request a Google Cloud coupon**: Use the [Student Coupon Retrieval Link](https://vector.my.salesforce-sites.com/GCPEDU?cid=hX%2B0lYOMn9PreaLGLA%2FJOXNutKSWwIktE0yiogMMYDZYfhGn9nk3r45POxAFYM%2B%2B/) to request a Google Cloud coupon. You will be asked for a name and email address, which needs to match your school domain (@ucsd.edu). A confirmation email will be sent to you with a coupon code.

2. **Access Google Cloud**: Once you have received your coupon, navigate to the [Google Cloud console](https://console.cloud.google.com/).

3. **Create a project**: One member of your team should use the [New Project](https://console.cloud.google.com/projectcreate) page to create a new GCP project for your CSE 190 labs. Choose a descriptive "Project name" and verify that the "Organization" is set to "ucsd.edu". Click `Create`.

4. **Share your project**: The creator of the new project should invite their teammates to the project. In the [Google Cloud console](https://console.cloud.google.com/) select `IAM & Admin` and then click `Grant Access`. Enter your teammate's UCSD email under "New principals", assign them the role of "Owner", and click `Save`. Repeat with your third team member if applicable. Check that all team members can now access the shared project (from the [Google Cloud console](https://console.cloud.google.com/), you can use the dropdown menu at the top left to select a project).

### Creating an Instance
Next you will set up and access a VM.

1. **Create a VM instance**: One member of your team should create a VM. From the [Google Cloud console](https://console.cloud.google.com/), click `Create a VM`. Click `Enable` to enable the Compute Engine API. Configure the following settings:

   - **Name**: Enter a name for your virtual machine.
   - **Region/Zone**: Select "us-west1 (Oregon)" as the region and leave the zone as "Any".
   - **Machine type**: Choose the preset type `e2-standard-8`. This has 8 vCPUs, 4 cores, and 32 GB of memory.
   - **Operating system and storage**: Click `OS and storage` from the left sidebar. Then click `Change` and select "Ubuntu" for the operating system and "Ubuntu 24.04 LTS" (x86/64) for the version. Select a boot disk of size 20 GB. Click `Select`.
   - **Provisioning model**: Click `Advanced` from the left sidebar. For the VM provisioning model, select "Spot". With [Spot VMs](https://cloud.google.com/solutions/spot-vms), GCP may preempt your instance under high load, but will give you a 30-second warning before doing so.

    Review the configuration and click the `Create` button to create the virtual machine. It may take a minute or so for it to be provisioned. Once it's ready, it will show up under "VM instances".

2. **Upload your SSH keys**: To access your instance via a terminal or IDE, you will need to use either an existing SSH key pair or generate a new one. You can upload a public key to your VM instance using either of these two approaches:

    - **Google Cloud UI**: Go to the instance details page (click on your VM name). Then click `Edit` towards the top of the page. Scroll down to the "Security and access" section and click `Add Item` under "SSH Keys". Enter your public key and then click `Save` at the bottom of the page.
    - **GCloud CLI**: Install the [gcloud CLI](https://cloud.google.com/sdk/docs/install#linux) and complete the setup steps from `gcloud init` if necessary. Running the following command should automatically generate and add an ssh key for you.
    
        ```console
        $ gcloud compute ssh <INSTANCE_NAME>
        ```
    Each team member will need to upload their own SSH key.

3. **SSH into the virtual machine**: From the "VM instances" page, find your instance's external IP address. Then ssh to the instance from a terminal:
    ```console
    # USERNAME: your username in the SSH key. 
    # EXTERNAL_IP: the external IP address of the VM.
	# PATH_TO_PRIVATE_KEY: the location of your private key, often ~/.ssh/<key name>
    $ ssh -i PATH_TO_PRIVATE_KEY USERNAME@EXTERNAL_IP
    ```

4. **Starting and stopping your instance**: By default your instance will keep running indefinitely. To stop your instance, click the check box to the left of it on the "VM instances" page and then click `Stop` in the menu above. All of your VM's state will be saved while it is stopped. Use the same approach to `Start/Resume` your instance later. You should stop your instance when none of your team members are using it to avoid incurring excessive costs.


## Assignment 2: Experiment Setup
Set up your VM by performing the following steps.

1. **Clone Repository**:
Clone your team's project repository. If prompted for a password, you will need to generate SSH keys and upload them to your GitHub profile (alternatively, forward existing keys from your personal computer by using the `-A` option when you SSH to your instance).
    ```console
    $ git clone <YOUR_REPO>.git
    $ cd <YOUR_REPO>
    ```

2. **Install Dependencies**:

    ***Kubernetes***: The following script installs all the key Kubernetes tools that you will need including `kubectl`, `kubeadm`, and `kubelet`. It also initializes your virtual machine as a control plane node. 
    ```console
    $ . ./scripts/k8s_setup.sh
    $ kubectl version # You should see the version of client and server if kubernetes is sucessfully installed
    ```
    Only one member of your group needs to run this script. Others can simply run the following:
    ```console
    $ mkdir -p $HOME/.kube
    $ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    $ sudo chown $(id -u):$(id -g) $HOME/.kube/config
    ```

    ***Golang***:  The following script installs Golang.
    ```console
    $ . ./scripts/golang_setup.sh
    $ go version
    ```
    After one member of your group runs the script, the other group members just need to update their path via
    ```console
    $ echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    $ source $HOME/.bashrc
    ```

    ***Conda and Python***: We will use [Python](https://python.org) for tasks such as populating our microservices application or checking whether actual application behavior matches expected behavior. [Miniconda](https://repo.anaconda.com/miniconda/) is a lightweight package manager that allows us to easily set up a Python environment. Python environments are essentially configurations of dependencies, Python libraries, and settings that provide isolated and self-contained spaces for running Python applications, ensuring consistency and avoiding conflicts between different projects or applications. The provided file `requirements.yaml` specifies all the needed dependencies and parameters for the cse190 environment. Each group member can run installation separately.
    ```console
    $ wget https://repo.anaconda.com/miniconda/Miniconda3-py310_23.3.1-0-Linux-x86_64.sh -O Miniconda.sh
    $ bash Miniconda.sh

	# This next command may take a while to run (about 10 minutes is normal)
    $ conda env create -f requirements.yaml
    $ conda activate cse190
    ```
    If you get an error about `conda` not being found when trying to create the env, you may need to run `source $HOME/.bashrc`.

3. **Configure Your VM**: For compatibility with Kubernetes, we need to revert to an older version of cgroups (the Linux component that manages containers).  To do so, one member of your group should replace the line containing `GRUB_CMDLINE_LINUX=""` in `/etc/default/grub` with:

    ```
    GRUB_CMDLINE_LINUX="systemd.unified_cgroup_hierarchy=0"
    ```

    After this, run `sudo update-grub` and then reboot your VM
    (for example, via `sudo reboot`).

Awesome! Now that we have finished setting up our VM, let's learn more about Go, Kubernetes, and gRPC!

## Assignment 3: Go, gRPC, and Kubernetes
Complete the following tutorials so that you are prepared for the remainder of the labs.

1. [A Tour of Go](https://go.dev/tour/list): A short, interactive tutorial that covers the basics of the Go programming language. You can complete this tutorial in the browser or on the virtual machine you set up in Assignment 1. 

2. [gRPC and Protobuf Tutorial](../tutorials/grpc.md): Complete this tutorial on the virtual machine you set up in Assignment 1. 

3. [Kubernetes Tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/): Complete this tutorial on the virtual machine you set up in Assignment 1. 
    - You can skip the minikube setup since we'll be using kubeadm for
      this lab.
    - `k8s_setup.sh` has already called `kubeadm` to create a
      cluster. Hence, you can skip the first step ("Create a
      Kubernetes Cluster")
    - Only the "Deploy an app" module is mandatory. However, you're encouraged to delve into other modules if you're interested.
	- When you're done with the tutorial, you can delete your
      deployment with:

```console
$ kubectl delete deployment kubernetes-bootcamp
```
