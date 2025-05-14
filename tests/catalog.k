import schemas.component_schema

# Configuration
_name = option("name") or "" #ask;min2+max30
_description = option("description") or "" #ask;min2+max30

component_schema.BackstageComponent {
    metadata = {
        name = _name
        description = _description
        tags = ["proxmox", "vm", "ansible"]
        annotations = {
            "backstage.io/managed-by-location": "url:https://github.com/example/repo/blob/main/catalog.yaml",
            "backstage.io/techdocs-ref": "dir:."
        }
        namespace = "default"
    }
    spec = {
        type = "infrastructure"
        lifecycle = "production"
        owner = "platform-team"
        system = "vm-management"
        dependsOn = ["resource:proxmox-vm-01"]
    }
}
