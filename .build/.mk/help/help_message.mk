.PHONY: help
help: help_message

help_message:
	@awk 'BEGIN {FS = ":.*##"; \
		printf "$(<b>)$(<cYellow>)Usage:$(</>)$(</>) \
		$(<br>)make $(<cGreen>)<target>$(</>) $(<cBlue>)\"<arguments>\"$(</>)️\
		$(<br>)$(<br>)$(<b>)$(<cYellow>)Available commands:$(</>)$(</>) \
		$(<br>)"}/^[a-zA-Z_-]+:.*?##/{ printf "  $(<cGreen>)%-10s$(</>) %s$(<br>)", $$1, $$2 } /^##@/{ printf "\
		$(<br>)$(<b>)%s$(</>)$(<br>)", substr($$0, 5) }'$(MAKEFILE_LIST)

