# Color settings
# See: https://gist.github.com/rsperl/d2dfe88a520968fbc1f49db0a29345b9
# define standard colors
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)
RESET := $(shell tput -Txterm sgr0)

B=feature
T=$(B)
pr:
	@echo
	@echo "${GREEN}# DELETE MERGED BRANCH ${RESET}"
	-git delete-merged-branch
	@echo
	@echo "${GREEN}# UPDATE MAIN BRANCH ${RESET}"
	git pull origin develop:develop
	@echo
	@echo "${GREEN}# CREATE NEW BRANCH ${RESET}"
	-git co -b $(B)
	@echo
	@echo "${GREEN}# CREATE PULL REQUEST ${RESET}"
	git commit --allow-empty -m ":tada: The first commit in $(B)"
	gh pr create -a @me -t "[PR] $(T)" -B develop
	gh pr view --web

doc:
	@echo "open in web browser http://localhost:6060"
	@godoc -http=:6060
