CREATE SCHEMA testtask;

CREATE TABLE testtask.locations (
    location_id bigint NOT NULL PRIMARY KEY,
    location varchar NOT NULL,
    coordinates point
);
COMMENT ON TABLE testtask.locations IS 'Адреса';
COMMENT ON COLUMN testtask.locations.location IS 'Информация об адресной точке в формате json';
COMMENT ON COLUMN testtask.locations.coordinates IS 'Адресная точка';

CREATE TABLE testtask.item_locations (
    item_id bigint NOT NULL,
    location_id bigint NOT NULL
        REFERENCES testtask.locations(location_id) ON UPDATE CASCADE ON DELETE CASCADE,
    UNIQUE(item_id, location_id)
);

COMMENT ON TABLE testtask.item_locations IS 'Хранит привязанные к объявлению адреса';
COMMENT ON COLUMN testtask.item_locations.item_id IS 'Объявление';
COMMENT ON COLUMN testtask.item_locations.location_id IS 'Адресная точка';
