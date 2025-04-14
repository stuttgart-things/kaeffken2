schema VsphereVmAnsible:
    """VsphereVmAnsible defines a Vsphere VM with Ansible provisioning."""
    apiVersion: str = "resources.stuttgart-things.com/v1alpha1"
    kind: str = "VsphereVmAnsible"
    metadata: {str:str}
    spec: VsphereVmAnsibleSpec

schema VsphereVmAnsibleSpec:
    vm: VMSpec
    vsphere: VsphereSpec
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
    firmware: "bios" | "uefi"
    template: str
    bootstrap?: str
    annotation?: str
    unverifiedSsl?: str

schema VsphereSpec:
    folderPath: str
    datacenter: str
    datastore: str
    resourcePool: str
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
    roles?: [str]
    collections: [str]