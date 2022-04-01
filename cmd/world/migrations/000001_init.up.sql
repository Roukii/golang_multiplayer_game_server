create type stats(
  level int, 
  maxHp Bigint, 
  hp Bigint, 
  maxMp Bigint, 
  Mp Bigint, 
  created_at Timestamp, 
  update_at Timestamp
);

create type static_entity(
  uuid Uuid, 
  name text, 
  positionX int, 
  positionY int, 
  entity_type int, 
  stats frozen<stats>, 
  entry_to_chunk_uuid uuid, 
  created_at Timestamp, 
  update_at Timestamp
);

create type dynamic_entity(
  uuid Uuid, 
  name text, 
  position_x int, 
  position_y int, 
  entity_type int, 
  stats frozen<stats>,
  created_at Timestamp, 
  Update_at Timestamp);

create type tile(
  type int,
  elevation int,
  created_at Timestamp, 
  update_at Timestamp
);


create type chunk(
  uuid Uuid, 
  name text, 
  positionX int,
  positionY int, 
  tiles list<frozen<tile>>, 
  static_entities set<frozen<static_entity>>, 
  dynamic_entities set<frozen<dynamic_entity>>, 
  state int, 
  created_at Timestamp, 
  update_at Timestamp
);



create type world(
  uuid Uuid, 
  name text, 
  level int, 
  length int, 
  width int, 
  chunks set<frozen<chunk>>, 
  seed text, 
  type int, 
  created_at Timestamp, 
  update_at Timestamp
);

create type spawn_point(
  world_uuid Uuid, 
  chunk_uuid Uuid, 
  position_x int, 
  position_y int, 
  created_at Timestamp, 
  update_at Timestamp
);


create type player(
  uuid Uuid, 
  name text, 
  level int, 
  stats frozen<stats>, 
  spawn_point frozen<spawn_point>, 
  type int, 
  created_at Timestamp, 
  update_at Timestamp
);



create TABLE universe(
  uuid Uuid,
  name text,
  worlds set<frozen<world>>,
  players set<frozen<player>>,
  created_at Timestamp, 
  PRIMARY KEY(uuid)
);


-- 