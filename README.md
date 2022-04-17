# pock_multiplayer
pock for a multiplayer game

TODO Upgrade :
  - Clean entity folder
  - Handle generation of StaticEntity with update and streaming 
  - Batch RPC Message from stream
  - Put Mutex where it's needed (worldService, chunkService, DynamicEntityService)
  - Move Player between world service when they change world
  - Batch client request and proccess ?
  - Start/Stop world service by the player count
  - Create Equipment, Inventory and Stats fonctionality
  - Create ennemy dynamic entity and add AI
  - Upgrade world generation

TODO Fonctionality:
  - Interactable Entity Generation
  - Enemy Entity Generation
  - Enemy AI
  - Player Movement
  - Player Attack
  - Collision system
  - Inventory
  - Equipment
  - Skill
  - Player Interaction With Interactable Entity
  - Building Creation


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
// write it down you lazy bastard

SQL Database : 
    - users
    - user_world_affs
    - connexions
    - worlds


### Install : 
Create databases with this command (need docker)
- docker compose -f docker-compose.dev.yml up -d

## Cassandra
  - connect to pock_multiplayer-DC1N1-1 docker container :
    - cqlsh -u cassandra -p cassandra
    - enter this command : '''
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

Generate Proto for GRPC Command
SERVER -> protoc --go_out=. --go_opt=paths=source_relative \
   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   proto/world.proto

CLIENT -> protoc --csharp_out=proto/ --grpc_out=proto/ --plugin=protoc-gen-grpc=proto/grpc_csharp_plugin proto/world.proto

INSERT INTO game.chunks_by_world (world_uuid,chunk_uuid,x,y,tiles,created_at) VALUES (e0a44f54-b4de-11ec-9e28-367dda4cfa8c, 813e3fac-b4de-11ec-be67-367dda4cfa8c, 1, 1, [{tile_type : 1,elevation:0}], '2018-02-07 14:07:00');
  
  SyntaxException: line 1:124 no viable alternative at input ',' (...y,tiles,created_at) VALUES (["e0a44f54-b4de-11ec-9e28-367dda4cfa8]c",...)