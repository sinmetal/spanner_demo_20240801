CREATE TABLE SampleMessages (
  SampleMessageID STRING(MAX) NOT NULL,
  Message STRING(MAX) NOT NULL,
  SampleMessages_Message_Tokens TOKENLIST AS (TOKENIZE_FULLTEXT(Message)) HIDDEN
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (SampleMessageID);

CREATE SEARCH INDEX SampleMessagesMessageIndex
  ON SampleMessages(SampleMessages_Message_Tokens);