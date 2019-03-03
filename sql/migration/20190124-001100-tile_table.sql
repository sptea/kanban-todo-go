create table if not exists 'tile' ('tile_id' text primary key, 'title' text, 'row_id' text, 'text' text, foreign key(row_id) references row(row_id));
