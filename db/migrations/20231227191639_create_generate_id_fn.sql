-- atlas:delimiter --end
CREATE OR REPLACE FUNCTION generate_id(OUT result BIGINT) AS $$
DECLARE
    our_epoch  BIGINT := 1667969486860;
    seq_id     BIGINT;
    now_millis BIGINT;
    -- the id of this DB shard, must be set for each
    -- schema shard you have - you could pass this as a parameter too
    shard_id   INT    := 1;
BEGIN
    SELECT nextval('global_id_seq') % 1024 INTO seq_id;

    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - our_epoch) << 23;
    result := result | (shard_id << 10);
    result := result | (seq_id);
END;
$$ LANGUAGE PLPGSQL;
--end