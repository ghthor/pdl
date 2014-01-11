DIRS = mysql

EXAMPLE_APP = example
EXAMPLE_APP_FILES = $(shell find $(EXAMPLE_APP).app/*.js) $(shell find $(EXAMPLE_APP).app/*.html)

test: test-deps
	go test -run=UnitSpecs -rebuild-deps=false

integration-test: test-deps
	go test -run=IntegrationSpecs -rebuild-deps=false

test-deps: mysql/schema.sql $(EXAMPLE_APP).pkg

mysql/%.sql : force_look
	$(MAKE) -C mysql $(subst mysql/,,$@) $(MFLAGS)

$(EXAMPLE_APP).pkg : $(EXAMPLE_APP_FILES)
	zip -r $@ $(EXAMPLE_APP).app

clean :
	@for d in $(DIRS) ; do \
		$(MAKE) -C $$d clean $(MFLAGS) ; \
	done
	go clean ./...
	rm -f $(EXAMPLE_APP).pkg

force_look :
	@true

.PHONY : test integration-test test-deps clean
