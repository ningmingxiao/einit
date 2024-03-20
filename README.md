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
      containers:
        - name: busybox-1
          image: docker.io/library/alpine:latest
          imagePullPolicy: IfNotPresent
          command: ['/home/main']
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "sleep 201 &"]
```

main.go
```
package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("sleep", "10000")

	// 启动命令
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 10000)

}
```
kill -9 `pidof sleep`
sleep 201 will become Zombie.
If set shareProcessNamespace: true pause will   reap Zombie process,but sometimes it is not safe enough.

```
root      531129  3.0  0.1 722524 13908 ?        Sl   18:39   0:00 /usr/bin/containerd-shim-runc-v2 -namespace k8s.io -id bdd8d75f2b7ffe2b48d9df6a55068527e0d270ff760b317318831d
65535     531151  3.0  0.0    996     4 ?        Ss   18:39   0:00  \_ /pause
root      531208  0.0  0.0 710244  2944 ?        Ssl  18:39   0:00  \_ /home/main
root      531225  0.0  0.0   1612     4 ?        S    18:39   0:00      \_ sleep 10000
root      531242  0.0  0.0   1612     4 ?        S    18:39   0:00      \_ sleep 201
```

kill -9 531242 (531242 will not be reaped by main)

```
root      531129  0.0  0.1 722524 13544 ?        Sl   18:39   0:00 /usr/bin/containerd-shim-runc-v2 -namespace k8s.io -id bdd8d75f2b7ffe2b48d9df6a55068527e0d270ff760b317318831d
65535     531151  0.0  0.0    996     4 ?        Ss   18:39   0:00  \_ /pause
root      531208  0.0  0.0 710244  2944 ?        Ssl  18:39   0:00  \_ /home/main
root      531225  0.0  0.0   1612     4 ?        S    18:39   0:00      \_ sleep 10000
root      531242  0.0  0.0      0     0 ?        Z    18:39   0:00      \_ [sleep] <defunct>

```
kill -9 531225 (531225 will be reaped by main)
```
root      531129  0.0  0.1 722524 14808 ?        Sl   18:39   0:00 /usr/bin/containerd-shim-runc-v2 -namespace k8s.io -id bdd8d75f2b7ffe2b48d9df6a55068527e0d270ff760b317318831d
65535     531151  0.0  0.0    996     4 ?        Ss   18:39   0:00  \_ /pause
root      531208  0.0  0.0 710244  2944 ?        Ssl  18:39   0:00  \_ /home/main
root      531242  0.0  0.0      0     0 ?        Z    18:39   0:00      \_ [sleep] <defunct>

```
