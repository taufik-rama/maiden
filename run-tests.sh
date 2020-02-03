#!/bin/bash

./build

# Fixture Cassandra
echo "Cassandra push"
./maiden -c examples/basic/maiden.yaml fixture push cassandra -v
echo "Cassandra remove"
./maiden -c examples/basic/maiden.yaml fixture remove cassandra -v

echo "Cassandra push (keyspace invalid)"
./maiden -c examples/basic/maiden.yaml fixture push cassandra --keyspace test -v
echo "Cassandra remove (keyspace invalid)"
./maiden -c examples/basic/maiden.yaml fixture remove cassandra --keyspace test -v

echo "Cassandra push (keyspace)"
./maiden -c examples/basic/maiden.yaml fixture push cassandra --keyspace my_keyspace -v
echo "Cassandra remove (keyspace)"
./maiden -c examples/basic/maiden.yaml fixture remove cassandra --keyspace my_keyspace -v

echo "Cassandra push (keyspace & table invalid)"
./maiden -c examples/basic/maiden.yaml fixture push cassandra --keyspace my_keyspace --table test -v
echo "Cassandra remove (keyspace & table invalid)"
./maiden -c examples/basic/maiden.yaml fixture remove cassandra --keyspace my_keyspace --table test -v

echo "Cassandra push (keyspace & table)"
./maiden -c examples/basic/maiden.yaml fixture push cassandra --keyspace my_keyspace --table my_table -v
echo "Cassandra remove (keyspace & table)"
./maiden -c examples/basic/maiden.yaml fixture remove cassandra --keyspace my_keyspace --table my_table -v

# Fixture Elasticsearch
echo "Elasticsearch push"
./maiden -c examples/basic/maiden.yaml fixture push elasticsearch -v
echo "Elasticsearch remove"
./maiden -c examples/basic/maiden.yaml fixture remove elasticsearch -v

echo "Elasticsearch push (index invalid)"
./maiden -c examples/basic/maiden.yaml fixture push elasticsearch --index test -v
echo "Elasticsearch remove (index invalid)"
./maiden -c examples/basic/maiden.yaml fixture remove elasticsearch --index test -v

echo "Elasticsearch push (index)"
./maiden -c examples/basic/maiden.yaml fixture push elasticsearch --index my-index -v
echo "Elasticsearch remove (index)"
./maiden -c examples/basic/maiden.yaml fixture remove elasticsearch --index my-index -v

# Fixture PostgreSQL
echo "PostgreSQL push"
./maiden -c examples/basic/maiden.yaml fixture push postgresql -v
echo "PostgreSQL remove"
./maiden -c examples/basic/maiden.yaml fixture remove postgresql -v

echo "PostgreSQL push (database invalid)"
./maiden -c examples/basic/maiden.yaml fixture push postgresql --database test -v
echo "PostgreSQL remove (database invalid)"
./maiden -c examples/basic/maiden.yaml fixture remove postgresql --database test -v

echo "PostgreSQL push (database)"
./maiden -c examples/basic/maiden.yaml fixture push postgresql --database my-database -v
echo "PostgreSQL remove (database)"
./maiden -c examples/basic/maiden.yaml fixture remove postgresql --database my-database -v

echo "PostgreSQL push (database & table invalid)"
./maiden -c examples/basic/maiden.yaml fixture push postgresql --database my-database --table test -v
echo "PostgreSQL remove (database & table invalid)"
./maiden -c examples/basic/maiden.yaml fixture remove postgresql --database my-database --table test -v

echo "PostgreSQL push (database & table)"
./maiden -c examples/basic/maiden.yaml fixture push postgresql --database my-database --table my-table -v
echo "PostgreSQL remove (database & table)"
./maiden -c examples/basic/maiden.yaml fixture remove postgresql --database my-database --table my-table -v

# Fixture Redis
echo "Redis push"
./maiden -c examples/basic/maiden.yaml fixture push redis -v
echo "Redis remove"
./maiden -c examples/basic/maiden.yaml fixture remove redis -v
