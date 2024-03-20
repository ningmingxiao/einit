# einit 
if a process write in go run in container as 1 pid(process A). If other process (process B) exited sometimes process 1 will be chosed to be it's ppid. If this B process is Zombie, process A will not reap Zombie process
for example  k8s yaml
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sx-heap
  labels:
    app: test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      hostNetwork: true
      containers:
        - name: busybox-1
          image: docker.io/library/alpine:latest
          imagePullPolicy: IfNotPresent
          command: ['/bin/sh', '-ec','top']
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "sleep 201 &"]
```
kill -9 `pidof sleep`
sleep 201 will become Zombie.
If set shareProcessNamespace: true pause will   reap Zombie process,but sometimes is not safe enough.
