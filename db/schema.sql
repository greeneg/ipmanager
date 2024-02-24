--
-- File generated with SQLiteStudio v3.4.4 on Sat Feb 24 16:22:52 2024
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: AssignedAddresses
DROP TABLE IF EXISTS AssignedAddresses;

CREATE TABLE IF NOT EXISTS AssignedAddresses (
    Id             INTEGER  PRIMARY KEY AUTOINCREMENT
                            UNIQUE
                            NOT NULL,
    Address        STRING   UNIQUE
                            NOT NULL,
    SubnetId       INTEGER  REFERENCES Subnets (Id) 
                            NOT NULL,
    AssignmentDate DATETIME NOT NULL
                            DEFAULT (CURRENT_TIMESTAMP),
    AssignmentId   INTEGER  NOT NULL
                            REFERENCES Users (Id),
    BitMask        INTEGER  NOT NULL,
    HostName       STRING   NOT NULL
                            UNIQUE,
    DomainId       INTEGER  REFERENCES Domains (Id) 
                            NOT NULL
);


-- Table: Domains
DROP TABLE IF EXISTS Domains;

CREATE TABLE IF NOT EXISTS Domains (
    Id           INTEGER  PRIMARY KEY AUTOINCREMENT
                          NOT NULL
                          UNIQUE,
    DomainName   STRING   UNIQUE
                          NOT NULL,
    CreationDate DATETIME NOT NULL
                          DEFAULT (CURRENT_TIMESTAMP),
    CreatorId    INTEGER  REFERENCES Users (Id) 
);


-- Table: Subnets
DROP TABLE IF EXISTS Subnets;

CREATE TABLE IF NOT EXISTS Subnets (
    Id                         INTEGER  NOT NULL
                                        UNIQUE
                                        PRIMARY KEY AUTOINCREMENT,
    NetworkName                STRING   NOT NULL
                                        UNIQUE,
    AddressStart               STRING   NOT NULL
                                        UNIQUE,
    BitMask                    INTEGER  NOT NULL,
    NumberOfAvailableAddresses INTEGER  NOT NULL,
    CreationDate               DATETIME NOT NULL
                                        DEFAULT (CURRENT_TIMESTAMP),
    CreatorId                  INTEGER  REFERENCES Users (Id) 
                                        NOT NULL,
    DomainId                   INTEGER  NOT NULL
                                        REFERENCES Domains (Id) 
);


-- Table: Users
DROP TABLE IF EXISTS Users;

CREATE TABLE IF NOT EXISTS Users (
    Id           INTEGER PRIMARY KEY AUTOINCREMENT
                         NOT NULL
                         UNIQUE,
    UserName     STRING  UNIQUE
                         NOT NULL,
    PasswordHash STRING  NOT NULL,
    CreationDate INTEGER NOT NULL
                         DEFAULT (CURRENT_TIMESTAMP) 
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
