USE [test]
GO

/****** Object:  Table [dbo].[People]    Script Date: 2016/3/12 14:27:41 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[People](
	[Age] [int] NOT NULL,
	[Name] [nvarchar](255) NOT NULL,
	[PeopleId] [int] IDENTITY(1,1) NOT NULL,
	[NonIndexA] [nchar](255) NULL,
	[NonIndexB] [nchar](255) NULL,
	[IndexAPart1] [int] NULL,
	[IndexAPart2] [int] NULL,
	[IndexAPart3] [int] NULL,
 CONSTRAINT [PK__people__3214EC27D30B5FDA] PRIMARY KEY CLUSTERED
(
	[PeopleId] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
