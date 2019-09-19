# Foghorn 📯

Foghorn is a fully Kubernetes-native CI system, inspired by [prow](https://github.com/kubernetes/test-infra/tree/master/prow). It's powered by [Tekton Pipelines](https://github.com/tektoncd/pipeline), which provides a number of handy abstractions for defining CI pipelines in a k8s-native way. Foghorn provides webhook handling and some git-centric CommentOps features, and it supports multiple git providers.

## Architecture

There are two logical groupings of components in Foghorn:

* Catcher
  - A highly-available webhook handler. Analogous to `hook` in Prow.
* Controllers
  - Controllers for Foghorn's CRDs, `GitEvent` and `Action`
  - contains the bulk of Foghorn's business logic

The way this works is that Catcher receives a webhook from a git provider (GitHub, Gitlab, Bitbucket, etc.) and creates a `GitEvent` resource, which contains a webhook that has been parsed into a generic struct using `go-scm`.

The `GitEvent` controller then picks up the newly created `GitEvent` and determines what should happen next, e.g., run a Tekton pipeline, label an issue on the git provider, etc. It then creates an `Action` resource that will encapsulate data about the required action, at which point the `Action` controller takes over and monitors the action, recording any changes to its state.

Catcher is designed to run in a dedicated binary, while the controllers can be run either as a single, monolithic binary or as separate, dedicated binaries. By running them separately, you have the option of replacing one of the provided controllers with a controller of your own if you wish to override the default behavior of that controller. You can also add a controller to extend Foghorn's capabilities in either scenario.

## Building

The fastest way to get started hacking on Foghorn is with [ko](https://github.com/google/ko).

If you take a look in `deployments`, everything should look normal except for the following field in `200-catcher-deployment.yaml`:

```yaml
image: github.com/foghornci/foghorn/cmd/catcher
```

The magic of `ko` is that it will: 

* automatically build the go binary whose code is specified in the `image` field
* package it into an image
* push the image to a configured image registry
* update the deployment manifest with the newly build image name and tag
* `kubectl apply` a collection of k8s manifest files

Once you've installed and configured `ko`, simply run the following from the project root:

```sh
ko apply -f deployments/
```

This will build and deploy Foghorn to your Kubernetes cluster. In addition to a deployment, it will create a service and ingress so that you can start receiving webhooks right away, and it will put all of it in a namespace called `foghorn`. Simply re-run the command to build and deploy your latest changes. It's idempotent, so running it multiple times with no changes in between will have no effect.
