
#!/bin/bash

database="inventory.db"

# if database doesn't exist, create it
if [ ! -f $database ]; then
  echo "Creating inventory.db"
  exec `sqlite3 inventory.db < DDL.sql`

  echo "Pouring something into it..."
  exec `sqlite3 inventory.db < DML.sql`

  # is it created?
  if [ -f $database ]; then
    echo "inventory.db is ready to Go"
  fi
fi

# When somebody press Ctrl-C
trap '{ echo "Hey, you pressed Ctrl-C. Time to quit." ; exit 1; }' INT

# start goventory microservice
echo "Starting goventory on Port :3000. Press Ctrl-C to quit."
exec `PORT=3000 go run microservice.go`
