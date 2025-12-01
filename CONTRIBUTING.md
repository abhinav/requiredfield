# Contributing

## (For maintainers) Releasing

To release a new version of requiredfield:

1. Batch the changes into a version and merge the changelog.

   ```bash
   changie batch minor && changie merge
   ```

   Replace `minor` with `major`, `patch`, or a specific version as appropriate.

2. Create a release branch and commit the changes.

   ```bash
   git add .changes CHANGELOG.md
   git checkout -b release-$(changie latest)
   git commit -m "Release $(changie latest)"
   ```

3. Push the branch and create a PR.

   ```bash
   git push -u origin release-$(changie latest)
   gh pr create --fill
   ```

4. Wait for CI to pass, then merge the PR.

   ```bash
   gh pr merge <number> --squash
   ```

5. Switch back to main and pull the changes.

   ```bash
   git checkout main
   git pull
   ```

6. Create the GitHub release.

   ```bash
   gh release create $(changie latest) -F .changes/$(changie latest).md
   ```
