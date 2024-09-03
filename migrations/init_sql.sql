create table ocr (
    id serial,
    image_url varchar(255) null,
    text text null,
    status varchar(255) null,
    created_at timestamp default now(),
    updated_at timestamp null,
    deleted_at timestamp null
)