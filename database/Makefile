DIRS = mysql

EXAMPLE_APP = example
EXAMPLE_APP_FILES = $(shell find $(EXAMPLE_APP).app/*.js) $(shell find $(EXAMPLE_APP).app/*.html)

test: test-deps
	go test -run=UnitSpecs -rebuild-deps=false

integration-test: test-deps
	go test -run=IntegrationSpecs -rebuild-deps=false

test-deps: mysql/schema.sql $(EXAMPLE_APP).pkg $(EXAMPLE_APP).pkg.modified

mysql/%.sql : force_look
	$(MAKE) -C mysql $(subst mysql/,,$@) $(MFLAGS)

$(EXAMPLE_APP).pkg : $(EXAMPLE_APP_FILES)
	zip -r $@ $(EXAMPLE_APP).app

$(EXAMPLE_APP).pkg.modified : $(EXAMPLE_APP_FILES)
	echo "some data" > $(EXAMPLE_APP).app/extra_data
	zip -r $@ $(EXAMPLE_APP).app
	rm -f $(EXAMPLE_APP).app/extra_data

clean :
	@for d in $(DIRS) ; do \
		$(MAKE) -C $$d clean $(MFLAGS) ; \
	done
	go clean ./...
	rm -f $(EXAMPLE_APP).pkg
	rm -f $(EXAMPLE_APP).pkg.modified

force_look :
	@true

.PHONY : test integration-test test-deps clean
