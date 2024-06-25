CREATE TABLE IF NOT EXISTS predict_model (
    id bigserial not null constraint predict_model_pk primary key,
    name varchar(255) not null,
    master_segment_id int not null,
    train_configurations jsonb,
    status varchar(255) not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

DROP TRIGGER IF EXISTS set_timestamp_predict_model on "predict_model";
CREATE TRIGGER set_timestamp_predict_model
    BEFORE UPDATE ON "predict_model"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();