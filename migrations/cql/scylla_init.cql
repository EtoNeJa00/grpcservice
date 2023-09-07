create keyspace grpcservice with replication = {
        'class': 'SimpleStrategy',
        'replication_factor': 1
    };
create table  grpcservice.records
(
    id uuid,
    record text,
    primary key (id)
);