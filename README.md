# Reproduction of k8s issue #74839

## Howto

1. Create a cluster with at least 2 nodes.
1. Deploy the app.
```console
kubectl create -f ./deploy.yaml
```

1. Check if the server crashed
```console
$ kubectl get pods
NAME                           READY   STATUS             RESTARTS   AGE
boom-server-59945555cd-8rwqk   0/1     CrashLoopBackOff   4          2m
startup-script                 1/1     Running            0          2m
```

## Reference

https://github.com/kubernetes/kubernetes/issues/74839


