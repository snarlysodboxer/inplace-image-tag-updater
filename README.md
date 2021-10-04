# Inplace Image Tag Updater

Update image tags in Kubernetes Specs.
Don't touch anything else, leave comments and indentation and everything else.

## Why
* The different `yq`s will either rearrange things or remove your comments or change your styling.
* `sed`/`awk` approaches are error prone and hard to read and maintain.
* `kustomize edit set images` also rearranges and restyles your `kustomization.yaml` files.

## Usage
```bash
go build
inplace-image-tag-updater -h
# assuming you have Kubernetes specs in a directory called kustomize
git grep -l snarlysodboxer/my-image: kustomize | inplace-image-tag-updater -image snarlysodboxer/my-image -newTag 1.2.3
# OR in Docker
git grep -l snarlysodboxer/my-image: kustomize | docker run -i --rm -v $(pwd):/code snarlysodboxer/inplace-image-tag-updater:v0.0.1 -image snarlysodboxer/my-image -newTag 1.2.3
```

## Customize the search regex and replacement string (Generally not needed)
* The value of the `-image` flag will have any backslashes escaped, and then be substituted into the `%s` in the `-searchRegex` flag's value. The `-searchRegex` flag's default value is `image:\\s+%s:\\S+`. (Note that the backslashes are escaped for the shell, the actual regex ends up bing `image:\s+%s:\S+` when passed to Go's `regexp.MustCompile`.)
* The values of the `-image` and `-newTag` flags will be substituted into the `%s`es in the `-replacementFormat` flag's value. The `-replacementFormat` flag's default value is `image: %s:%s`.
