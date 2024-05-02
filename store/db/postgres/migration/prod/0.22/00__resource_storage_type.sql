ALTER TABLE `resource`
  ADD COLUMN `storage_type` VARCHAR(256) NOT NULL DEFAULT '',
  ADD COLUMN `reference` VARCHAR(256) NOT NULL DEFAULT '',
  ADD COLUMN `payload` TEXT NOT NULL;

UPDATE `resource`
SET `storage_type` = 'LOCAL', `reference` = `internal_path`
WHERE `internal_path` IS NOT NULL AND `internal_path` != '';

UPDATE `resource`
SET `storage_type` = 'EXTERNAL', `reference` = `external_link`
WHERE `external_link` IS NOT NULL AND `external_link` != '';

ALTER TABLE `resource`
  DROP COLUMN `internal_path`,
  DROP COLUMN `external_link`;
