--
-- File generated with SQLiteStudio v3.4.4 on Sat Mar 23 23:23:14 2024
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: AssignedAddresses
DROP TABLE IF EXISTS AssignedAddresses;

CREATE TABLE IF NOT EXISTS AssignedAddresses (
    Id           INTEGER  PRIMARY KEY AUTOINCREMENT
                          UNIQUE
                          NOT NULL,
    Address      STRING   UNIQUE
                          NOT NULL,
    HostNameId   INTEGER  NOT NULL
                          REFERENCES Hosts (Id),
    DomainId     INTEGER  REFERENCES Domains (Id) 
                          NOT NULL,
    SubnetId     INTEGER  REFERENCES Subnets (Id) 
                          NOT NULL,
    CreatorId    INTEGER  NOT NULL
                          REFERENCES Users (Id),
    CreationDate DATETIME NOT NULL
                          DEFAULT (CURRENT_TIMESTAMP) 
);


-- Table: Domains
DROP TABLE IF EXISTS Domains;

CREATE TABLE IF NOT EXISTS Domains (
    Id           INTEGER  PRIMARY KEY AUTOINCREMENT
                          NOT NULL
                          UNIQUE,
    DomainName   STRING   UNIQUE
                          NOT NULL,
    CreatorId    INTEGER  REFERENCES Users (Id) 
                          NOT NULL,
    CreationDate DATETIME NOT NULL
                          DEFAULT (CURRENT_TIMESTAMP) 
);


-- Table: Hosts
DROP TABLE IF EXISTS Hosts;

CREATE TABLE IF NOT EXISTS Hosts (
    Id           INTEGER  PRIMARY KEY AUTOINCREMENT
                          NOT NULL
                          UNIQUE,
    HostName     STRING   NOT NULL
                          UNIQUE,
    MacAddresses JSON     NOT NULL,
    CreatorId    INTEGER  REFERENCES Users (Id) 
                          NOT NULL,
    CreationDate DATETIME NOT NULL
                          DEFAULT (CURRENT_TIMESTAMP) 
);


-- Table: Subnets
DROP TABLE IF EXISTS Subnets;

CREATE TABLE IF NOT EXISTS Subnets (
    Id             INTEGER  NOT NULL
                            UNIQUE
                            PRIMARY KEY AUTOINCREMENT,
    NetworkName    STRING   NOT NULL
                            UNIQUE,
    NetworkPrefix  STRING   NOT NULL
                            UNIQUE,
    BitMask        INTEGER  NOT NULL,
    GatewayAddress STRING   NOT NULL,
    DomainId       INTEGER  NOT NULL
                            REFERENCES Domains (Id),
    CreatorId      INTEGER  REFERENCES Users (Id) 
                            NOT NULL,
    CreationDate   DATETIME NOT NULL
                            DEFAULT (CURRENT_TIMESTAMP) 
);


-- Table: Users
DROP TABLE IF EXISTS Users;

CREATE TABLE IF NOT EXISTS Users (
    Id              INTEGER  PRIMARY KEY AUTOINCREMENT
                             NOT NULL
                             UNIQUE,
    UserName        STRING   UNIQUE
                             NOT NULL,
    Status          STRING   DEFAULT enabled
                             NOT NULL,
    PasswordHash    STRING   NOT NULL,
    CreationDate    DATETIME NOT NULL
                             DEFAULT (CURRENT_TIMESTAMP),
    LastChangedDate DATETIME NOT NULL
                             DEFAULT (CURRENT_TIMESTAMP) 
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
