# Releasing

Steps to release a new version
------------------------------

- Make sure repo is up to date: `git pull`
- Verify CI tests: `make fmtcheck errcheck lint`
- Verify build and tests: `make build test`
- Update CHANGELOG.md (TBD)
- Tag with annotation: `git tag -a -m "Release v0.0.0" v0.0.0`
- Verify tag: `git tag -n`
- Push tag: `git push --tags`
- Build and release binaries: `./scripts/release.sh`

Helper to create CHANGELOG entries
----------------------------------

`git log --reverse --pretty=format:"%s" | tail -100 | sed 's/^/* /'`
