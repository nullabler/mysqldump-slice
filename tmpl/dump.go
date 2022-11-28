package tmpl

	const dump = `-- Host: {{ .Host }}    Database: {{ .Database }}
-- ------------------------------------------------------
-- Server version	{{ .ServerVersion }}
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: {{ .Database }}
--

{{ .CreateDatabase }}

USE {{ .Database }};

{{ .Tables }}
`

type TemplateParams struct {
	Host string
	Database string
	CreateDatabase string
	ServerVersion string
	Tables string
}

func NewTemplateParams(host string, db string) *TemplateParams {
	return &TemplateParams{
		Host: host,
		Database: db, 
		ServerVersion: "0.0.1",
	}
}

func Dump() string {
	return dump
}
