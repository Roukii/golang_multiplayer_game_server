create type stats(
  level int, 
  maxHp Bigint, 
  hp Bigint, 
  maxMp Bigint, 
  Mp Bigint,
);

create type spawn_point(
  world_uuid Uuid, 
  x float, 
  y float, 
  z float,
  updated_at Timestamp
);

create TABLE world(
  world_uuid Uuid,
  max_player int,
  name text,
  width int,
  length int,
  seed text,
  created_at Timestamp, 
  PRIMARY KEY (world_uuid) 
);


create TABLE players_by_user(
  user_uuid Uuid,
  player_uuid Uuid,
  name text,
  stats stats, 
  spawn_point spawn_point, 
  created_at Timestamp, 
  updated_at Timestamp,
  PRIMARY KEY (user_uuid, player_uuid) 
) WITH CLUSTERING ORDER BY (player_uuid DESC);

create type tile(
  type int,
  elevation float,
);

create TABLE chunks_by_world(
  chunk_uuid Uuid, 
  world_uuid uuid,
  x int,
  y int, 
  created_at Timestamp,
  tiles list<frozen<tile>>,
  PRIMARY KEY (world_uuid) 
);

create TABLE static_entity_by_chunk(
  static_entity_uuid Uuid, 
  chunk_uuid Uuid, 
  name text, 
  x float, 
  y float,
  z float,
  entity_type int, 
  stats stats, 
  entry_to_chunk_uuid uuid, 
  created_at Timestamp, 
  updated_at Timestamp,  
  PRIMARY KEY (chunk_uuid, static_entity_uuid)
) WITH CLUSTERING ORDER BY (static_entity_uuid DESC);

create TABLE dynamic_entity_by_chunk(
  dynamic_entity_uuid Uuid, 
  chunk_uuid Uuid, 
  name text, 
  x float, 
  y float, 
  z float,
  entity_type int, 
  stats stats, 
  entry_to_chunk_uuid uuid, 
  created_at Timestamp, 
  updated_at Timestamp,  
  PRIMARY KEY (chunk_uuid, dynamic_entity_uuid)
) WITH CLUSTERING ORDER BY (dynamic_entity_uuid DESC);

CREATE INDEX  dynamic_entity_by_chunk_x_idx ON dynamic_entity_by_chunk (x);
CREATE INDEX  dynamic_entity_by_chunk_y_idx ON dynamic_entity_by_chunk (y);

CREATE INDEX chunks_by_world_y_idx ON chunks_by_world (y);
CREATE INDEX chunks_by_world_x_idx ON chunks_by_world (x);

CREATE INDEX  static_entity_by_chunk_y_idx ON static_entity_by_chunk (y);
CREATE INDEX static_entity_by_chunk_x_idx ON static_entity_by_chunk (x);
