TYPE=default

.SILENT:

all:
	go run main.go -type=$(TYPE)

wq:
	go run main.go -type=workerqueue

dr:
	go run main.go -type=durable