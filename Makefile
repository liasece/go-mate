
.PHONY: all
all: run


.PHONY: run
run:
	go run main/main.go buildRunner -f R:\lilithgames\Avatar\Server\src\solarland\backendv2\cluster\gamedev\domain\entity\game.go -n GameDetail1 -n GameDetail2

