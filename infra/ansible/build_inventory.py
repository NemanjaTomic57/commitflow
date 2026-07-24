import json
import subprocess

import yaml

# Default SSH user used by Ansible to connect to all hosts.
ANSIBLE_USER = "admin"


def terraform_output():
    """
    Retrieve Terraform outputs as a Python dictionary.

    Executes:
        terraform output -json

    Returns:
        dict: Terraform outputs in JSON format.
    """
    result = subprocess.run(
        ["terraform", "output", "-json"],
        capture_output=True,
        text=True,
        cwd="../terraform/base",
        check=True,
    )

    return json.loads(result.stdout)


def get_value(state, key):
    """
    Extract a value from the Terraform output.

    If the output is a string, it is returned directly.
    If the output is a map/object, its values are returned as a list.

    Args:
        state (dict): Terraform output dictionary.
        key (str): Name of the Terraform output.

    Returns:
        str | list: The extracted value.
    """
    value = state[key]["value"]

    # Return primitive string outputs directly.
    if isinstance(value, str):
        return value

    # Convert map/object values into a list.
    values = []

    for item in value.values():
        values.append(item)

    return values


def build_inventory(tf):
    """
    Build an Ansible inventory from Terraform outputs.

    The generated inventory contains:
      - Global Ansible variables.
      - SSH ProxyCommand configuration using the bastion host.
      - A 'kafka' host group populated with all Kafka nodes.

    Args:
        tf (dict): Terraform output dictionary.

    Returns:
        dict: Ansible inventory structure.
    """
    bastion_ip = tf["bastion_public_ip"]["value"]
    kafka_nodes = tf["kafka_private_ips"]["value"]

    inventory = {
        "all": {
            "vars": {
                "ansible_user": ANSIBLE_USER,
                "ansible_ssh_private_key_file": "~/.ssh/aws.pem",

                # Route all SSH connections through the bastion host.
                "ansible_ssh_common_args": (
                    "-o StrictHostKeyChecking=no "
                    "-o ProxyCommand=\"ssh "
                    "-i ~/.ssh/aws.pem "
                    "-o StrictHostKeyChecking=no "
                    f"-W %h:%p {ANSIBLE_USER}@{bastion_ip}\""
                ),
            },
            "children": {
                "kafka": {
                    "hosts": {}
                }
            },
        }
    }

    # Add every Kafka node to the inventory.
    for name, ip in kafka_nodes.items():
        inventory["all"]["children"]["kafka"]["hosts"][name] = {
            "ansible_host": ip
        }

    return inventory


def main():
    """
    Generate an Ansible inventory from Terraform outputs.

    The inventory is written to 'inventory.yml' in the current directory.
    """
    # Read Terraform outputs.
    tf = terraform_output()

    # Convert Terraform outputs into an Ansible inventory.
    inventory = build_inventory(tf)

    # Write the inventory to disk.
    with open("inventory.yml", "w") as f:
        yaml.safe_dump(inventory, f, sort_keys=False)

    print("Generated inventory.yml")


if __name__ == "__main__":
    main()

