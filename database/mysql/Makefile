MIGRATIONS = $(shell find migrations/*.sql)

schema.sql : $(MIGRATIONS)
	database-util -config="../config.json" > schema.sql

clean :
	rm -f schema.sql

.PHONY : clean
