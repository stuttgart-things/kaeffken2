schema ProxmoxVMAnsible:
    """ProxmoxVMAnsible defines a Proxmox VM with Ansible provisioning."""
    apiVersion: str = "resources.stuttgart-things.com/v1alpha1"
    kind: str = "ProxmoxVMAnsible"
    metadata: {str:str}
    spec: ProxmoxVMAnsibleSpec

schema ProxmoxVMAnsibleSpec:
    vm: VMSpec
    proxmox: ProxmoxSpec
    tfvars: TFVarsSpec
    connectionSecret: ConnectionSecretSpec
    compositionRef: CompositionRefSpec
    providerRef: ProviderRefSpec
    ansible: AnsibleSpec

schema VMSpec:
    count: str
    name: str
    cpu: str
    ram: str
    disk: str
    template: str
    annotation?: str

schema ProxmoxSpec:
    node: str
    datastore: str
    folderPath: str
    network: str

schema TFVarsSpec:
    secretName: str
    secretNamespace: str
    secretKey: str

schema ConnectionSecretSpec:
    name: str
    namespace: str

schema CompositionRefSpec:
    name: str

schema ProviderRefSpec:
    name: str

schema AnsibleSpec:
    provisioningName: str
    playbooks: [str]
    ansibleVarsFile: [str]
    gitRepoUrl: str
    gitRevision: str
    providerRef: {str:str}
    vaultSecretName: str
    pipelineNamespace: str
    workingImage: str
    collections: [str]