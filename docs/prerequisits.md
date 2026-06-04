# Prerequisites

To run this application in your development environment, you will need to fulfill the following prerequisites:

First, git pull the project.

- Install Virtualbox
- Install Vagrant
- Install Ansible
- pip install jmespath
- Install Golang
- S3 bucket in AWS
- AWS IAM user with permissions to the bucket
- GitHub Personal Access Token
- GitLab Personal Access Token

Ansible requires you to install 

## AWS IAM

Before we provision the virtual machine with Vagrant, we first have to add the credentials for the S3 bucket into the .env file.

For this, follow the following steps:

1. Create a bucket in your account.
2. Create an IAM user with access keys. Note down the access keys, you will need them later.
3. Assign the following policy to the IAM user.

```json
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"s3:ListAllMyBuckets"
			],
			"Resource": "arn:aws:s3:::*"
		},
		{
			"Effect": "Allow",
			"Action": [
				"s3:ListBucket",
				"s3:GetBucketLocation"
			],
			"Resource": "arn:aws:s3:::<S3_BUCKET_NAME>"
		},
		{
			"Effect": "Allow",
			"Action": [
				"s3:PutObject",
				"s3:GetObject",
				"s3:DeleteObject",
				"s3:AbortMultipartUpload",
				"s3:PutObjectTagging"
			],
			"Resource": "arn:aws:s3:::<S3_BUCKET_NAME>/*"
		}
	]
}
```

Don't forget to fill in the placeholder for the bucket name in the policy.

After you created the bucket and the IAM user with the access keys, copy the access keys as well as the bucket name into the .env file. You can use the template at .env.example.

This bucket will be used to sink the data from Git to S3. The IAM access keys will be used by the Kafka Connect S3 Sink plugin to connect and upload the data to JSON objects.

The costs for uploading and storing the data in S3 depend on the size of your Git accounts, based on how many projects you **own**. Projects you maintain but do not own are not being considered and no data will be fetched for those projects.

## GitHub and GitLab Personal Access Token

We also need a personal access token for GitLab and GitHub APIs. Create them in your account and add them to the project by copying the environment file example and filling out the placeholders.

GitLab PAT permissions need to include:

- api

GitHub PAT (classic) permsissions need to include:

- repo


## Provision the VM with Vagrant

Before we can provision the VM, we first have to fill out a few placeholder variables in the following files.

- /.env

To run the application, first create and provision the virtual machine for Kafka.

```bash
cd vagrant
vagrant up
ansible-playbook -i hosts playbook.yml
```

This will start the Kafka broker on a virtual machine.

You can automate the provisioning of the environemnt variables and the virtual machine by running the `provision.sh` script in the project root.

## PostgreSQL Database

If you want to connect to the PostgreSQL database, you can do so from within the VM.

```bash
sudo -U postgres psql
```

## Grafana Dashboard

The Grafana dashboard is available under http://192.168.56.10:3000.

Default credentials:
username: admin
password: admin

You have to change the password at first login.
