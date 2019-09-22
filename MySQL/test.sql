Use test_db;

Create table `TechTools` (id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100) NOT NULL);
Create table ResourceList (id int, link varchar(200), tech_id SMALLINT UNSIGNED NOT NULL, constraint `fk_techtool_resource` FOREIGN KEY (tech_id) REFERENCES TechTools (id));

Insert into TechTools (id, name) values (2, 'Golang');
Insert into TechTools (id, name) values (1, 'Python');
Insert into TechTools (id, name) values (3, 'MySQL');
Insert into TechTools (id, name) values (4, 'Docker');

Insert into ResourceList (id, link, tech_id) values (1, 'link_python1', 1);
Insert into ResourceList (id, link, tech_id) values (2, 'link_python2', 1);
Insert into ResourceList (id, link, tech_id) values (3, 'link_golang', 2);
Insert into ResourceList (id, link, tech_id) values (4, 'link_golang2', 2);
Insert into ResourceList (id, link, tech_id) values (5, 'link_MySQL', 3);
Insert into ResourceList (id, link, tech_id) values (6, 'link_docker', 4);
Insert into ResourceList (id, link, tech_id) values (7, 'link_MySQL', 4);

