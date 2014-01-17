DIRS = mysql

EXAMPLE_APP = example
EXAMPLE_APP_FILES = $(shell find $(EXAMPLE_APP).app/*.js) $(shell find $(EXAMPLE_APP).app/*.html)

test: test-deps
	go test -run=UnitSpecs

integration-test: test-deps
	go test -run=IntegrationSpecs

test-deps: mysql/schema.sql $(EXAMPLE_APP).pkg $(EXAMPLE_APP).pkg.1 $(EXAMPLE_APP).pkg.2 another.pkg

mysql/%.sql : force_look
	$(MAKE) -C mysql $(subst mysql/,,$@) $(MFLAGS)

$(EXAMPLE_APP).pkg : $(EXAMPLE_APP_FILES)
	zip -r $@ $(EXAMPLE_APP).app

$(EXAMPLE_APP).pkg.1 : $(EXAMPLE_APP_FILES)
	echo "1" > $(EXAMPLE_APP).app/extra_data
	zip -r $@ $(EXAMPLE_APP).app
	rm -f $(EXAMPLE_APP).app/extra_data

$(EXAMPLE_APP).pkg.2 : $(EXAMPLE_APP_FILES)
	echo "2" > $(EXAMPLE_APP).app/extra_data
	zip -r $@ $(EXAMPLE_APP).app
	rm -f $(EXAMPLE_APP).app/extra_data

another.pkg : $(EXAMPLE_APP).pkg
	echo "another" > $(EXAMPLE_APP).app/extra_data
	zip -r $@ $(EXAMPLE_APP).app
	rm -f $(EXAMPLE_APP).app/extra_data

clean :
	@for d in $(DIRS) ; do \
		$(MAKE) -C $$d clean $(MFLAGS) ; \
	done
	go clean ./...
	rm -f $(EXAMPLE_APP).pkg
	rm -f $(EXAMPLE_APP).pkg.1
	rm -f $(EXAMPLE_APP).pkg.2
	rm -f another.pkg

force_look :
	@true

.PHONY : test integration-test test-deps clean
