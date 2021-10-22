# Releasing

Steps to Release a New Version
------------------------------

- Make sure your local copy of the repository is up to date:

	```
	git pull --rebase
	```

- Verify CI tests:

	```
	make fmtcheck errcheck lint
	```

- Verify build and tests:

	```
	make build test
	```

- Tag with annotation:

	```
	git tag -a v0.0.0 -m "Release v0.0.0"
	```

- Verify the tag:

	```
	git tag -l -n
	```

- Push the tag:

	```
	git push --tags
	```

- A [Release Github Action](https://github.com/vmware/terraform-provider-vra/actions/workflows/release.yml) will automatically be triggered, and a [draft release](https://github.com/vmware/terraform-provider-vra/releases) will be created.
