# Release Flow Version

The idea of Release Flow Version is to generate SemVer version numbers based on the Release Flow.

## Release Flow

The Release Flow follows these rules:

* Trunk based development on main branch
* Feature branches are created off main and merged into main
* Release branches are created off main and never merged
* Hotfix branches are created off main, merged into main and cherry picked into release branch
* If cherry picking is not possible or desired then and separate hotfix branch can be created off release branch and merged into release branch
* Release branch have the release number in the branch name (release/1.0.0)
* Each commit (merged hotfix) on a release branch increases the patch version (1.0.1, 1.0.2, ...)
* The version on the main branch uses the highest release version as base and increases the minor version. It will also append a prerelease tag and the number of commits since the release branch was created (1.1.0-beta.5). After release/1.1.0 is created the version for new builds on main will be 1.2.0-beta.1, 1.2.0-beta.2, ...
* The version on feature branches uses the highest release version as base and increases the minor version. It will also append the branch name and the number of commits since the branch was created (1.1.0-newfeature.2)