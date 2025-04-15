schema AnsibleRun:
    """AnsibleRun defines an Ansible execution configuration."""
    apiVersion: str = "resources.stuttgart-things.com/v1alpha1"
    kind: str = "AnsibleRun"
    metadata: {str:str}
    spec: AnsibleRunSpec

schema AnsibleRunSpec:
    pipelineRunName: str
    createInventory: str
    varsFile?: str
    inventoryFile?: str
    playbooks: [str]
    ansibleVarsFile?: [str]
    gitRepoUrl: str
    gitRevision: str
    providerRef: {str:str}
    vaultSecretName: str
    pipelineNamespace: str
    workingImage: str
    roles?: [str]
    collections?: [str]