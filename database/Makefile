DIRS = mysql

EXAMPLE_APP = example
EXAMPLE_APP_FILES = $(shell find $(EXAMPLE_APP).app/*.js) $(shell find $(EXAMPLE_APP).app/*.html)

test:
	go test -run=UnitSpecs

integration-test: mysql/schema.sql
	go test -run=IntegrationSpecs

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

.PHONY : test integration-test clean
