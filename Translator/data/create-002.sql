#  The true identity of an Entry field is its Original and Part Of fields.
#  But MySQL/Maria don't like using blobs as primary keys, so instead we've
#  got a `Key` field with an md5 hash of it.
# 
#  Maximum key length: 767 bytes -> 255 characters
#  Timestamp = 4 bytes -> 2 characters
#  Language = 2 characters


create table Entries (
	EntryID bigint unique not null,
	Original text not null,
	PartOf text not null,
	primary key entry_key (EntryID)
);

create index Entries_PartOf on Entries (PartOf(255));

create table Sources (
	SourceID bigint unique not null primary key,
	Filepath text not null,
	Page varchar(255) not null,
	Volume varchar(255) not null,
	Level int not null,
	Game varchar(64) not null
);

create table EntrySources (
	EntryID bigint not null,
	SourceID bigint not null,
	Count int not null,
	primary key (EntryID, SourceID)
);

create index EntrySources_Reverse on EntrySources (SourceID, EntryID);

create table Translations (
	TranslationID bigint unique not null primary key,
	EntryID bigint not null,
	Language char(2) not null,
	Translator varchar(128) not null,
	Translation text not null,
	IsPreferred boolean not null,
	IsConflicted boolean not null
);

create index Translations_EntryID on Translations (EntryID);

create table Users (
	Email varchar(128) unique not null primary key,
	Password varchar(255) not null,
	Secret varchar(255) not null,
	Name varchar(255) not null,
	IsAdmin boolean not null,
	Language char(2) not null,
	IsLanguageLead boolean not null
);

create index Users_Language on Users (Language);

create table Votes (
	TranslationID bigint not null,
	Voter varchar(128) not null,
	Vote boolean not null,
	primary key (TranslationID, Voter)
);

create index Votes_Reverse on Votes (Voter, TranslationID);