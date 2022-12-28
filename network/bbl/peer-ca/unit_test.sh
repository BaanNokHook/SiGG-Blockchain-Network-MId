#!/bin/bash -l

echo "Waiting for ca server starting up ..."
while [ "$(docker inspect -f {{.State.Running}} $(docker ps -q -f name=ca.*.mobileid.com))" == false ]; do sleep 1; done

echo " echo \"Testing database connection\"
  echo \"List certificate files\"
  ls /etc/database-client/
  echo \"Connection string: \${FABRIC_CA_SERVER_DB_DATASOURCE}\"
  psql \"\${FABRIC_CA_SERVER_DB_DATASOURCE}\" -c \"SELECT id FROM users\"
" | docker exec -i $(docker ps -q -f name=ca.*.mobileid.com) bash --
