CREATE KEYSPACE IF NOT EXISTS keyspace_1 WITH REPLICATION = { 
    'class' : 'SimpleStrategy', 
    'replication_factor' : 1 
};
