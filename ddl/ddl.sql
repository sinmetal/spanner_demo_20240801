CREATE TABLE SampleMessages (
  SampleMessageID STRING(MAX) NOT NULL,
  Category STRING(MAX),
  Title STRING(MAX),
  Message STRING(MAX) NOT NULL,
  Tags ARRAY<STRING(MAX)>,
  SampleMessages_Title_Tokens TOKENLIST AS (TOKENIZE_FULLTEXT(Title)) HIDDEN,
  SampleMessages_Message_Tokens TOKENLIST AS (TOKENIZE_FULLTEXT(Message)) HIDDEN,
  SampleMessages_Tags_Tokens TOKENLIST AS (TOKENIZE_FULLTEXT(Tags)) HIDDEN,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (SampleMessageID);

CREATE INDEX SampleMessagesCategoryIndex
    ON SampleMessages(Category);

CREATE SEARCH INDEX SampleMessagesTitleIndex
  ON SampleMessages(SampleMessages_Title_Tokens);

CREATE SEARCH INDEX SampleMessagesMessageIndex
  ON SampleMessages(SampleMessages_Message_Tokens);

CREATE SEARCH INDEX SampleMessagesTitleMessageIndex ON SampleMessages(SampleMessages_Title_Tokens, SampleMessages_Tags_Tokens);

CREATE SEARCH INDEX SampleMessagesTagsIndex
  ON SampleMessages(SampleMessages_Tags_Tokens);