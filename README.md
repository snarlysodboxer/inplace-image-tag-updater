# Inplace Image Tag Updater

Update image tags in Kubernetes Specs.
Don't touch anything else, leave comments and indentation and everything else.

## Why
* The different `yq`s will either rearrange things or remove your comments or change your styling.
* `sed`/`awk` approaches are error prone and hard to read and maintain.

`go build`
`git grep -l snarlysodboxer/my-image: kustomize | ./inplace-image-tag-updater -image snarlysodboxer/my-image -newTag 1.2.3`
