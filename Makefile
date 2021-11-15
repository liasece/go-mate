
.PHONY: all
all: run


.PHONY: run
run:
	go run main/main.go buildRunner -f /home/user/testgame.go -n e1 -n e2

