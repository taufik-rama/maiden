CREATE KEYSPACE IF NOT EXISTS keyspace_2 WITH REPLICATION = { 
    'class' : 'SimpleStrategy', 
    'replication_factor' : 1 
};
