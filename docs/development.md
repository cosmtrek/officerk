# development

## Prepare MySQL database

```bash
# enter mysql console
mysql -uroot

# create user
CREATE USER 'officerk'@'localhost' IDENTIFIED BY 'officerkpass';
# create database
CREATE DATABASE officerk_development;

# grant permissions
GRANT ALL on officerk_development.* to 'officerk'@'localhost' IDENTIFIED BY 'officerkpass';
```