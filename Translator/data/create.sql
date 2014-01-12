/**
 *  The true identity of an Entry field is its Original and Part Of fields.
 *  But MySQL/Maria don't like using blobs as primary keys, so instead we've
 *  got a `Key` field with an md5 hash of it.
 *
 *  Maximum key length: 767 bytes -> 255 characters
 */

create table Entries (
	Original text not null,
	PartOf text not null,
	Key char(32) unique not null primary key
);

create index Entries_PartOf on Entries (PartOf(255))

create table EntrySources (
	Entry char(32) not null,
	SourcePath text not null,
	primary key (Entry, SourcePath(240))
);

create table Sources (
	Filepath text not null,
	Page varchar(255) not null,
	Volume varchar(255) not null,
	SourceGroup varchar(255) not null,
	Game varchar(255) not null,
	primary key source_filepath (Filepath(255))
);

create table Translations (
	Entry char(32) not null,
	Language char(2) not null,
	Translator varchar(128) not null,
	Translation text not null,
	primary key (Entry, Language, Translator)
);

create table Comments (
	Entry char(32) not null,
	Language char(2) not null,
	Commenter varchar(128) not null,
	Comment text not null,
	CommentDate timestamp not null,
	primary key (Entry, Language, Commenter, CommentDate)
);

create table Users (
	Email varchar(128) unique not null primary key,
	Password varchar(255) not null,
	Secret varchar(255) not null,
	Name varchar(255) not null,
	IsAdmin boolean not null,
	Language char(2) not null,
);
