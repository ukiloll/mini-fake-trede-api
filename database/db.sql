-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema faketradeapp
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema faketradeapp
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `faketradeapp` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
USE `faketradeapp` ;

-- -----------------------------------------------------
-- Table `faketradeapp`.`assets`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `faketradeapp`.`assets` (
  `asset_id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `type` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`asset_id`),
  UNIQUE INDEX `name` (`name` ASC) VISIBLE,
  UNIQUE INDEX `type` (`type` ASC) VISIBLE)
ENGINE = InnoDB
AUTO_INCREMENT = 3
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `faketradeapp`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `faketradeapp`.`users` (
  `user_id` INT NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(100) NOT NULL,
  `auth_host` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `email` (`email` ASC) VISIBLE,
  UNIQUE INDEX `auth_host` (`auth_host` ASC) VISIBLE)
ENGINE = InnoDB
AUTO_INCREMENT = 2
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `faketradeapp`.`transition`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `faketradeapp`.`transition` (
  `trade_id` INT NOT NULL AUTO_INCREMENT,
  `trade_type` ENUM('Buy', 'Sell') NOT NULL,
  `price` FLOAT NOT NULL,
  `quantity` FLOAT NOT NULL,
  `user_id` INT NOT NULL,
  `asset_id` INT NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`trade_id`),
  INDEX `user_id` (`user_id` ASC) VISIBLE,
  INDEX `asset_id` (`asset_id` ASC) VISIBLE,
  CONSTRAINT `transition_ibfk_1`
    FOREIGN KEY (`user_id`)
    REFERENCES `faketradeapp`.`users` (`user_id`),
  CONSTRAINT `transition_ibfk_2`
    FOREIGN KEY (`asset_id`)
    REFERENCES `faketradeapp`.`assets` (`asset_id`))
ENGINE = InnoDB
AUTO_INCREMENT = 18
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `faketradeapp`.`user_assets`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `faketradeapp`.`user_assets` (
  `user_Assets_id` INT NOT NULL AUTO_INCREMENT,
  `asset_id` INT NOT NULL,
  `user_id` INT NOT NULL,
  `quantity` FLOAT NOT NULL,
  `price` FLOAT NOT NULL,
  PRIMARY KEY (`user_Assets_id`),
  INDEX `user_id` (`user_id` ASC) VISIBLE,
  INDEX `asset_id` (`asset_id` ASC) VISIBLE,
  CONSTRAINT `user_assets_ibfk_1`
    FOREIGN KEY (`user_id`)
    REFERENCES `faketradeapp`.`users` (`user_id`),
  CONSTRAINT `user_assets_ibfk_2`
    FOREIGN KEY (`asset_id`)
    REFERENCES `faketradeapp`.`assets` (`asset_id`))
ENGINE = InnoDB
AUTO_INCREMENT = 5
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
