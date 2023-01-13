# Create test user
CREATE USER 'testuser'@'%' IDENTIFIED BY 'testpassword';
GRANT ALL PRIVILEGES ON GeoPhoto.* TO 'testuser'@'%';

# Photo table
create table photo
(
    id           int unsigned auto_increment
        primary key,
    uuid         uuid      default uuid()              not null,
    user_id      int unsigned                          not null,
    filename     text                                  not null,
    description  text                                  null,
    address_name text                                  null,
    address      text                                  null,
    coordinates  point                                 not null,
    timestamp    timestamp default current_timestamp() not null,
    created_at   timestamp default current_timestamp() not null,
    updated_at   timestamp default current_timestamp() not null on update current_timestamp(),
    constraint photo_uuid_unique
        unique (uuid)
);

create index photo_created_at_index
    on photo (created_at);

create index photo_timestamp_index
    on photo (timestamp);

create index photo_updated_at_index
    on photo (updated_at);

create index photo_user_id_index
    on photo (user_id);

create index photo_uuid_index
    on photo (uuid);

# User table
create table user
(
    id            int auto_increment
        primary key,
    uuid          uuid      default uuid()              not null,
    phone_number  varchar(20)                           not null,
    name          text                                  null,
    thumbnail_url text                                  null,
    created_at    timestamp default current_timestamp() not null,
    updated_at    timestamp default current_timestamp() not null on update current_timestamp()
);

create index user_phone_number_index
    on user (phone_number);

create index user_uuid_index
    on user (uuid);

# Verification code table
create table verification_code
(
    id           int auto_increment
        primary key,
    phone_number varchar(20)                            not null,
    hashed_code  text                                   not null,
    is_voided    tinyint(1) default 0                   not null,
    created_at   timestamp  default current_timestamp() not null,
    updated_at   timestamp  default current_timestamp() not null on update current_timestamp()
);

create index verification_code_phone_number_index
    on verification_code (phone_number);
