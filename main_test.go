package main

import "testing"

func TestUpdateContents(t *testing.T) {
	img := "snarlysodboxer/my-image"
	tag := "1.2.3"
	input := `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-stuff
spec:
  template:
    spec:
      initContainers:
      - name: test1
        image: snarlysodboxer/my-image:1.19-v0.0.1
      - name: test2
        image: snarlysodboxer/my-image:0.1.209  # this is a comment that should stay
      containers:
      - name: test33
        image: snarlysodboxer/my-image:7.5.10-v0.0.21`
	expected := `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-stuff
spec:
  template:
    spec:
      initContainers:
      - name: test1
        image: snarlysodboxer/my-image:1.2.3
      - name: test2
        image: snarlysodboxer/my-image:1.2.3  # this is a comment that should stay
      containers:
      - name: test33
        image: snarlysodboxer/my-image:1.2.3`

	updated := updateContents([]byte(input), &img, &tag)
	if string(updated) != expected {
		t.Fatalf("Expected: \n%s\nGot: \n%s\n", expected, updated)
	}
}
