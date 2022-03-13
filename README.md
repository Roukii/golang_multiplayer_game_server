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


NoSQL Database :
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
protoc --go_out=. --go_opt=paths=source_relative \
   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   proto/world.proto