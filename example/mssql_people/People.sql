USE [test];
/****** Object:  Table [dbo].[People]    Script Date: 2016/3/18 17:11:31 ******/
SET ANSI_NULLS ON;
SET QUOTED_IDENTIFIER ON;
CREATE TABLE [dbo].[People](
	[Age] [int] NOT NULL,
	[Name] [nvarchar](255) NOT NULL,
	[PeopleId] [int] IDENTITY(1,1) NOT NULL,
	[NonIndexA] [nvarchar](255) NULL,
	[NonIndexB] [nvarchar](255) NULL,
	[IndexAPart1] [int] NULL,
	[IndexAPart2] [int] NULL,
	[IndexAPart3] [int] NULL,
	[UniquePart1] [int] NOT NULL,
	[UniquePart2] [int] NOT NULL,
	[CreateDate] [datetime] NOT NULL,
	[UpdateDate] [datetime] NOT NULL,
 CONSTRAINT [PK_People] PRIMARY KEY CLUSTERED
(
	[PeopleId] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY];

/****** Object:  Index [IX_AGE]    Script Date: 2016/3/18 17:11:31 ******/
CREATE INDEX [IX_AGE] ON [dbo].[People](
	[Age]
)
/****** Object:  Index [IX_INDEX_A]    Script Date: 2016/3/18 17:11:31 ******/
CREATE INDEX [IX_INDEX_A] ON [dbo].[People](
	[NonIndexA]
)

/****** Object:  Index [IX_NAME]    Script Date: 2016/3/18 17:11:31 ******/
CREATE UNIQUE INDEX [IX_NAME] ON [dbo].[People]
(
	[Name] ASC
)
/****** Object:  Index [IX_UNIQ_A]    Script Date: 2016/3/18 17:11:31 ******/
CREATE UNIQUE INDEX [IX_UNIQ_A] ON [dbo].[People]
(
	[UniquePart1] ASC,
	[UniquePart2] ASC
)
