# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "debian-13"

  config.vm.network("private_network", ip: "192.168.56.10")

  config.vm.provider("virtualbox") do |v|
    v.memory = 4096
    v.cpus = 4
  end
end
