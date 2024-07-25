# Releasing

## Steps to Release a New Version

- Make sure your local copy of the repository is up to date:

  ```shell
  git pull --rebase
  ```

- Verify CI tests:

  ```shell
  make fmtcheck errcheck lint
  ```

- Verify build and tests:

  ```shell
  make build test
  ```

- Tag with annotation:

  ```shell
  git tag -a v0.0.0 -m "Release v0.0.0"
  ```

- Verify the tag:

  ```shell
  git tag -l -n
  ```

- Push the tag:

  ```shell
  git push --tags
  ```

- A [Release Github Action](https://github.com/vmware/terraform-provider-vra/actions/workflows/release.yml) will automatically be triggered, and a [draft release](https://github.com/vmware/terraform-provider-vra/releases) will be created.
