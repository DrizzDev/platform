CREATE TABLE journal (
    id       INTEGER PRIMARY KEY,
    trace    TEXT    NOT NULL,
    span     TEXT    NOT NULL,
    parent   TEXT    NOT NULL,
    sequence INTEGER NOT NULL,
    mark     TEXT    NOT NULL,
    origin   TEXT    NOT NULL,
    fidelity TEXT    NOT NULL,
    category TEXT    NOT NULL,
    payload  BLOB    NOT NULL,
    artifact TEXT    NOT NULL DEFAULT '',
    state    TEXT    NOT NULL,
    stamped  INTEGER NOT NULL
) STRICT;

CREATE INDEX journal_trace ON journal (trace, sequence);

CREATE TABLE lease (
    trace TEXT    PRIMARY KEY,
    until INTEGER NOT NULL
) STRICT;
