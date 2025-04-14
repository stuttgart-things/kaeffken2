import schemas.vsphere_vm_ansible

# Configuration with interactive options
_count = option("count") or "1" #select;1,2,3,4,5
_name = option("name") or "" #ask;min2+max30
_cpu = option("cpu") or "8" #select;4,8,12,16
_ram = option("ram") or "4096" #select;4096,8192,12288
_disk = option("disk") or "64G" #select;32G,64G,96G,128G,196G,256G
_template = option("template") or "ubuntu24" #select;ubuntu24,rocky9
_annotation = option("annotation") or "build w/ crossplane" #ask;min2+max30
_namespace = option("namespace") or "proxmox" #ask;min2+max30
_firmware = option("firmware") or "bios"
_bootstrap = option("bootstrap") or "['echo STUTTGART-THINGS']"
_unverifiedSsl = option("unverifiedSsl") or "true"
_datastore = option("datastore") or "v3700"
_folderPath = option("folderPath") or "stuttgart-things"
_network = option("network") or "vmbr0" #select;vmbr101,vmbr102,vmbr103,vmbr104
_datacenter = option("datacenter") or "/LabUL"
_resourcePool = option("resourcePool") or "/LabUL/host/Cluster-V6.5/Resources"



# Proxmox VM Configuration
vsphere_vm_ansible.VsphereVmAnsible {
    metadata = {
        name = _name
        namespace = _namespace
    }
    spec = {
        vm = {
            count = _count
            name = _name
            cpu = _cpu
            ram = _ram
            disk = _disk
            firmware = _firmware
            template = _template
            bootstrap = _bootstrap
            unverifiedSsl = _unverifiedSsl
            annotation = _annotation
        }
        vsphere = {
            folderPath = _folderPath
            datacenter = _datacenter
            datastore = _datastore
            resourcePool = _resourcePool
            network = _network
        }
        tfvars = {
            secretName = "proxmox-tfvars"
            secretNamespace = "proxmox"
            secretKey = "terraform.tfvars"
        }
        connectionSecret = {
            name = _name
            namespace = _namespace
        }
        compositionRef = {
            name = "proxmox-vm-ansible"
        }
        providerRef = {
            name = "proxmox-vm"
        }
        ansible = {
            provisioningName = _name
            playbooks = [
                "sthings.baseos.prepare_env",
                "sthings.baseos.setup",
                "sthings.baseos.dev"
            ]
            ansibleVarsFile = [
                "golang_version+-1.24.1",
                "manage_filesystem+-true",
                "update_packages+-true",
                "install_requirements+-true",
                "install_motd+-true",
                "username+-sthings",
                "lvm_home_sizing+-'15%'",
                "lvm_root_sizing+-'35%'",
                "lvm_var_sizing+-'50%'",
                "event_author+-crossplane",
                "event_tags+-ansible,baseos,crossplane,tekton",
                "send_to_msteams+-true",
                "reboot_all+-false"
            ]
            gitRepoUrl = "https://github.com/stuttgart-things/ansible.git"
            gitRevision = "main"
            providerRef = {
                name = "kubernetes-incluster"
            }
            vaultSecretName = "vault"
            pipelineNamespace = "tekton-pipelines"
            workingImage = "ghcr.io/stuttgart-things/sthings-ansible:11.3.0-3"
            collections = [
                "https://github.com/stuttgart-things/ansible/releases/download/sthings-rke-25.3.610/sthings-rke-25.3.610.tar.gz",
                "https://github.com/stuttgart-things/ansible/releases/download/sthings-awx-25.5.617/sthings-awx-25.5.617.tar.gz",
                "https://github.com/stuttgart-things/ansible/releases/download/sthings-container-25.1.813.tar.gz/sthings-container-25.1.813.tar.gz",
                "https://github.com/stuttgart-things/ansible/releases/download/sthings-baseos-25.1.814.tar.gz/sthings-baseos-25.1.814.tar.gz"
            ]
        }
    }
}