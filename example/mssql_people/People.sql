/*
Navicat SQL Server Data Transfer

Source Server         : localtest
Source Server Version : 110000
Source Host           : localhost:1433
Source Database       : test
Source Schema         : dbo

Target Server Type    : SQL Server
Target Server Version : 110000
File Encoding         : 65001

Date: 2016-03-12 12:27:00
*/


-- ----------------------------
-- Table structure for People
-- ----------------------------
DROP TABLE [dbo].[People]
GO
CREATE TABLE [dbo].[People] (
[Age] int NOT NULL ,
[Name] nvarchar(255) NOT NULL ,
[PeopleId] int NOT NULL IDENTITY(1,1) ,
[NonIndexA] nchar(255) NULL ,
[NonIndexB] nchar(255) NULL ,
[IndexAPart1] int NULL ,
[IndexAPart2] int NULL ,
[IndexAPart3] int NULL 
)


GO
DBCC CHECKIDENT(N'[dbo].[People]', RESEED, 1590)
GO

-- ----------------------------
-- Indexes structure for table People
-- ----------------------------
CREATE INDEX [IX_people_age] ON [dbo].[People]
([Age] ASC) 
GO
CREATE UNIQUE INDEX [IX_people_name] ON [dbo].[People]
([Name] ASC) 
WITH (IGNORE_DUP_KEY = ON)
GO
CREATE INDEX [IX_People_A] ON [dbo].[People]
([IndexAPart1] ASC, [IndexAPart2] ASC, [IndexAPart3] ASC) 
GO

-- ----------------------------
-- Primary Key structure for table People
-- ----------------------------
ALTER TABLE [dbo].[People] ADD PRIMARY KEY ([PeopleId])
GO
