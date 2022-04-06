# pock_multiplayer
pock for a multiplayer game


Gateway/Gateway

Route :
  - /auth/login
  - /auth/logout
  - /auth/refresh
  - /auth/reset_password
  - /auth/register
  - /world/user
  - /world/list
  - /world/join
  - /world/delete_user


Cassandra Database :
  - id_world : 
    - id_user :
      - name
      - posXY
      - customization :
        - hair
        - nose
        - ...
      - stats
      - equipment :
        - weapon
        - armor
        - ...
      - inventory
      - job

SQL Database : 
    - users
    - user_world_affs
    - connexions
    - worlds


RPC Command
SERVER -> protoc --go_out=. --go_opt=paths=source_relative \
   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   proto/world.proto

CLIENT -> protoc --csharp_out=proto/ --grpc_out=proto/ --plugin=protoc-gen-grpc=proto/grpc_csharp_plugin proto/world.proto

Fonctionality:
  - World Procedural Generation (Noise map)
  - Interactable Entity Generation
  - Enemy Entity Generation
  - Enemy AI
  - Player Movement
  - Player Attack
  - Inventory
  - Equipment
  - Skill
  - Player Interaction With Interactable Entity
  - Building Creation

## Univers Generation
# World :
UUID string
Name string
Level smallint
Length smallint
Width smallint
Chunks []Chunk
seed    string
type    WorldType
...

# Chunk :
UUID string
updatedAt timestamp
createdAt timestamp
PositionX int
PositionY int
Tiles []Tile
Entity []Entity
state    ChunkState

# Chunk State :
- Normal
- Combat

# Tile
type TileType
elevation smallint 
updatedAt timestamp

# Tile type:
- dirt
- grass
- rock
- forest
- sand
- snow

# Entity
positionX smallint
positionY smallint
Entity Entity
type      EntityType
Stats     Stats
entryChunkUUID string
updatedAt timestamp
createdAt timestamp

# Entity Type :
- Empty
- Creature
- Player
- Building
- Chunk
- Ressource
- Item

# Stats :
Level smallint
MaxHP int
HP int
MaxMP int
MP int


## Cassandra
  - connect to pock_multiplayer-DC1N1-1 docker container :
    - cqlsh -u cassandra -p cassandra
    - '''sql
          CREATE KEYSPACE IF NOT EXISTS game
        WITH REPLICATION = {
        'class' : 'SimpleStrategy',
        'replication_factor' : 1
        };
      '''
  - brew install golang-migrate
  - create migration file : migrate create -ext sql -dir migrations -seq create_users_table
  -
  - migrate -database "cassandra://127.0.0.1:9042/game?sslmode=disable&x-multi-statement=true" -path cmd/world/migrations up
  - migrate -database "cassandra://127.0.0.1:9042/game?sslmode=disable&x-multi-statement=true" -path cmd/world/migrations down
  - migrate -database "cassandra://127.0.0.1:9042/game?sslmode=disable&x-multi-statement=true" -path cmd/world/migrations drop

INSERT INTO game.chunks_by_world (world_uuid,chunk_uuid,x,y,tiles,created_at) VALUES (e0a44f54-b4de-11ec-9e28-367dda4cfa8c, 813e3fac-b4de-11ec-be67-367dda4cfa8c, 1, 1, [{tile_type : 1,elevation:0}], '2018-02-07 14:07:00');
  
  SyntaxException: line 1:124 no viable alternative at input ',' (...y,tiles,created_at) VALUES (["e0a44f54-b4de-11ec-9e28-367dda4cfa8]c",...)