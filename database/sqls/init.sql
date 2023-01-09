CREATE USER 'testuser'@'%' IDENTIFIED BY 'testpassword';
GRANT ALL PRIVILEGES ON GeoPhoto.* TO 'testuser'@'%';

create table GeoPhoto.photo
(
    id          int UNSIGNED auto_increment,
    user_id     int UNSIGNED                        not null,
    photo_url    text                                not null,
    description text                                null,
    timestamp   timestamp default current_timestamp not null,
    coordinates point                               not null,
    created_at  timestamp default current_timestamp not null,
    updated_at  timestamp default current_timestamp not null on update current_timestamp,
    constraint id
        primary key (id)
);

create index photo_created_at_index
    on GeoPhoto.photo (created_at);

create index photo_timestamp_index
    on GeoPhoto.photo (timestamp);

create index photo_updated_at_index
    on GeoPhoto.photo (updated_at);

create index photo_user_id_index
    on GeoPhoto.photo (user_id);

