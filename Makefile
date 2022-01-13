
.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/ github.com/liasece/go-mate/main
	@ rm -rf /Users/jansen/lilith/server/backendv2/cluster/gamedev/internal/domain/gameEntryOpt.go && \
		./bin/main buildRunner -f /Users/jansen/lilith/server/backendv2/cluster/gamedev/internal/domain/entity/gameEntry.go -o /Users/jansen/lilith/server/backendv2/cluster/gamedev/internal/domain/gameEntryOpt.go -n GameEntry
	@ rm -rf /Users/jansen/lilith/server/backendv2/cluster/social/internal/domain/entity/userRelationOpt.go && \
		./bin/main buildRunner -f /Users/jansen/lilith/server/backendv2/cluster/social/internal/domain/entity/userRelation.go -o /Users/jansen/lilith/server/backendv2/cluster/social/internal/domain/entity/userRelationOpt.go -n UserRelation
	@ rm -rf /Users/jansen/lilith/server/backendv2/cluster/gamedev/internal/domain/gameDetailOpt.go && \
		./bin/main buildRunner \
			-f /Users/jansen/lilith/server/backendv2/cluster/gamedev/internal/domain/entity/ \
			--entity-pkg solarland/backendv2/cluster/gamedev/internal/domain/entity \
			-n GameDetail \
			-o /Users/jansen/lilith/server/backendv2/cluster/gamedev/internal/domain/gameDetailOpt.go \
			--out-rep-inf-file /Users/jansen/lilith/server/backendv2/cluster/gamedev/internal/domain/gameDetailBase.go \
			--out-rep-adp-file /Users/jansen/lilith/server/backendv2/cluster/gamedev/internal/adapter/repository/game/detail/baseGameDetail.go

.PHONY: run
run:
	go run main/main.go buildRunner -f /home/user/testgame.go -n e1 -n e2

