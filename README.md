# kaeffken2
gitops scaffolder

## USAGE

<details><summary>RENDER REQUEST INTERACTIVE</summary>

```bash
kaeffken2 render \
--config tests/vmRequestConfig.yaml \
--request tests/vmRequest.yaml
```

</details>


<details><summary>RENDER NON-INTERACTIVE (w/ RANDOM VALUES)</summary>

```bash
kaeffken2 render \
--config tests/vmRequestConfig.yaml \
--request tests/vmRequest.yaml \
--survey=false \
--destination /tmp/output2.yaml
```

</details>

## TESTS

<details><summary>RENDER w/ KCL CLI AND PARAMETERS</summary>

```bash
kcl catalog.k -D name="vm-02"
```

</details>
