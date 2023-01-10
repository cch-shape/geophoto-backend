# Create test user
CREATE USER 'testuser'@'%' IDENTIFIED BY 'testpassword';
GRANT ALL PRIVILEGES ON GeoPhoto.* TO 'testuser'@'%';

# Create photo table
create table photo
(
    id          uuid      default uuid()              not null
        primary key,
    user_id     int unsigned                          not null,
    filename    text                                  not null,
    description text                                  null,
    coordinates point                                 not null,
    timestamp   timestamp default current_timestamp() not null,
    created_at  timestamp default current_timestamp() not null,
    updated_at  timestamp default current_timestamp() not null on update current_timestamp()
);

create index photo_created_at_index
    on photo (created_at);

create index photo_timestamp_index
    on photo (timestamp);

create index photo_updated_at_index
    on photo (updated_at);

create index photo_user_id_index
    on photo (user_id);
    