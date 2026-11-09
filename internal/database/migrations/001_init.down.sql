DROP TRIGGER IF EXISTS trg_requests_inc ON requests;
DROP FUNCTION IF EXISTS bump_bin_counter();

DROP TABLE IF EXISTS replays;
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS bins;