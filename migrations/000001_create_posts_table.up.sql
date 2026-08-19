CREATE TABLE IF NOT EXISTS `posts` (
  `Id` int unsigned NOT NULL AUTO_INCREMENT,
  `Title` varchar(200) DEFAULT NULL,
  `Content` text,
  `Category` varchar(100) DEFAULT NULL,
  `Created_date` timestamp NULL DEFAULT NULL,
  `Updated_date` timestamp NULL DEFAULT NULL,
  `Status` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;