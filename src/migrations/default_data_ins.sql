-- graphs
INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (191,'Отсутствует легенда на графике',1,'no_graph_leg');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (192,'Отсутствует подпись к осям на графике',1,'no_graph_annot');


-- schemes


INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (161,'Неверное расположение стрелок на графиках',1,'wrong_scheme_arrows');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (162,'Отсутствует подпись (да, нет) к блоку ветвления',1,'wrong_scheme_if');


INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (163,'Неверный формат терминаторов схемы алгоритма',1,'wrong_terminators');

-- tables

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (171,'Отсутствует подпись таблицы',1,'no_table_annot');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (172,'Подпись таблицы неверна',1,'wrong_table_annot');

-- formulas coming soon

-- extra

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (0,'Ошибок нет, все хорошо))',1,'no_errors');






