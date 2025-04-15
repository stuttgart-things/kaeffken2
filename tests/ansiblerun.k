import schemas.ansible_run

# Configuration
_name = option("name") or "ansiblerun" #ask;min2+max35
_namespace = option("namespace") or "crossplane-system" #ask;min2+max35
_vault_secret = option("vault_secret") or "vault2" #ask;min2+max35
_createInventory = option("createInventory") or "false" #select;true,false
_workingImage = option("workingImage") or "ghcr.io/stuttgart-things/sthings-ansible:11.3.0-3" #ask;min2+max60
_pipelineRunNamespace = option("pipelineRunNamespace") or "tekton-pipelines" #ask;min2+max35
_ansible_vars = option("ansible_vars") or """
manage_filesystem+-true
update_packages+-true
""" #list

ansible_run.AnsibleRun {
    metadata = {
        name = _name
        namespace = _namespace
    }
    spec = {
        pipelineRunName = _name
        createInventory = "false"
        varsFile = "bmFtZToga29sbGUK"
        inventoryFile = "W2luaXRpYWxfbWFzdGVyX25vZGVdCjEwLjMxLjEwMy40MwoKW2FkZGl0aW9uYWxfbWFzdGVyX25vZGVzXQo="
        playbooks = [var for var in _playbooks.split("\n") if var]
        ansibleVarsFile = [var for var in _ansible_vars.split("\n") if var]
        gitRepoUrl = "https://github.com/stuttgart-things/ansible.git"
        gitRevision = "main"
        providerRef = {
            name = "kubernetes-incluster"
        }
        vaultSecretName = _vault_secret
        pipelineNamespace = _pipelineRunNamespace
        workingImage = _workingImage
        roles = [
            "https://github.com/stuttgart-things/install-requirements.git,2024.05.11"
        ]
        collections = [
            "community.crypto:2.22.3",
            "community.general:10.1.0",
            "ansible.posix:2.0.0",
            "kubernetes.core:5.0.0",
            "community.docker:4.1.0",
            "community.vmware:5.2.0",
            "awx.awx:24.6.1",
            "community.hashi_vault:6.2.0",
            "ansible.netcommon:7.1.0",
            "https://github.com/stuttgart-things/ansible/releases/download/sthings-baseos-25.3.1202.tar.gz/sthings-baseos-25.3.1202.tar.gz",
            "https://github.com/stuttgart-things/ansible/releases/download/sthings-container-25.6.1311.tar.gz/sthings-container-25.6.1311.tar.gz"
        ]
    }
}