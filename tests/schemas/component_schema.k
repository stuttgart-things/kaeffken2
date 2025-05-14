schema BackstageComponent:
    """BackstageComponent defines a Backstage catalog component."""
    apiVersion: str = "backstage.io/v1alpha1"
    kind: str = "Component"
    metadata: Metadata
    spec: ComponentSpec

schema Metadata:
    name: str
    description?: str
    tags?: [str]
    annotations?: {str: str}
    labels?: {str: str}
    namespace?: str

schema ComponentSpec:
    type: str
    lifecycle: str
    owner: str
    system?: str
    providesApis?: [str]
    consumesApis?: [str]
    dependsOn?: [str]
    subcomponentOf?: str
