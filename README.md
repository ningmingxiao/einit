# einit 
if a process write in go run in container as 1 pid(process A). If other process (process B) exited sometimes process 1 will be chosed to be it's ppid. If this B process is Zombie, process A will not reap Zombie process
