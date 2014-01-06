DIRS = mysql

test:
	go test -run=UnitSpecs

integration-test: mysql/schema.sql
	go test -run=IntegrationSpecs

mysql/%.sql : force_look
	$(MAKE) -C mysql $(subst mysql/,,$@) $(MFLAGS)

clean :
	@for d in $(DIRS) ; do \
		$(MAKE) -C $$d clean $(MFLAGS) ; \
	done
	go clean ./...

force_look :
	@true

.PHONY : test integration-test clean
