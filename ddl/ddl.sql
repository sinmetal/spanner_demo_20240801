CREATE TABLE SampleMessages (
  SampleMessageID STRING(MAX) NOT NULL,
  Message STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (SampleMessageID);