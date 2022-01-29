# CC Controller

## Why this?

This is an implementation of a controller using only [client-go](https://github.com/kubernetes/client-go). The `CC` controller code is based on [sample-controller](https://github.com/kubernetes/sample-controller).

The idea is to play around with different ways for defining controllers. This implementation is the approach where more tuning is required in order to make the controller work. But it's also a the best approach for understanding the controller internals.

Other options for defining a k8s controller are:

- Via [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime).
- Via [kubebuilder](https://book.kubebuilder.io/introduction.html)

## What is this controller about?

This controller reconciles on events(_create,update_) of [BuildRun](https://github.com/shipwright-io/build/blob/main/samples/buildrun/buildrun_buildah_cr.yaml) CRD's. The controller doesn't do too much, it just logs some messages for educational purposes.

## How to use this?

- A Kubernetes cluster, e.g. [kind](https://github.com/shipwright-io/build/blob/main/hack/install-kind.sh).
- Install Tekton:

  ```sh
  kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.30.0/release.yaml
  ```

- Install Shipwright

  ```sh
  kubectl apply --filename https://github.com/shipwright-io/build/releases/download/v0.7.0/release.yaml
  kubectl apply --filename https://github.com/shipwright-io/build/releases/download/v0.7.0/sample-strategies.yaml
  ```

- Build the go binary

  ```sh
  go build -o cc-controller cmd/cc-controller/main.go
  ```

- Run the binary

  ```sh
  ./cc-controller
  ```

- Create the following two CRD's, to see the controller reconciling:

   ```sh
   kubectl apply -f https://raw.githubusercontent.com/shipwright-io/build/main/samples/build/build_buildah_cr.yaml
   kubectl apply -f https://raw.githubusercontent.com/shipwright-io/build/main/samples/buildrun/buildrun_buildah_cr.yaml
   ```
