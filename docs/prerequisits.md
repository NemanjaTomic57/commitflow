# Prerequisites

To run this application in your development environment, you will need to fulfill the following prerequisites:

- Install Virtualbox
- Install Vagrant
- Install Ansible
- Install Golang
- S3 bucket in AWS
- AWS IAM user with permissions to the bucket
- GitLab Personal Access Token

## AWS IAM

Before we provision the virtual machine with Vagrant, we first have to add the credentials for the S3 bucket into /vagrant/.aws/credentials.

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
                "s3:AbortMultipartUpload",
                "s3:PutObjectTagging"
            ],
            "Resource": "arn:aws:s3:::<S3_BUCKET_NAME>/*"
        }
    ]
}
```

Don't forget to fill in the placeholder for the bucket name in the policy.

## GitLab Personal Access Token

We also need a personal access token for GitLab and GitHub APIs. Create them in your account and add them to the project by copying the environment file example and filling out the placeholders.

The personal access token needs the API permissions.

## Provision the VM with Vagrant

Before we can provision the VM, we first have to fill out a few placeholder variables in the following files.

- /.env
- /vagrant/.aws/credentials
- /vagrant/kafka/connect/s3.json

To run this application, first create and provision the virtual machine for Kafka.

```bash
cd vagrant
vagrant up
ansible-playbook -i hosts playbook.yml
```

This will start the Kafka broker on a virtual machine.

You can automate the provisioning of the AWS credentials and the virtual machine by running the `provision.sh` script in the project root. This will then ask you for the bucket name, access key id, secret access key and then fill out all the placeholders with the sed CLI tool.
